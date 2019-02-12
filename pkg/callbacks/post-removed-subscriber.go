package callbacks

import (
	"notification-service/pkg/constants"
	"fmt"
	"notification-service/pkg/firebase"
	"notification-service/pkg/i18n"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"
	"notification-service/pkg/utils"
	"strings"

	"github.com/globalsign/mgo/bson"
)

func PostRemovedSubscriberCB(subj, reply string, p *subscribers.Post) {
	fmt.Printf("Received a post on subject %s! with Post %+v\n", subj, p)
	err, post := mongo.GetPostById(p.Id)
	if err != nil {
		return
	}
	
	// Remove all vote and share scheduleds
	mongo.RemoveVoteScheduleByResource(post.Id)
	mongo.RemoveShareScheduleByResource(post.Id)

	return
	
	postCreator := mongo.GetProfileById(post.Created.ProfileId)

	tokens := mongo.GetTokensByProfiles([]bson.ObjectId{post.Created.ProfileId})

	entities := []mongo.Entity{
		{
			EntityId: post.Id,
			EntityType: "post",
		},
	}

	notification := mongo.CreateNotification(mongo.NotificationModel{
		Receiver:        postCreator.Id,
		Identifier:      post.Id.Hex() + "_remove",
		Participants:    []bson.ObjectId{postCreator.Id},
		DisplayTemplate: constants.NotificationTemplate["TRANSACTIONAL"],
		EntityGroupId:   post.Id.Hex(),
		ActionId:        post.Id,
		ActionType:      "post",
		Purpose:         constants.NotificationPurpose["POST_REMOVE"],
		Entities:        entities,
		NUUID:           "",
	})

	reason := post.IgnoreReason
	language := postCreator.Language
	deleteReason := i18n.DeleteReason[language][reason]

	postTitle := utils.TruncateTitle(post.Title, 4)
	data := i18n.DataModel{
		Name:         postCreator.Name,
		Post:         postTitle,
		DeleteReason: deleteReason,
	}
	var msgStr string
	templateName := "post_removed"
	msgStr = i18n.GetString(language, templateName, data)
	fmt.Println(msgStr)
	title := i18n.GetAppTitle(language)
	msgStr = strings.Replace(msgStr, "\"\" ", "", 1)

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
		DeepLink: "manch://posts/terms_" + post.Id.Hex(),
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
	fmt.Printf("Processed a post on subject %s! with Post ID %s\n", subj, p.Id)

}
