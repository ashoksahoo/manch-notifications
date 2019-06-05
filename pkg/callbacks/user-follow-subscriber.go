package callbacks

import (
	"fmt"
	"notification-service/pkg/constants"
	"notification-service/pkg/firebase"
	"notification-service/pkg/i18n"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"
	"notification-service/pkg/utils"
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

	follower := mongo.GetProfileById(bson.ObjectIdHex(uf.ProfileId))

	userFollow := mongo.GetUserFollowById(uf.Id)
	// fmt.Printf("\nuser follow %+v\n", userFollow)
	followsTo := mongo.GetProfileById(userFollow.ResourceId)
	// fmt.Printf("\nfollowsTo %+v\n", followsTo)

	entities := []mongo.Entity{
		{
			EntityId:   userFollow.Id,
			EntityType: "user_follow",
		},
	}
	count := followsTo.FollowerCount
	data := i18n.DataModel{
		Name:  follower.Name,
		Count: count,
	}
	var msgStr string
	var templateName string
	if count > 1 {
		templateName = "follow_user_multi"
	} else {
		templateName = "follow_user_one"
	}

	msgStr = i18n.GetString(followsTo.Language, templateName, data)
	htmlMsgStr := i18n.GetHtmlString(followsTo.Language, templateName, data)
	msgStr = strings.Replace(msgStr, "\"\" ", "", 1)
	title := i18n.GetAppTitle(followsTo.Language)

	messageMeta := mongo.MessageMeta{
		TemplateName: templateName,
		Template:     i18n.Strings[followsTo.Language][templateName],
		Data:         data,
	}
	deepLink := "manch://profile/" + followsTo.Id.Hex()

	notification := mongo.CreateNotification(mongo.NotificationModel{
		Receiver:        followsTo.Id,
		Identifier:      followsTo.Id.Hex() + "_follow",
		Participants:    []bson.ObjectId{follower.Id},
		DisplayTemplate: constants.NotificationTemplate["TRANSACTIONAL"],
		EntityGroupId:   userFollow.Id.Hex(),
		ActionId:        userFollow.Id,
		ActionType:      "userfollow",
		Purpose:         constants.NotificationPurpose["USER_FOLLOW"],
		Entities:        entities,
		NUUID:           "",
		Message:         msgStr,
		MessageMeta:     messageMeta,
		MessageHtml:     htmlMsgStr,
		DeepLink:        deepLink,
	})

	icon := follower.Avatar
	msg := firebase.ManchMessage{
		Title:    title,
		Message:  msgStr,
		Icon:     icon,
		DeepLink: deepLink,
		Id:       notification.NId,
	}

	followBatch := []int{1, 2, 3, 5, 15, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	if utils.Contains(followBatch, count+1) || ((count+1)%50 == 0) {
		fmt.Printf("\nGCM Message %+v\n", msg)
		tokens := mongo.GetTokensByProfiles([]bson.ObjectId{userFollow.ResourceId})
		if tokens != nil {
			for _, token := range tokens {
				go firebase.SendMessage(msg, token.Token, notification)
			}
		} else {
			fmt.Printf("No token\n")
		}
	}
	fmt.Printf("Processed a User follow on subject %s! with user follow ID %s\n", subj, uf.Id)
}
