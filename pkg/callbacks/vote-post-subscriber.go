package callbacks

import (
	"fmt"
	"math/rand"
	"notification-service/pkg/constants"
	"notification-service/pkg/firebase"
	"notification-service/pkg/i18n"
	"notification-service/pkg/mongo"
	"notification-service/pkg/redis"
	"notification-service/pkg/subscribers"
	"notification-service/pkg/utils"
	"strconv"
	"strings"
	"time"

	"github.com/globalsign/mgo/bson"
)

/**
This processes Upvotes from Posts
1) Get Voting Details
2) Validate only upvote & self vote
3) Get Who created the post -> He gets the notification and we need his current lang
4) Get tokens from the above profile (Supports multiple device tokens.)
5) Create/Update Notification Table which has the meta info for the notificaiotn
6) Construct Data for i18n template
7) Generate template using template data and String Formatter
8) Create push notification
9) Fire the notifications in routines.

*/
func VotePostSubscriberCB(subj, reply string, v *subscribers.Vote) {
	//fmt.Printf("\nNats MSG %+v", v)
	fmt.Printf("Received a vote on subject %s! with vote %+v\n", subj, v)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in subscribers.VotePostSubscriber", r)
		}
	}()
	dir, err := strconv.Atoi(v.Direction)

	if err != nil {
		fmt.Println("Invalid vote")
		return
	}

	err, post := mongo.GetPostById(v.Resource)
	if err != nil {
		return
	}

	fmt.Println("***VOTE COUNT**: ", v.UpVotes)

	vote := post.GetVote(v.Id)

	// create community stats
	community := mongo.GetCommunityById(post.CommunityIds[0].Hex())
	mongo.CreateCommunityStats(mongo.CommunityStatsModel{
		CommunityId:           post.CommunityIds[0],
		Action:                "vote",
		EntityId:              vote.Id,
		EntityType:            "vote",
		ProfileId:             post.Created.ProfileId,
		ActionSource:          post.SourcedBy,
		CommunityCreatorType:  community.Created.Type,
		ActorType:             vote.Created.UserType,
		ParticipatingEntityId: post.Id,
	})

	if vote.Created.ProfileId == post.Created.ProfileId {
		//Self Vote
		fmt.Println("Self Vote")
		return
	}

	// create postCreator's coin
	mongo.CreateUserCoin(mongo.UserCoinsModel{
		ProfileId:   post.Created.ProfileId,
		CoinsEarned: 1,
		Action:      "vote",
	})

	if dir < 1 {
		//Do not process downvotes and unvote
		mongo.RemoveParticipants((post.Id.Hex() + "_vote"), false, vote.Created.ProfileId)
		return
	}

	entities := []mongo.Entity{
		{
			EntityId:   post.Id,
			EntityType: "post",
		},
		{
			EntityId:   vote.Id,
			EntityType: "vote",
		},
	}

	// update feed base time stamp feed_base_ts
	if post.UpVotes == 100 {
		mongo.UpdateOnePostsByQuery(bson.M{"_id": post.Id}, bson.M{"$set": bson.M{"feed_base_ts": time.Now()}})
	}

	userLikesNo := mongo.CountVoteByQuery(bson.M{
		"resource":     post.Id,
		"created.type": bson.M{"$ne": "bot"},
	})

	// schedule likes if it is good posts and liked by 10 users
	if post.PostLevel == "1" && userLikesNo == 10 {
		// count scheduled likes
		totalScheduledLikes := mongo.CountScheduledVotesByQuery(bson.M{
			"resource": post.Id,
			"deleted":  bson.M{"$ne": true},
		})

		min := 101 - totalScheduledLikes - post.UpVotes
		max := 120 - totalScheduledLikes - post.UpVotes

		m, botProfilesHi := mongo.GetBotProfilesIds("hi")
		n, botProfilesTe := mongo.GetBotProfilesIds("te")
		n = m + n
		botProfilesIds := append(botProfilesHi, botProfilesTe...)
		// shuffle profiles
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(n, func(i, j int) { botProfilesIds[i], botProfilesIds[j] = botProfilesIds[j], botProfilesIds[i] })

		voteCreatorsList := mongo.GetAllVotedUserIncludingScheduled(bson.M{"resource": post.Id})

		botProfiles := utils.Difference(botProfilesIds, voteCreatorsList)
		// get unique bot profiles

		noOfVotes := utils.Random(min, max)
		t := utils.SplitTimeInRange(1, 30, noOfVotes, time.Minute)
		j := 0
		for k := 0; j < noOfVotes; j, k = j+1, k+1 {
			vote := mongo.CreateVotesSchedulePost(t[k], bson.ObjectIdHex(post.Id.Hex()), bson.ObjectIdHex(botProfiles[j]))
			mongo.AddVoteSchedule(vote)
		}
	}

	// schedule vote for the user likes
	if vote.Created.UserType != "bot" {
		botProfiles := mongo.GetBotProfileByBucketId(userLikesNo - 1)
		if len(botProfiles) != 0 {
			// schedule two votes
			randomIndexes := utils.GetNRandom(0, 50, 2)
			randomProfile := []string{botProfiles[randomIndexes[0]], botProfiles[randomIndexes[1]]}
			j := 0
			noOfVotes := 2
			t := utils.SplitTimeInRange(1, 15, noOfVotes, time.Minute)
			for k := 0; j < noOfVotes; j, k = j+1, k+1 {
				vote := mongo.CreateVotesSchedulePost(t[k], bson.ObjectIdHex(post.Id.Hex()), bson.ObjectIdHex(randomProfile[j]))
				mongo.AddVoteSchedule(vote)
			}
		}
	}

	// increase post views for bot likes
	if vote.Created.UserType == "bot" {
		mongo.UpdateOnePostsByQuery(bson.M{
			"_id": post.Id,
		}, bson.M{
			"$inc": bson.M{"no_of_views": utils.Random(5, 10)},
		})
	}

	postCreator := mongo.GetProfileById(post.Created.ProfileId)
	count := post.UpVotes
	// notification for karma points
	// karma notification
	karmaBucket := []int{100}
	karmaNextBucket := getNextBucket(post.NotifiedKarmaBuckets, karmaBucket)
	if isValidBucket(karmaNextBucket, count, post.NotifiedKarmaBuckets) {
		bucketStr := strconv.Itoa(karmaNextBucket)
		key := post.Id.Hex() + ":karma:" + bucketStr
		result, err := redis.AcquireLock(key, bucketStr, 60)
		if err == nil && result == 1 {
			templateName := "post_karma_points"
			data := i18n.DataModel{
				Count: post.UpVotes,
			}
			msgStr := i18n.GetString(postCreator.Language, templateName, data)
			htmlMsgStr := i18n.GetHtmlString(postCreator.Language, templateName, data)
			title := i18n.GetAppTitle(postCreator.Language)

			messageMeta := mongo.MessageMeta{
				TemplateName: templateName,
				Template:     i18n.Strings[postCreator.Language][templateName],
				Data:         data,
			}
			// update notification message
			deepLink := "manch://posts/" + post.Id.Hex()
			notification := mongo.CreateNotification(mongo.NotificationModel{
				Receiver:        postCreator.Id,
				Identifier:      postCreator.Id.Hex() + "_karma_points",
				Participants:    []bson.ObjectId{postCreator.Id},
				DisplayTemplate: constants.NotificationTemplate["TRANSACTIONAL"],
				EntityGroupId:   post.Id.Hex(),
				ActionId:        post.Id,
				ActionType:      "comment",
				Purpose:         constants.NotificationPurpose["KARMA_POINTS"],
				Entities:        entities,
				Message:         msgStr,
				MessageMeta:     messageMeta,
				MessageHtml:     htmlMsgStr,
				DeepLink:        deepLink,
			})

			icon := mongo.ExtractThumbNailFromPost(post)

			if icon == "" {
				icon = vote.Created.Avatar
			}

			msg := firebase.ManchMessage{
				Title:    title,
				Message:  msgStr,
				Icon:     icon,
				DeepLink: deepLink,
				Id:       notification.NId,
			}

			tokens := mongo.GetTokensByProfiles([]bson.ObjectId{post.Created.ProfileId})
			fmt.Printf("\nGCM Message %+v\n", msg)
			if tokens != nil {
				for _, token := range tokens {
					go firebase.SendMessage(msg, token.Token, notification)
				}
			} else {
				fmt.Printf("No token")
			}
			// update post
			mongo.UpdateOnePostsByQuery(bson.M{"_id": post.Id}, bson.M{
				"$addToSet": bson.M{"notified_karma_buckets": karmaNextBucket},
			})
			// delete lock
			redis.ReleaseLock(key)
		}
	}

	// schedule follow
	if (vote.Created.UserType == "bot" && post.Created.UserType != "bot") ||
		(vote.Created.UserType != "bot" && post.Created.UserType == "bot") {
		// bot follow user
		var profileId, resourceId bson.ObjectId
		if vote.Created.UserType == "bot" {
			profileId = vote.Created.ProfileId
			resourceId = post.Created.ProfileId
		} else {
			profileId = post.Created.ProfileId
			resourceId = vote.Created.ProfileId
		}
		randomNumber := utils.Random(0, 100)
		if randomNumber > 40 {
			t := time.Now().Add(time.Duration(utils.Random(1, 24)) * time.Hour)
			followSchedule := mongo.CreateFollowSchedule(t, profileId, resourceId)
			mongo.AddFollowSchedule(followSchedule)
		}
	}

	notification := mongo.NotificationModel{
		Receiver:        postCreator.Id,
		Identifier:      post.Id.Hex() + "_vote",
		Participants:    []bson.ObjectId{vote.Created.ProfileId},
		DisplayTemplate: constants.NotificationTemplate["TRANSACTIONAL"],
		EntityGroupId:   post.Id.Hex(),
		ActionId:        vote.Id,
		ActionType:      "vote",
		Purpose:         constants.NotificationPurpose["VOTE"],
		Entities:        entities,
		NUUID:           "",
	}

	postTitle := utils.TruncateTitle(post.Title, 4)
	data := i18n.DataModel{
		Name:  vote.Created.Name,
		Post:  postTitle,
		Count: count,
	}

	var msgStr string
	var templateName string
	if count == 0 {
		templateName = "post_like_one"
	} else if count == 1 {
		templateName = "post_like_two"
		// get other upvoter and udpate data.Name2
		votes := mongo.GetAllVoteByQuery(bson.M{
			"resource":           post.Id,
			"created.profile_id": bson.M{"$nin": []bson.ObjectId{vote.Created.ProfileId}},
		})
		fmt.Printf("all upvotes for template post_like_two \n %+v\n", votes)
		if len(votes) > 0 {
			profileId := votes[0].Created.ProfileId
			profile := mongo.GetProfileById(profileId)
			data.Name2 = profile.Name
			notification.Participants = append(notification.Participants, profileId)
		}
	} else if count == 2 {
		templateName = "post_like_three"
		// get other two upvoters and update data.Name2 and data.Name3
		votes := mongo.GetAllVoteByQuery(bson.M{
			"resource":           post.Id,
			"created.profile_id": bson.M{"$nin": []bson.ObjectId{vote.Created.ProfileId}},
		})
		if len(votes) > 1 {
			profileId1 := votes[0].Created.ProfileId
			profile1 := mongo.GetProfileById(profileId1)
			data.Name2 = profile1.Name
			notification.Participants = append(notification.Participants, profileId1)
			profileId2 := votes[1].Created.ProfileId
			profile2 := mongo.GetProfileById(profileId2)
			data.Name3 = profile2.Name
			notification.Participants = append(notification.Participants, profileId2)
		}
	} else {
		templateName = "post_like_multi"
	}

	msgStr = i18n.GetString(postCreator.Language, templateName, data)
	htmlMsgStr := i18n.GetHtmlString(postCreator.Language, templateName, data)
	msgStr = strings.Replace(msgStr, "\"\" ", "", 1)
	title := i18n.GetAppTitle(postCreator.Language)

	messageMeta := mongo.MessageMeta{
		TemplateName: templateName,
		Template:     i18n.Strings[postCreator.Language][templateName],
		Data:         data,
	}
	// update notification message
	deepLink := "manch://posts/" + post.Id.Hex()

	notification.Message = msgStr
	notification.MessageMeta = messageMeta
	notification.MessageHtml = htmlMsgStr
	notification.DeepLink = deepLink

	notification = mongo.CreateNotification(notification)

	icon := mongo.ExtractThumbNailFromPost(post)

	if icon == "" {
		icon = vote.Created.Avatar
	}

	msg := firebase.ManchMessage{
		Title:    title,
		Message:  msgStr,
		Icon:     icon,
		DeepLink: deepLink,
		Id:       notification.NId,
	}

	// upvote notification
	upvoteBucket := []int{5, 25, 50}
	nextBucket := getNextBucket(post.NotifiedVoteBuckets, upvoteBucket)
	if isValidBucket(nextBucket, count, post.NotifiedVoteBuckets) {
		// send notification
		bucketStr := strconv.Itoa(nextBucket)
		key := post.Id.Hex() + ":" + bucketStr
		result, err := redis.AcquireLock(key, bucketStr, 60)
		if err != nil || result == 0 {
			return
		}
		// send notification
		tokens := mongo.GetTokensByProfiles([]bson.ObjectId{post.Created.ProfileId})
		fmt.Printf("\nGCM Message %+v\n", msg)
		if tokens != nil {
			for _, token := range tokens {
				go firebase.SendMessage(msg, token.Token, notification)
			}
		} else {
			fmt.Printf("No token")
		}
		// update post
		mongo.UpdateOnePostsByQuery(bson.M{"_id": post.Id}, bson.M{
			"$addToSet": bson.M{"notified_vote_buckets": nextBucket},
		})
		// delete lock
		redis.ReleaseLock(key)
	}

	fmt.Printf("Processed a vote on subject %s! with vote Id %s\n", subj, v.Id)

}

func getNextBucket(archivedBucket []int, upvoteBucket []int) int {
	n := len(archivedBucket)
	if n == 0 {
		return upvoteBucket[0]
	}
	lastBucket := archivedBucket[len(archivedBucket)-1]
	var nextBucket int
	if lastBucket >= 50 {
		nextBucket = lastBucket + 50
		return nextBucket
	}

	for index, bucket := range upvoteBucket {
		if bucket == lastBucket {
			nextBucket = upvoteBucket[index+1]
		}
	}
	return nextBucket
}

func isValidBucket(nextBucket, upVotes int, archivedBucket []int) bool {
	if upVotes >= nextBucket && !utils.IncludesInt(archivedBucket, nextBucket) {
		return true
	}
	return false
}
