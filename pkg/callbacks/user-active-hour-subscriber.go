package callbacks

import (
	"fmt"
	"notification-service/pkg/constants"
	"notification-service/pkg/firebase"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"

	"github.com/globalsign/mgo/bson"
)

func UserActiveHourCB(subj, reply string, u *subscribers.UserActiveHour) {
	fmt.Printf("Received a New User created on subject %s! with User %+v\n", subj, u)

	n := mongo.CountUserActiveHour(bson.M{
		"profile_id": bson.ObjectIdHex(u.ProfileId),
	})
	user := mongo.GetUserByProfileId(u.ProfileId)
	profile := user.Profiles[0]
	createdAt := user.CreatedAt
	diff := u.LastActiveHour.Sub(createdAt)

	fmt.Println("n", n)
	fmt.Println("profile.ratingnotifed", profile.RatingNotified)
	fmt.Println("diff", diff.Hours())

	if n > 2 && !profile.RatingNotified && diff.Hours() >= 0 {
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
			mongo.UpdateProfileById(profile.Id, bson.M{"$set": bson.M{"rating_notified": true}})
		} else {
			fmt.Printf("No token")
		}

	}

	fmt.Printf("Processed a New User created on subject %s! with User Id %s\n", subj, u.Id)
}
