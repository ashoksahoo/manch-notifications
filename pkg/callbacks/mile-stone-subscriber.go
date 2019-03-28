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

		tokens := mongo.GetTokensByProfiles([]bson.ObjectId{profile.Id})
		var msgStrTitle, msgStrText string
		var templateTitle, templateText string
		templateTitle = "100_coin_milestone_title"
		templateText = "100_coin_milestone_text"
		data := i18n.DataModel{
			Name: profile.Name,
		}
		msgStrTitle = i18n.GetString(profile.Language, templateTitle, data)
		msgStrText = i18n.GetString(profile.Language, templateText, data)
		htmlMsgStr := i18n.GetHtmlString(profile.Language, templateTitle, data)

		bigPicture := i18n.GetString(profile.Language, "100_coin_milestone_image", data)

		messageMeta := mongo.MessageMeta{
			TemplateName: templateTitle,
			Template:     i18n.Strings[profile.Language][templateTitle],
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
			Message:         msgStrTitle,
			MessageMeta:     messageMeta,
			MessageHtml:     htmlMsgStr,
			DeepLink:        deepLink,
		})

		msg := firebase.ManchMessage{
			Title:      msgStrTitle,
			Message:    msgStrText,
			DeepLink:   deepLink,
			BigPicture: bigPicture,
			Id:         notification.NId,
		}

		fmt.Printf("\nGCM Message %+v\n", msg)
		if tokens != nil {
			for _, token := range tokens {
				go firebase.SendMessage(msg, token.Token, notification)
			}
		} else {
			fmt.Printf("No token")
		}

		// update referrer's coin
		err, referralData := mongo.GetReferralsByQuery(bson.M{
			"profile_id":       m.ProfileId,
			"referring_params": bson.M{"$exists": true},
		})
		if err != nil {
			fmt.Println("error", err)
		} else {
			fmt.Printf("referal\n%+v", referralData)
			referredBy := referralData.ReferringParams["profile_id"].(string)
			mongo.UpdateProfileById(bson.ObjectIdHex(referredBy), bson.M{
				"$inc": bson.M{"profiles.$.total_coins": 100},
			})
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
