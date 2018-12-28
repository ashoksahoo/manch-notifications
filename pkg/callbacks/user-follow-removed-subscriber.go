package callbacks

import (
	"fmt"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"

	"github.com/globalsign/mgo/bson"
)

func UserFollowRemovedSubscriberCB(subj, reply string, uf *subscribers.Subscription) {
	// update unique user remove (resource, type).uniquUsers set to profile id
	fmt.Printf("Received a User follow on subject %s! with user follow %+v\n", subj, uf)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in subscribers.UserFollowSubscriber", uf)
		}
	}()

	if uf.ResourceType != "user" {
		fmt.Println("Not a user resource follows")
		return
	}

	follower := mongo.GetProfileById(bson.ObjectIdHex(uf.ProfileId))
	// fmt.Printf("\nfollower %+v\n", follower)
	followsTo := mongo.GetProfileById(bson.ObjectIdHex(uf.Resource))
	mongo.RemoveNotificationUser(followsTo.Id, "follows", follower.Id)
}
