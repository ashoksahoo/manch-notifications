package callbacks

import (
	"fmt"
	"notification-service/pkg/subscribers"
	"notification-service/pkg/mongo"
	"github.com/globalsign/mgo/bson"
	"notification-service/pkg/i18n"
	"notification-service/pkg/firebase"
	"notification-service/pkg/constants"
)

func ReferralSubscriberCB(subj, reply string, referralData *subscribers.Referral) {
	fmt.Printf("Received a referral event on subject %s! with Post %+v\n", subj, referralData)

	if _, isErrorMessage := referralData.ReferringParams["error_message"]; isErrorMessage {
		return
	}

	referredByProfileId := referralData.ReferringParams["profile_id"].(string)

	mongo.CreateUserCoin(mongo.UserCoinsModel{
		ProfileId:   bson.ObjectIdHex(referredByProfileId),
		CoinsEarned: 100,
		Action:      "referral",
	})

	profile := mongo.GetProfileById(bson.ObjectIdHex(referralData.ProfileId))
	referrer := mongo.GetProfileById(bson.ObjectIdHex(referredByProfileId))
	tokens := mongo.GetTokensByProfiles([]bson.ObjectId{referrer.Id})
	var msgStrTitle, msgStrText string
	var templateTitle, templateText string
	templateTitle = "100_coin_referral_title"
	templateText = "100_coin_referral_text"
	data := i18n.DataModel{
		Name:  referrer.Name,
		Name2: profile.Name,
	}
	msgStrTitle = i18n.GetString(referrer.Language, templateTitle, data)
	msgStrText = i18n.GetString(referrer.Language, templateText, data)
	htmlMsgStr := i18n.GetHtmlString(referrer.Language, templateTitle, data)

	messageMeta := mongo.MessageMeta{
		TemplateName: templateTitle,
		Template:     i18n.Strings[referrer.Language][templateTitle],
		Data:         data,
	}
	deepLink := "manch://profile/" + referrer.Id.Hex()

	notification := mongo.CreateNotification(mongo.NotificationModel{
		Receiver:        referrer.Id,
		Identifier:      referrer.Id.Hex() + profile.Id.Hex() + "_100_referral",
		Participants:    []bson.ObjectId{referrer.Id},
		DisplayTemplate: constants.NotificationTemplate["TRANSACTIONAL"],
		EntityGroupId:   referrer.Id.Hex(),
		ActionId:        referrer.Id,
		ActionType:      "referral",
		Purpose:         constants.NotificationPurpose["100_COIN_REFERRAL"],
		Message:         msgStrTitle,
		MessageMeta:     messageMeta,
		MessageHtml:     htmlMsgStr,
		DeepLink:        deepLink,
	})

	msg := firebase.ManchMessage{
		Title:      msgStrTitle,
		Message:    msgStrText,
		DeepLink:   deepLink,
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

	fmt.Printf("Processed a referral event on subject %s! with Post Id %s\n", subj, referralData.Id)

}
