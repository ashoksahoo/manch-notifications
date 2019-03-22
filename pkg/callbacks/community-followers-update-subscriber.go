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

func CommunityFollowersUpdateCB(subj, reply string, C *subscribers.Community) {
	fmt.Printf("Received a Community followers update on subject %s! with vote %+v\n", subj, C)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in subscribers.CommunityFollowersUpdateCB", r)
		}
	}()

	if C.FollowersCount != 10 || C.FollowersCount != 100 {
		return
	}

	community := mongo.GetCommunityById(C.Id)

	if community.Type == "m_manch" {
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
		if C.FollowersCount == 10 {
			templateName = "manch_10_members_title"
			templateText = "manch_10_members_text"
			purpose = constants.NotificationPurpose["MANCH_10_MEMBERS"]
		} else if C.FollowersCount == 100 {
			templateName = "manch_100_members"
			purpose = constants.NotificationPurpose["MANCH_100_MEMBERS"]
		}
		deepLink := ""
		for _, adminProfile := range adminProfiles {
			var htmlMsgStr, msgStr, title string
			if C.FollowersCount == 10 {
				msgStr = i18n.GetString(adminProfile.Language, templateText, data)
				htmlMsgStr = i18n.GetHtmlString(adminProfile.Language, templateName, data)
				title = i18n.GetString(adminProfile.Language, templateName, data)
			} else {
				msgStr = i18n.GetString(adminProfile.Language, templateName, data)
				htmlMsgStr = i18n.GetHtmlString(adminProfile.Language, templateName, data)
				title = i18n.GetAppTitle(adminProfile.Language)
			}
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

	fmt.Printf("Processed a community followers update on subject %s! with vote Id %s\n", subj, C.Id)

}
