package callbacks

import (
	"fmt"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"

	"github.com/globalsign/mgo/bson"
)

func LiveTopicsWinnerSubscriberCB(subj, reply string, W *subscribers.LiveTopicsWinner) {
	fmt.Printf("Received a live topics winners on subject %s! with Value %+v\n", subj, W)

	winners := W.Winners

	for _, winner := range winners {
		mongo.CreateUserCoin(mongo.UserCoinsModel{
			ProfileId: bson.ObjectIdHex(winner),
			CoinsEarned: 500,
			Action: "live-topics.winner",
		})
	}

	fmt.Printf("Processed a live topics winners on subject %s! with Id %s\n", subj, W.Id)
}
