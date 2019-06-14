package callbacks

import (
	"fmt"
	"notification-service/pkg/constants"
	"notification-service/pkg/firebase"
	"notification-service/pkg/i18n"
	"notification-service/pkg/mongo"
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

	// // create postCreator's coin
	// mongo.CreateUserCoin(mongo.UserCoinsModel{
	// 	ProfileId:   post.Created.ProfileId,
	// 	CoinsEarned: 1,
	// 	Action:      "vote",
	// })

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

	// schedule vote for the user likes
	if vote.Created.UserType != "bot" {
		userNo := mongo.CountVoteByQuery(bson.M{
			"resource":     post.Id,
			"created.type": bson.M{"$ne": "bot"},
		})
		botProfiles := mongo.GetBotProfileByBucketId(userNo - 1)
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

	postCreator := mongo.GetProfileById(post.Created.ProfileId)
	upVotes := post.CoinsEarned
	// notification for karma points
	if upVotes != 0 && upVotes%50 == 0 {
		templateName := "post_karma_points"
		data := i18n.DataModel{
			Count: upVotes,
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

	count := post.UpVotes
	fmt.Println("post upvotes:", count)

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
	if count > 25 {
		msgStr = "❤️ " + msgStr
	}
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

	upvoteBucket := []int{1, 5, 25, 50, 75, 100}

	if utils.Contains(upvoteBucket, count+1) || (((count + 1) % 50) == 0) {
		tokens := mongo.GetTokensByProfiles([]bson.ObjectId{post.Created.ProfileId})
		fmt.Printf("\nGCM Message %+v\n", msg)
		if tokens != nil {
			for _, token := range tokens {
				go firebase.SendMessage(msg, token.Token, notification)
			}
		} else {
			fmt.Printf("No token")
		}
	}

	fmt.Printf("Processed a vote on subject %s! with vote Id %s\n", subj, v.Id)

}
