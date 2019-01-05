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
	// notification := mongo.CreateNotification(followsTo.Id, "follows", "user", follower.Id)
	
	entities := []mongo.Entity{
		{
			EntityId: userFollow.Id,
			EntityType: "user_follow",
		},
	}
	notification := mongo.CreateNotification(mongo.NotificationModel{
		Receiver:        followsTo.Id,
		Identifier:      followsTo.Id.Hex() + "_follow",
		Participants:    []bson.ObjectId{follower.Id},
		DisplayTemplate: "transactional",
		EntityGroupId:   userFollow.Id.Hex(),
		ActionId:        userFollow.Id,
		ActionType:      "userfollow",
		Purpose:         "follow",
		Entities:        entities,
		NUUID:           "",
	})
	count := followsTo.FollowerCount - 1
	data := i18n.DataModel{
		Name:  follower.Name,
		Count: count,
	}
	var msgStr string
	var templateName string
	if count > 0 {
		templateName = "follow_user_multi"
	} else {
		templateName = "follow_user_one"
	}

	msgStr = i18n.GetString(followsTo.Language, templateName, data)	
	msgStr = strings.Replace(msgStr, "\"\" ", "", 1)
	title := i18n.GetAppTitle(followsTo.Language)

	messageMeta := mongo.MessageMeta{
		Template: templateName,
		Data: data,
	}
	// update notification message
	mongo.UpdateNotification(bson.M{"_id": notification.Id}, bson.M{
		"message": msgStr,
		"message_meta": messageMeta,
	})

	msg := firebase.ManchMessage{
		Title:    title,
		Message:  msgStr,
		DeepLink: "manch://profile/" + followsTo.Id.Hex(),
		Id:       notification.NId,
	}
	//firebase.SendMessage(msg, "frgp37gfvFg:APA91bHbnbfoX-bp3M_3k-ceD7E4fZ73fcmVL4b5DGB5cQn-fFEvfbj3aAI9g0wXozyApIb-6wGsJauf67auK1p3Ins5Ff7IXCN161fb5JJ5pfBnTZ4LEcRUatO6wimsbiS7EANoGDr4")
	fmt.Printf("\nGCM Message %+v\n", msg)
	if tokens != nil {
		for _, token := range tokens {
			go firebase.SendMessage(msg, token.Token, notification)
		}
	} else {
		fmt.Printf("No token\n")
	}

	fmt.Printf("Processed a User follow on subject %s! with user follow ID %s\n", subj, uf.Id)

}
