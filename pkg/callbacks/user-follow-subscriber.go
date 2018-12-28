package callbacks

import (
	"fmt"
	"notification-service/pkg/firebase"
	"notification-service/pkg/i18n"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"
	"strings"

	"github.com/globalsign/mgo/bson"
)

func UserFollowSubscriberCB(subj, reply string, uf *subscribers.Subscription) {
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

	userFollow := mongo.GetUserFollowById(uf.Id)
	// fmt.Printf("\nuser follow %+v\n", userFollow)
	follower := mongo.GetProfileById(userFollow.ProfileId)
	// fmt.Printf("\nfollower %+v\n", follower)
	followsTo := mongo.GetProfileById(userFollow.ResourceId)
	// fmt.Printf("\nfollowsTo %+v\n", followsTo)
	tokens := mongo.GetTokensByProfiles([]bson.ObjectId{userFollow.ResourceId})
	notification := mongo.CreateNotification(followsTo.Id, "follows", "user", follower.Id)

	count := len(notification.UniqueUsers) - 1
	data := i18n.DataModel{
		Name:  follower.Name,
		Count: count,
	}
	var msgStr string

	if len(notification.UniqueUsers) > 1 {
		msgStr = i18n.GetString(followsTo.Language, "follow_user_multi", data)
	} else {
		msgStr = i18n.GetString(followsTo.Language, "follow_user_one", data)
	}

	title := i18n.GetAppTitle(followsTo.Language)
	msgStr = strings.Replace(msgStr, "\"\" ", "", 1)
	msg := firebase.ManchMessage{
		Title:    title,
		Message:  msgStr,
		DeepLink: "manch://profile/" + followsTo.Id.Hex(),
		Id:       notification.Identifier,
	}
	//firebase.SendMessage(msg, "frgp37gfvFg:APA91bHbnbfoX-bp3M_3k-ceD7E4fZ73fcmVL4b5DGB5cQn-fFEvfbj3aAI9g0wXozyApIb-6wGsJauf67auK1p3Ins5Ff7IXCN161fb5JJ5pfBnTZ4LEcRUatO6wimsbiS7EANoGDr4")
	fmt.Printf("\nGCM Message %+v\n", msg)
	if tokens != nil {
		for _, token := range tokens {
			go firebase.SendMessage(msg, token.Token)
		}
	} else {
		fmt.Printf("No token\n")
	}

	fmt.Printf("Processed a User follow on subject %s! with user follow ID %s\n", subj, uf.Id)

}
