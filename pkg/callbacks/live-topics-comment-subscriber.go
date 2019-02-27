package callbacks

import (
	"fmt"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"

	"github.com/globalsign/mgo/bson"
)

func LiveTopicsCommentSubscriberCB(subj, reply string, comment *subscribers.LiveTopicComment) {
	fmt.Printf("Received a live topics comment on subject %s! with Comment %+v\n", subj, comment)

	action := "live-topics.comment"
	if comment.IsReply {
		action = "live-topics.reply"
	}

	if comment.IsReply {
		mongo.CreateUserCoin(mongo.UserCoinsModel{
			ProfileId:   bson.ObjectIdHex(comment.CreatedBy),
			CoinsEarned: 4,
			Action:      action,
		})
	} else {
		mongo.CreateUserCoin(mongo.UserCoinsModel{
			ProfileId:   bson.ObjectIdHex(comment.CreatedBy),
			CoinsEarned: 2,
			Action:      action,
		})
		upVoteCoins := (2 * comment.UpVotes)
		replyCoins := (10 * comment.ReplyCount)
		mongo.CreateUserCoin(mongo.UserCoinsModel{
			ProfileId:   bson.ObjectIdHex(comment.CreatedBy),
			CoinsEarned: upVoteCoins,
			Action:      "live-topics.comment.upvotes",
		})
		mongo.CreateUserCoin(mongo.UserCoinsModel{
			ProfileId:   bson.ObjectIdHex(comment.CreatedBy),
			CoinsEarned: replyCoins,
			Action:      "live-topics.comment.replies",
		})
	}

	fmt.Printf("Processed a live topics comment on subject %s! with Comment %s\n", subj, comment.Id)
}
