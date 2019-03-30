package callbacks

import (
	"fmt"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"
	"time"

	"github.com/globalsign/mgo/bson"
)

func UserStreakCB(subj, reply string, userStreak *subscribers.UserStreak) {
	fmt.Printf("Received a user streak on subject %s! with Value %+v\n", subj, userStreak)

	if userStreak.CurrentStreak.StreakLength == 1 {
		badge := mongo.Badge{
			ResourceName: "1 day streak",
			Icon:         "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/badges/1_day_streak.jpg",
		}
		milestone := mongo.Milestone{
			MileStoneId: "0",
			Name:        "1 day streak",
			Badge:       badge,
			Value:       1,
			Type:        "streak",
			Date:        time.Now(),
		}

		query := bson.M{
			"profiles.0._id":                    bson.ObjectIdHex(userStreak.ProfileId),
			"profiles.0.achieved_milestones.id": bson.M{"$ne": "0"},
		}

		update := bson.M{
			"profiles.$.current_badge": badge,
			"$push":                    bson.M{"profiles.$.achieved_milestones": milestone},
		}
		mongo.UpdateUser(query, update)
	}

	if userStreak.CurrentStreak.StreakLength == 7 {

	}

	if userStreak.CurrentStreak.StreakLength == 30 {

	}

	if userStreak.CurrentStreak.StreakLength == 100 {

	}
	fmt.Printf("Processed a user streak on subject %s! with profile Id %s\n", subj, userStreak.ProfileId)
}
