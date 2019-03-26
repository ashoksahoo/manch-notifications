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

	// mongo.CreateUserCoin(mongo.UserCoinsModel{
	// 	ProfileId:   bson.ObjectIdHex(comment.CreatedBy),
	// 	CoinsEarned: upVoteCoins,
	// 	Action:      "live-topics.comment.upvotes",
	// })
	granularity := "ipl"
	granularityStart := time.Now() // TODO: fix granularity start and end
	granularityEnd := time.Now()
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
