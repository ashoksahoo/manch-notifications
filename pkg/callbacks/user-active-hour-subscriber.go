package callbacks

import (
	"fmt"
	"notification-service/pkg/constants"
	"notification-service/pkg/firebase"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"
	"time"

	"github.com/globalsign/mgo/bson"
)

func UserActiveHourCB(subj, reply string, u *subscribers.UserActiveHour) {
	fmt.Printf("Received a New User Active Hour subject %s! with User %+v\n", subj, u)
	profile := mongo.GetProfileById(bson.ObjectIdHex(u.ProfileId))

	fmt.Println("profile.ratingnotifed", profile.RatingNotified)
	if profile.RatingNotified {
		return
	}

	currentTime := time.Now()
	threedaysAgo := currentTime.AddDate(0, 0, -3)
	fmt.Println("currentTime", currentTime)
	fmt.Println("three days ago", threedaysAgo)
	noOfSessions := mongo.CountUserActiveHour(bson.M{
		"profile_id": bson.ObjectIdHex(u.ProfileId),
		"createdAt": bson.M{"$gte": threedaysAgo},
	})
	fmt.Println("n", noOfSessions)
	if noOfSessions > 4 {
		// notify profile
		notification := mongo.CreateNotification(mongo.NotificationModel{
			Receiver:        profile.Id,
			Identifier:      profile.Id.Hex() + "_user_review",
			Participants:    []bson.ObjectId{profile.Id},
			DisplayTemplate: constants.NotificationTemplate["TRANSACTIONAL"],
			ActionId:        profile.Id,
			ActionType:      "rating_popup",
			Purpose:         constants.NotificationPurpose["USER_REVIEW"],
			PushType:        "manch:D",
		})

		tokens := mongo.GetTokensByProfiles([]bson.ObjectId{profile.Id})
		msg := firebase.ManchMessage{
			Title:     "",
			Message:   "",
			Namespace: "manch:D",
			Id:        notification.NId,
		}
		if tokens != nil {
			for _, token := range tokens {
				fmt.Println("successfully sent data message")
				go firebase.SendMessage(msg, token.Token, notification)
			}
			mongo.UpdateProfileById(profile.Id, bson.M{
				"$set": bson.M{"profiles.$.rating_notified": true, "profiles.$.rating_notified_at": currentTime},
			})
		} else {
			fmt.Printf("No token")
		}

	}

	fmt.Printf("Processed a New User Active Hour subject %s! with User Id %s\n", subj, u.Id)
}
