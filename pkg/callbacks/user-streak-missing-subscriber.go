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

func UserSTreakMissingCB(subj, reply string, userStreak *subscribers.UserStreak) {
	fmt.Printf("Received a user streak missing on subject %s! with Value %+v\n", subj, userStreak)
	profile := mongo.GetProfileById(bson.ObjectIdHex(userStreak.ProfileId))

	var msgStrTitle, msgStrText string
	var templateTitle, templateText string
	templateTitle = "streak_missing_title"
	templateText = "streak_missing_text"
	data := i18n.DataModel{
		Name:  profile.Name,
		Count: userStreak.CurrentStreak.StreakLength,
	}
	msgStrTitle = i18n.GetString(profile.Language, templateTitle, data)
	msgStrText = i18n.GetString(profile.Language, templateText, data)
	htmlMsgStr := i18n.GetHtmlString(profile.Language, templateTitle, data)

	bigPicture := i18n.GetString(profile.Language, "streak_missing_image", data)

	messageMeta := mongo.MessageMeta{
		TemplateName: templateTitle,
		Template:     i18n.Strings[profile.Language][templateTitle],
		Data:         data,
	}
	deepLink := "manch://profile/" + profile.Id.Hex()

	notification := mongo.CreateNotification(mongo.NotificationModel{
		Receiver:        profile.Id,
		Identifier:      profile.Id.Hex() + "_streak_missing_" + userStreak.CurrentStreak.EndDate.Format("20060102150405"),
		Participants:    []bson.ObjectId{profile.Id},
		DisplayTemplate: constants.NotificationTemplate["TRANSACTIONAL"],
		EntityGroupId:   profile.Id.Hex(),
		ActionId:        profile.Id,
		ActionType:      "streak_missing",
		Purpose:         constants.NotificationPurpose["STREAK_MISSING"],
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
