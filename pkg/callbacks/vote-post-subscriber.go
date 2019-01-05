package callbacks

import (
	"fmt"
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

	post := mongo.GetPostById(v.Resource)
	vote := post.GetVote(v.Id)
	if vote.Created.ProfileId == post.Created.ProfileId {
		//Self Vote
		fmt.Println("Self Vote")
		return
	}
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
			EntityId: post.Id,
			EntityType: "post",
		},
		{
			EntityId: vote.Id,
			EntityType: "vote",
		},
	}
	notification := mongo.CreateNotification(mongo.NotificationModel{
		Receiver:        postCreator.Id,
		Identifier:      post.Id.Hex() + "_vote",
		Participants:    []bson.ObjectId{vote.Created.ProfileId},
		DisplayTemplate: "transactional",
		EntityGroupId:   post.Id.Hex(),
		ActionId:        vote.Id,
		ActionType:      "vote",
		Purpose:         "vote",
		Entities:        entities,
		NUUID:           "",
	})

	postTitle := utils.TruncateTitle(post.Title, 4)
	data := i18n.DataModel{
		Name:  vote.Created.Name,
		Post:  postTitle,
		Count: post.UpVotes,
	}
	
	var msgStr string
	var templateName string
	if post.UpVotes > 1 {
		templateName = "post_like_multi"
	} else {
		templateName = "post_like_one"
	}

	msgStr = i18n.GetString(postCreator.Language, templateName, data)	
	msgStr = strings.Replace(msgStr, "\"\" ", "", 1)
	title := i18n.GetAppTitle(postCreator.Language)

	messageMeta := mongo.MessageMeta{
		Template: templateName,
		Data: data,
	}
	// update notification message
	mongo.UpdateNotification(bson.M{"_id": notification.Id}, bson.M{
		"message": msgStr,
		"message_meta": messageMeta,
	})

	msg := firebase.ManchMessage{
		Title:    title,
		Message:  msgStr,
		Icon:     mongo.ExtractThumbNailFromPost(post),
		DeepLink: "manch://posts/" + post.Id.Hex(),
		Id:       notification.NId,
	}

	fmt.Printf("\nGCM Message %+v\n", msg)
	if tokens != nil {
		for _, token := range tokens {
			go firebase.SendMessage(msg, token.Token, notification)
		}
	} else {
		fmt.Printf("No token")
	}
	fmt.Printf("Processed a vote on subject %s! with vote Id %s\n", subj, v.Id)

}
