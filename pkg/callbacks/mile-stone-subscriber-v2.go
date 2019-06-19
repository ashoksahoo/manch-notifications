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

func MileStoneSubscriberCBV2(subj, reply string, m *subscribers.MileStone) {
	fmt.Printf("Received a new milestone on subject %s! with Value %+v\n", subj, m)

	profile := mongo.GetProfileById(bson.ObjectIdHex(m.ProfileId))

	if m.MileStone == constants.MileStones["1000_COIN_MILESTONE"] {
		titleTemplate := "1000_milestone_title"
		textTemplate := "1000_milestone_text"
		data := i18n.DataModel{
			Name: profile.Name,
		}
		msgStrTitle := i18n.GetString(profile.Language, titleTemplate, data)
		msgStrText := i18n.GetString(profile.Language, textTemplate, data)
		htmlMsgStr := i18n.GetString(profile.Language, titleTemplate, data)
		bigPictureTemplateName := "1000_milestone_image"
		bigPicture := i18n.GetString(profile.Language, bigPictureTemplateName, data)

		messageMeta := mongo.MessageMeta{
			TemplateName: titleTemplate,
			Template:     i18n.Strings[profile.Language][titleTemplate],
			Data:         data,
		}
		deepLink := "manch://profile/" + profile.Id.Hex()

		notification := mongo.CreateNotification(mongo.NotificationModel{
			Receiver:        profile.Id,
			Identifier:      profile.Id.Hex() + "1000_coin_milestone",
			Participants:    []bson.ObjectId{profile.Id},
			DisplayTemplate: constants.NotificationTemplate["TRANSACTIONAL"],
			EntityGroupId:   profile.Id.Hex(),
			ActionId:        profile.Id,
			ActionType:      "coin_milestone",
			Purpose:         constants.NotificationPurpose["1000_COIN_MILESTONE"],
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
			Icon:       profile.Avatar,
			Id:         notification.NId,
		}
		tokens := mongo.GetTokensByProfiles([]bson.ObjectId{profile.Id})
		fmt.Printf("\nGCM Message %+v\n", msg)
		if tokens != nil {
			for _, token := range tokens {
				go firebase.SendMessage(msg, token.Token, notification)
			}
		} else {
			fmt.Printf("No token")
		}
	}

	if m.MileStone == constants.MileStones["10000_COIN_MILESTONE"] {
		titleTemplate := "10000_milestone_title"
		textTemplate := "10000_milestone_text"
		data := i18n.DataModel{
			Name: profile.Name,
		}
		msgStrTitle := i18n.GetString(profile.Language, titleTemplate, data)
		msgStrText := i18n.GetString(profile.Language, textTemplate, data)
		htmlMsgStr := i18n.GetString(profile.Language, titleTemplate, data)
		bigPictureTemplateName := "10000_milestone_image"
		bigPicture := i18n.GetString(profile.Language, bigPictureTemplateName, data)

		messageMeta := mongo.MessageMeta{
			TemplateName: titleTemplate,
			Template:     i18n.Strings[profile.Language][titleTemplate],
			Data:         data,
		}
		deepLink := "manch://profile/" + profile.Id.Hex()

		notification := mongo.CreateNotification(mongo.NotificationModel{
			Receiver:        profile.Id,
			Identifier:      profile.Id.Hex() + "10000_coin_milestone",
			Participants:    []bson.ObjectId{profile.Id},
			DisplayTemplate: constants.NotificationTemplate["TRANSACTIONAL"],
			EntityGroupId:   profile.Id.Hex(),
			ActionId:        profile.Id,
			ActionType:      "coin_milestone",
			Purpose:         constants.NotificationPurpose["10000_COIN_MILESTONE"],
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
			Icon:       profile.Avatar,
			Id:         notification.NId,
		}
		tokens := mongo.GetTokensByProfiles([]bson.ObjectId{profile.Id})
		fmt.Printf("\nGCM Message %+v\n", msg)
		if tokens != nil {
			for _, token := range tokens {
				go firebase.SendMessage(msg, token.Token, notification)
			}
		} else {
			fmt.Printf("No token")
		}
	}

	if m.MileStone == constants.MileStones["25000_COIN_MILESTONE"] {
		titleTemplate := "25000_milestone_title"
		textTemplate := "25000_milestone_text"
		data := i18n.DataModel{
			Name: profile.Name,
		}
		msgStrTitle := i18n.GetString(profile.Language, titleTemplate, data)
		msgStrText := i18n.GetString(profile.Language, textTemplate, data)
		htmlMsgStr := i18n.GetString(profile.Language, titleTemplate, data)
		bigPictureTemplateName := "25000_milestone_image"
		bigPicture := i18n.GetString(profile.Language, bigPictureTemplateName, data)

		messageMeta := mongo.MessageMeta{
			TemplateName: titleTemplate,
			Template:     i18n.Strings[profile.Language][titleTemplate],
			Data:         data,
		}
		deepLink := "manch://profile/" + profile.Id.Hex()

		notification := mongo.CreateNotification(mongo.NotificationModel{
			Receiver:        profile.Id,
			Identifier:      profile.Id.Hex() + "25000_coin_milestone",
			Participants:    []bson.ObjectId{profile.Id},
			DisplayTemplate: constants.NotificationTemplate["TRANSACTIONAL"],
			EntityGroupId:   profile.Id.Hex(),
			ActionId:        profile.Id,
			ActionType:      "coin_milestone",
			Purpose:         constants.NotificationPurpose["25000_COIN_MILESTONE"],
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
			Icon:       profile.Avatar,
			Id:         notification.NId,
		}
		tokens := mongo.GetTokensByProfiles([]bson.ObjectId{profile.Id})
		fmt.Printf("\nGCM Message %+v\n", msg)
		if tokens != nil {
			for _, token := range tokens {
				go firebase.SendMessage(msg, token.Token, notification)
			}
		} else {
			fmt.Printf("No token")
		}
	}

	if m.MileStone == constants.MileStones["100000_COIN_MILESTONE"] {
		titleTemplate := "100000_milestone_title"
		textTemplate := "100000_milestone_text"
		data := i18n.DataModel{
			Name: profile.Name,
		}
		msgStrTitle := i18n.GetString(profile.Language, titleTemplate, data)
		msgStrText := i18n.GetString(profile.Language, textTemplate, data)
		htmlMsgStr := i18n.GetString(profile.Language, titleTemplate, data)
		bigPictureTemplateName := "100000_milestone_image"
		bigPicture := i18n.GetString(profile.Language, bigPictureTemplateName, data)

		messageMeta := mongo.MessageMeta{
			TemplateName: titleTemplate,
			Template:     i18n.Strings[profile.Language][titleTemplate],
			Data:         data,
		}
		deepLink := "manch://profile/" + profile.Id.Hex()

		notification := mongo.CreateNotification(mongo.NotificationModel{
			Receiver:        profile.Id,
			Identifier:      profile.Id.Hex() + "100000_coin_milestone",
			Participants:    []bson.ObjectId{profile.Id},
			DisplayTemplate: constants.NotificationTemplate["TRANSACTIONAL"],
			EntityGroupId:   profile.Id.Hex(),
			ActionId:        profile.Id,
			ActionType:      "coin_milestone",
			Purpose:         constants.NotificationPurpose["100000_COIN_MILESTONE"],
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
			Icon:       profile.Avatar,
			Id:         notification.NId,
		}
		tokens := mongo.GetTokensByProfiles([]bson.ObjectId{profile.Id})
		fmt.Printf("\nGCM Message %+v\n", msg)
		if tokens != nil {
			for _, token := range tokens {
				go firebase.SendMessage(msg, token.Token, notification)
			}
		} else {
			fmt.Printf("No token")
		}
	}

	fmt.Printf("Processed a new milestone on subject %s! with Id %s\n", subj, m.ProfileId)
}
