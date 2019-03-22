package callbacks

import (
	"fmt"
	"notification-service/pkg/constants"
	"notification-service/pkg/firebase"
	"notification-service/pkg/i18n"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"
	"strings"

	"github.com/globalsign/mgo/bson"
)

func UserFollowApprovedCB(subj, reply string, uf *subscribers.Subscription) {
	// update unique user remove (resource, type).uniquUsers set to profile id
	fmt.Printf("Received a User follow Approved subject %s! with user follow %+v\n", subj, uf)

	userProfile := mongo.GetProfileById(bson.ObjectIdHex(uf.ProfileId))

	community := mongo.GetCommunityById(uf.Resource)

	entities := []mongo.Entity{
		{
			EntityId:   community.Id,
			EntityType: "community",
		},
	}
	data := i18n.DataModel{
		Name:      userProfile.Name,
		Community: community.Name,
	}

	templateName := "join_manch_approved"
	deepLink := ""
	msgStr := i18n.GetString(userProfile.Language, templateName, data)
	htmlMsgStr := i18n.GetHtmlString(userProfile.Language, templateName, data)
	msgStr = strings.Replace(msgStr, "\"\" ", "", 1)
	title := i18n.GetAppTitle(userProfile.Language)

	messageMeta := mongo.MessageMeta{
		TemplateName: templateName,
		Template:     i18n.Strings[userProfile.Language][templateName],
		Data:         data,
	}
	purpose := constants.NotificationPurpose["JOIN_MANCH_APPROVED"]
	notification := mongo.CreateNotification(mongo.NotificationModel{
		Receiver:        userProfile.Id,
		Identifier:      userProfile.Id.Hex() + purpose,
		Participants:    []bson.ObjectId{userProfile.Id},
		DisplayTemplate: constants.NotificationTemplate["TRANSACTIONAL"],
		EntityGroupId:   community.Id.Hex(),
		ActionId:        community.Id,
		ActionType:      "userfollow",
		Purpose:         purpose,
		Entities:        entities,
		Message:         msgStr,
		MessageHtml:     htmlMsgStr,
		DeepLink:        deepLink,
		MessageMeta:     messageMeta,
	})

	tokens := mongo.GetTokensByProfiles([]bson.ObjectId{userProfile.Id})

	msg := firebase.ManchMessage{
		Title:    title,
		Message:  msgStr,
		DeepLink: deepLink,
		Id:       notification.NId,
	}

	fmt.Printf("\nGCM Message %+v\n", msg)
	if tokens != nil {
		for _, token := range tokens {
			go firebase.SendMessage(msg, token.Token, notification)
		}
	} else {
		fmt.Printf("No token\n")
	}

}
