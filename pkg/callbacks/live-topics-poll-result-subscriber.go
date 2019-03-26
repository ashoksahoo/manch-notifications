package callbacks

import (
	"fmt"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"
	"time"

	"github.com/globalsign/mgo/bson"
)

func LiveTopicsPollResultCB(subj, reply string, data *subscribers.LiveTopicPoll) {
	fmt.Printf("Received a live topics poll result on subject %s! with Data %+v\n", subj, data)

	granularity := "ipl"
	loc, _ := time.LoadLocation("Asia/Kolkata")
	granularityStart := time.Date(2019, 3, 23, 20, 0, 0, 0, loc)
	granularityEnd := time.Date(2019, 5, 19, 11, 59, 59, 999, loc)
	key := data.ParticipantId + "_" + granularity

	var coins int
	if data.AnswerId == data.ResultId {
		coins = data.CoinsEarned
	} else {
		coins = data.CoinsLost
	}
	mongo.CreateUserLeaderBoard(
		key,
		bson.ObjectIdHex(data.ParticipantId),
		granularity,
		coins,
		granularityStart,
		granularityEnd,
	)

	fmt.Printf("Processed a live topics poll result on subject %s! with ID %s\n", subj, data.TopicId)
}
