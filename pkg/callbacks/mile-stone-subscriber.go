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

func MileStoneSubscriberCB(subj, reply string, m *subscribers.MileStone) {
	fmt.Printf("Received a new milestone on subject %s! with Value %+v\n", subj, m)

	profile := mongo.GetProfileById(bson.ObjectIdHex(m.ProfileId))
	if m.MileStone == constants.MileStones["100_COIN_MILESTONE"] {
		// send notification for this milestone

		if m.IsRefered {
			// update referer's coin
			fmt.Println("updating referar's coin", m.ReferedBy)
			mongo.UpdateProfileById(bson.ObjectIdHex(m.ReferedBy), bson.M{
				"$inc": bson.M{"profiles.$.total_coins": m.ReferalCoin},
			})
		}

		tokens := mongo.GetTokensByProfiles([]bson.ObjectId{profile.Id})
		var msgStr string
		var templateName string
		templateName = "100_coin_milestone"
		data := i18n.DataModel{
			Name: profile.Name,
		}
		msgStr = i18n.GetString(profile.Language, templateName, data)
		htmlMsgStr := i18n.GetHtmlString(profile.Language, templateName, data)
		title := i18n.GetAppTitle(profile.Language)

		messageMeta := mongo.MessageMeta{
			TemplateName: templateName,
			Template:     i18n.Strings[profile.Language][templateName],
			Data:         data,
		}
		deepLink := "manch://posts/"

		notification := mongo.CreateNotification(mongo.NotificationModel{
			Receiver:        profile.Id,
			Identifier:      profile.Id.Hex() + "_100_milestone",
			Participants:    []bson.ObjectId{profile.Id},
			DisplayTemplate: constants.NotificationTemplate["TRANSACTIONAL"],
			EntityGroupId:   profile.Id.Hex(),
			ActionId:        profile.Id,
			ActionType:      "milestone",
			Purpose:         constants.NotificationPurpose["100_COIN_MILESTONE"],
			Message:         msgStr,
			MessageMeta:     messageMeta,
			MessageHtml:     htmlMsgStr,
			DeepLink:        deepLink,
		})

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
			fmt.Printf("No token")
		}
	}

	if m.MileStone == constants.MileStones["1000_COIN_MILESTONE"] {

	}

	if m.MileStone == constants.MileStones["10000_COIN_MILESTONE"] {

	}

	if m.MileStone == constants.MileStones["25000_COIN_MILESTONE"] {

	}

	fmt.Printf("Processed a new milestone on subject %s! with Id %s\n", subj, m.ProfileId)
}
