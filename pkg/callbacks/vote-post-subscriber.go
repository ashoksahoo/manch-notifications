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
	postCreator := mongo.GetProfileById(post.Created.ProfileId)
	tokens := mongo.GetTokensByProfiles([]bson.ObjectId{post.Created.ProfileId})
	// notification := mongo.CreateNotification(post.Id, "like", "post", vote.Created.ProfileId)

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
			"created.profile_id": bson.M{"$nin": []bson.ObjectId{vote.Created.ProfileId}},
		})
		fmt.Printf("all upvotes for template post_like_two \n %+v\n", votes)
		if len(votes) > 0 {
			fmt.Println("upvote length is greater than 0")
			profileId := votes[0].Created.ProfileId
			fmt.Println("profile id is", profileId)
			profile := mongo.GetProfileById(profileId)
			fmt.Println("profile is ", profile)
			data.Name2 = profile.Name
			fmt.Printf("udpated data is %+v", data)
			notification.Participants = append(notification.Participants, profileId)
			fmt.Println("after updating participants is", notification.Participants);
		}
	} else if count == 2 {
		templateName = "post_like_three"
		// get other two upvoters and update data.Name2 and data.Name3
		votes := mongo.GetAllVoteByQuery(bson.M{
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

	notification = mongo.CreateNotification(notification)
	fmt.Println("notification created")
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
	mongo.UpdateNotification(bson.M{"_id": notification.Id}, bson.M{
		"message":      msgStr,
		"message_meta": messageMeta,
		"message_html": htmlMsgStr,
		"deep_link":    deepLink,
	})

	msg := firebase.ManchMessage{
		Title:    title,
		Message:  msgStr,
		Icon:     mongo.ExtractThumbNailFromPost(post),
		DeepLink: deepLink,
		Id:       notification.NId,
	}

	upvoteNumbers := []int{1, 2, 5, 10, 25, 50, 75, 100}

	if utils.Contains(upvoteNumbers, count + 1) {
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
