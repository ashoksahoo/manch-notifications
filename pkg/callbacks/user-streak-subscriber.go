package callbacks

import (
	"fmt"
	"notification-service/pkg/subscribers"
)

func UserStreakCB(subj, reply string, userStreak *subscribers.UserStreak) {
	fmt.Printf("Received a user streak on subject %s! with Value %+v\n", subj, userStreak)
	// if userStreak.CurrentStreak.StreakLength == 7 {
	// 	mongo.UpdateProfileById(bson.ObjectIdHex(userStreak.ProfileId), bson.M)
	// }
	fmt.Printf("Processed a user streak on subject %s! with profile Id %s\n", subj, userStreak.ProfileId)
}
