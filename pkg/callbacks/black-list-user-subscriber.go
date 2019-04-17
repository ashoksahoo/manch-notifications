package callbacks

import (
	"fmt"
	"notification-service/pkg/constants"
	"notification-service/pkg/firebase"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"
	"notification-service/pkg/utils"

	"github.com/globalsign/mgo/bson"
)

func BlackListUserSubscriberCB(subj, reply string, blacklist *subscribers.BlackListProfile) {
	fmt.Printf("Received a blacklist update on subject %s! with Info %+v\n", subj, blacklist)

	status := blacklist.Status
	profileId := blacklist.ProfileId
	reason := blacklist.Reason
	entities := []mongo.Entity{
		{
			EntityId:   bson.ObjectIdHex(profileId),
			EntityType: "profile",
		},
	}
	notification := mongo.CreateNotification(mongo.NotificationModel{
		Receiver:        bson.ObjectIdHex(profileId),
		Identifier:      profileId + "_user_blocked",
		Participants:    []bson.ObjectId{bson.ObjectIdHex(profileId)},
		DisplayTemplate: "transactional",
		ActionId:        bson.ObjectIdHex(profileId),
		ActionType:      "profile",
		Purpose:         constants.NotificationPurpose["USER_BLOCKED"],
		Entities:        entities,
		PushType:        "manch:D",
		NUUID:           "",
	})

	blockedStatus := map[string]string{}
	if status == constants.BlackListStatus["BLOCKED"] {
		blockedStatus["status"] = constants.BlackListStatus["BLOCKED"]
		blockedStatus["blocked_on"] = utils.ISOFormat(blacklist.BlockedOn)
		blockedStatus["blocked_till"] = utils.ISOFormat(blacklist.BlockedTill)
	} else if status == constants.BlackListStatus["WARNING"] {
		notification.Purpose = constants.NotificationPurpose["USER_WARNED"]
		blockedStatus["status"] = constants.BlackListStatus["WARNING"]
		blockedStatus["last_warned_on"] = utils.ISOFormat(blacklist.LastWarnedOn)
		notification.Identifier = profileId + "_user_warned"
	} else if status == constants.BlackListStatus["UN_BLOCKED"] {
		notification.Purpose = constants.NotificationPurpose["USER_UNBLOCKED"]
		blockedStatus["status"] = constants.BlackListStatus["UN_BLOCKED"]
		notification.Identifier = profileId + "_user_unblocked"
	}

	tokens := mongo.GetTokensByProfiles([]bson.ObjectId{bson.ObjectIdHex(profileId)})
	msg := firebase.ManchMessage{
		Title:     "",
		Message:   "",
		Namespace: "manch:D",
		Id:        notification.NId,
		Reason:    reason,
	}

	if blockedStatus["status"] == constants.BlackListStatus["WARNING"] {
		msg.LastWarned = blockedStatus["last_warned_on"]
	} else if blockedStatus["status"] == constants.BlackListStatus["BLOCKED"] {
		msg.BlockedTill = blockedStatus["blocked_till"]
		msg.BlockedOn = blockedStatus["blocked_on"]
	}
	msg.Status = blockedStatus["status"]

	fmt.Printf("notification is %+v\n\n\n", notification)
	fmt.Printf("manch message is %+v\n\n\n", msg)

	if tokens != nil {
		for _, token := range tokens {
			go firebase.SendMessage(msg, token.Token, notification)
		}
	} else {
		fmt.Printf("No token")
	}
	fmt.Printf("Processed a blacklist udpate on subject %s! with ProfileId %s\n", subj, blacklist.ProfileId)
}
