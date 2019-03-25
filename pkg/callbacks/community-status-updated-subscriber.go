package callbacks

import (
	"fmt"
	"notification-service/pkg/constants"
	"notification-service/pkg/firebase"
	"notification-service/pkg/i18n"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"

	"github.com/globalsign/mgo/bson"
)

func CommunityStatusUpdatedCB(subj, reply string, C *subscribers.Community) {
	fmt.Printf("Received a Community Status update on subject %s! with Community %+v\n", subj, C)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in subscribers.CommunityStatusUpdatedCB", r)
		}
	}()

	if C.Status == constants.CommunityStatus["ACTIVATED"] && C.Type == "m_manch" {
		community := mongo.GetCommunityById(C.Id)
		fmt.Println("community type is", community.Type)
		admins := community.Admins
		adminProfilesIds := []string{}

		for _, admin := range admins {
			adminProfilesIds = append(adminProfilesIds, admin.ProfileId.Hex())
		}

		adminProfiles := mongo.GetProfilesByIds(adminProfilesIds)

		entities := []mongo.Entity{
			{
				EntityId:   community.Id,
				EntityType: "community",
			},
		}

		data := i18n.DataModel{
			Community: community.Name,
		}

		var purpose, templateName, templateText string

		templateName = "manch_activation_title"
		templateText = "manch_activation_text"
		purpose = constants.NotificationPurpose["MANCH_ACTIVATION"]
		deepLink := "manch://manch/" + community.Id.Hex()
		for _, adminProfile := range adminProfiles {
			var htmlMsgStr, msgStr, title string
			msgStr = i18n.GetString(adminProfile.Language, templateText, data)
			htmlMsgStr = i18n.GetHtmlString(adminProfile.Language, templateName, data)
			title = i18n.GetString(adminProfile.Language, templateName, data)
			fmt.Println("message str is ", msgStr)
			fmt.Println("title is", title)
			messageMeta := mongo.MessageMeta{
				TemplateName: templateName,
				Template:     i18n.Strings[adminProfile.Language][templateName],
				Data:         data,
			}

			notification := mongo.CreateNotification(mongo.NotificationModel{
				Receiver:        adminProfile.Id,
				Identifier:      adminProfile.Id.Hex() + purpose,
				Participants:    []bson.ObjectId{adminProfile.Id},
				DisplayTemplate: constants.NotificationTemplate["TRANSACTIONAL"],
				EntityGroupId:   community.Id.Hex(),
				ActionId:        community.Id,
				ActionType:      "community",
				Purpose:         purpose,
				Entities:        entities,
				Message:         msgStr,
				MessageHtml:     htmlMsgStr,
				DeepLink:        deepLink,
				MessageMeta:     messageMeta,
			})

			tokens := mongo.GetTokensByProfiles([]bson.ObjectId{adminProfile.Id})

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

	}

	fmt.Printf("Processed a Community Status Update on subject %s! with Id %s\n", subj, C.Id)

}
