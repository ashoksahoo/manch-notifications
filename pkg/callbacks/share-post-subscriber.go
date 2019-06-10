package callbacks

import (
	"fmt"
	"notification-service/pkg/constants"
	"strings"

	"notification-service/pkg/firebase"
	"notification-service/pkg/i18n"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"

	"github.com/globalsign/mgo/bson"
)

func SharePostSubscriberCB(subj, reply string, share *subscribers.SharePost) {

	fmt.Printf("Received a Post Share on Subject%s! with postShare %+v\n", subj, share)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in subscribers.SharePostSubscriber", share)
		}
	}()

	err, post := mongo.GetPostById(share.Id)
	if err != nil {
		return
	}

	profile := mongo.GetProfileById(bson.ObjectIdHex(share.ProfileId))

	postCreator := mongo.GetProfileById(post.Created.ProfileId)

	entities := []mongo.Entity{
		{
			EntityId:   post.Id,
			EntityType: "post",
		},
	}

	count := share.ShareCount
	data := i18n.DataModel{
		Name:  profile.Name,
		Count: count,
	}
	var msgStr, htmlMsgStr, templateName string
	if count == 1 {
		templateName = "share_post_one"
	} else {
		templateName = "share_post_multi"
	}

	msgStr = i18n.GetString(postCreator.Language, templateName, data)
	htmlMsgStr = i18n.GetHtmlString(postCreator.Language, templateName, data)
	msgStr = strings.Replace(msgStr, "\"\" ", "", 1)
	title := i18n.GetAppTitle(postCreator.Language)

	messageMeta := mongo.MessageMeta{
		TemplateName: templateName,
		Template:     i18n.Strings[postCreator.Language][templateName],
		Data:         data,
	}
	deepLink := "manch://posts/" + post.Id.Hex()
	// update notification message
	notification := mongo.CreateNotification(mongo.NotificationModel{
		Receiver:        postCreator.Id,
		Identifier:      post.Id.Hex() + "_share",
		Participants:    []bson.ObjectId{profile.Id},
		DisplayTemplate: constants.NotificationTemplate["TRANSACTIONAL"],
		EntityGroupId:   post.Id.Hex(),
		ActionId:        post.Id,
		ActionType:      "post",
		Purpose:         constants.NotificationPurpose["POST_SHARE"],
		Entities:        entities,
		NUUID:           "",
		Message:         msgStr,
		MessageMeta:     messageMeta,
		MessageHtml:     htmlMsgStr,
		DeepLink:        deepLink,
	})

	icon := mongo.ExtractThumbNailFromPost(post)

	if icon == "" {
		icon = profile.Avatar
	}

	msg := firebase.ManchMessage{
		Title:    title,
		Message:  msgStr,
		Icon:     mongo.ExtractThumbNailFromPost(post),
		DeepLink: deepLink,
		Id:       notification.NId,
	}
	//firebase.SendMessage(msg, "frgp37gfvFg:APA91bHbnbfoX-bp3M_3k-ceD7E4fZ73fcmVL4b5DGB5cQn-fFEvfbj3aAI9g0wXozyApIb-6wGsJauf67auK1p3Ins5Ff7IXCN161fb5JJ5pfBnTZ4LEcRUatO6wimsbiS7EANoGDr4")
	tokens := mongo.GetTokensByProfiles([]bson.ObjectId{postCreator.Id})
	fmt.Printf("\nGCM Message %+v\n", msg)
	if tokens != nil {
		for _, token := range tokens {
			go firebase.SendMessage(msg, token.Token, notification)
		}
	} else {
		fmt.Printf("No token\n")
	}

	fmt.Printf("Processed a Post Share on subject %s! with Post ID %s\n", subj, share.Id)

}
