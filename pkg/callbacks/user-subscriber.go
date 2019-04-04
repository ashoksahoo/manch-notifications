package callbacks

import (
	"fmt"
	"math/rand"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"
	"notification-service/pkg/utils"
	"time"

	"github.com/globalsign/mgo/bson"
)

func UserSubscriberCB(subj, reply string, u *subscribers.User) {
	fmt.Printf("Received a New User on subject %s! with User %+v\n", subj, u)
	// create follow schedule for this user
	var resourceId bson.ObjectId

	// get user from db
	user := mongo.GetUserById(u.Id)
	userProfileId := user.Profiles[0].Id

	m, botProfilesHi := mongo.GetBotProfilesIds("hi")
	n, botProfilesTe := mongo.GetBotProfilesIds("te")
	n = m + n
	botProfilesIds := append(botProfilesHi, botProfilesTe...)
	// shuffle profiles
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(n, func(i, j int) { botProfilesIds[i], botProfilesIds[j] = botProfilesIds[j], botProfilesIds[i] })

	// set user to resource
	resourceId = userProfileId

	// 0-5th minute - +5 followes
	j := 0
	randomFollowers := utils.Random(1, 3)
	followers := randomFollowers
	t := utils.SplitTimeInRange(1, 5, randomFollowers, time.Minute)
	for k := 0; j < followers; j, k = j+1, k+1 {
		doc := mongo.CreateFollowSchedule(t[k], bson.ObjectIdHex(botProfilesIds[j]), resourceId)
		// fmt.Printf("saving doc:%+v\n", doc)
		mongo.AddFollowSchedule(doc)
	}

	// 5 minuts to 1 hours - +5
	randomFollowers = utils.Random(1, 3)
	t = utils.SplitTimeInRange(5, 1440, randomFollowers, time.Minute)
	followers += randomFollowers
	for k := 0; j < followers; j, k = j+1, k+1 {
		doc := mongo.CreateFollowSchedule(t[k], bson.ObjectIdHex(botProfilesIds[j]), resourceId)
		// fmt.Printf("saving doc:%+v\n", doc)
		mongo.AddFollowSchedule(doc)
	}

	// 1 Hr to 6Hr +5-10 followers
	randomFollowers = utils.Random(3, 5)
	t = utils.SplitTimeInRange(24, 48, randomFollowers, time.Hour)
	followers += randomFollowers
	for k := 0; j < followers; j, k = j+1, k+1 {
		doc := mongo.CreateFollowSchedule(t[k], bson.ObjectIdHex(botProfilesIds[j]), resourceId)
		// fmt.Printf("saving doc:%+v\n", doc)
		mongo.AddFollowSchedule(doc)
	}

	// 6 Hr to 24 Hr +5-10 followers
	randomFollowers = utils.Random(3, 5)
	t = utils.SplitTimeInRange(48, 72, randomFollowers, time.Hour)
	followers += randomFollowers
	for k := 0; j < followers; j, k = j+1, k+1 {
		doc := mongo.CreateFollowSchedule(t[k], bson.ObjectIdHex(botProfilesIds[j]), resourceId)
		// fmt.Printf("saving doc:%+v\n", doc)
		mongo.AddFollowSchedule(doc)
	}

	// // 1st to 3rd day +10-15 followers
	// randomFollowers = utils.Random(20, 30)
	// t = utils.SplitTimeInRange(24, 72, randomFollowers, time.Hour)
	// followers += randomFollowers
	// for k := 0; j < followers; j, k = j+1, k+1 {
	// 	doc := mongo.CreateFollowSchedule(t[k], bson.ObjectIdHex(botProfilesIds[j]), resourceId)
	// 	mongo.AddFollowSchedule(doc)
	// }

	// 3rd to 7th day +10-20 followers
	randomFollowers = utils.Random(6, 10)
	t = utils.SplitTimeInRange(72, 168, randomFollowers, time.Hour)
	followers += randomFollowers
	for k := 0; j < followers; j, k = j+1, k+1 {
		doc := mongo.CreateFollowSchedule(t[k], bson.ObjectIdHex(botProfilesIds[j]), resourceId)
		mongo.AddFollowSchedule(doc)
	}
	fmt.Println("total followers added:", followers)
	fmt.Printf("Processed a New User on subject %s! with User Id %s\n", subj, u.Id)
}
