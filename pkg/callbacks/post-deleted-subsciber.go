package callbacks

import (
	"fmt"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"

	"github.com/globalsign/mgo/bson"
)

func PostDeletedSubscriberCB(subj, reply string, p *subscribers.Post) {
	fmt.Printf("Received a post on subject %s! with Post %+v\n", subj, p)
	// remove vote schedule
	mongo.RemoveVoteScheduleByResource(bson.ObjectIdHex(p.Id))
	// remove share schedule
	mongo.RemoveShareScheduleByResource(bson.ObjectIdHex(p.Id))
	// remove comment schedule
	mongo.RemoveCommentScheduleByPostId(bson.ObjectIdHex(p.Id))
	fmt.Printf("Processed a post on subject %s! with Post ID %s\n", subj, p.Id)
}
