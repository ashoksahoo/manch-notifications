package callbacks

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"
	"notification-service/pkg/utils"
	"time"
)

func ShareSubscriberCB(subj, reply string, share *subscribers.Share) {
	fmt.Printf("Received Share on Subject%s! with Share %+v\n", subj, share)

	if share.ResourceType != "post" {
		return
	}
	// count todayskey and purpose
	dayKey, _, _, _ := utils.GetCurrentDateKeys()
	err, _ := mongo.GetUserCoinsByQuery(bson.M{
		"profile_id": bson.ObjectIdHex(share.ProfileId),
		"action":     "share",
		"day_key":    dayKey,
	})
	fmt.Println("Error is", err)
	if err == nil {
		return
	}

	// count today's share of this profile id
	n := mongo.CountShareByQuery(bson.M{
		"profile_id":   bson.ObjectIdHex(share.ProfileId),
		"created.date": bson.M{"$gte": utils.GetStartOfDay(time.Now()), "$lt": utils.GetEndOfDay(time.Now())},
	})

	fmt.Println("Count of Share is ", n)
	if n >= 10 {
		mongo.CreateUserCoin(mongo.UserCoinsModel{
			ProfileId:   bson.ObjectIdHex(share.ProfileId),
			Action:      "share",
			CoinsEarned: 50,
		})
	}

	fmt.Printf("Processed Share on subject %s! with Share ID %s\n", subj, share.Id)

}
