package callbacks

import (
	"fmt"
	"math/rand"
	"notification-service/pkg/constants"
	"notification-service/pkg/firebase"
	"notification-service/pkg/i18n"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"
	"notification-service/pkg/utils"
	"time"

	"github.com/globalsign/mgo/bson"
)

func UserInactiveSubscriberCB(subj, reply string, u *subscribers.UserInactive) {
	fmt.Printf("Received user inactive event on subject %s! with Value %+v\n", subj, u)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in subscribers.UserInactiveSubscriber", r)
		}
	}()

	profileId := u.ProfileId

	err, posts := mongo.GetAllPostsByQuery(bson.M{
		"created.profile_id": bson.ObjectIdHex(profileId),
	}, 5)

	fmt.Println("posts are", posts)
	if err != nil {
		return
	}
	postLength := len(posts)
	randomIndex := utils.Random(0, postLength)
	fmt.Println("random index", randomIndex)
	post := posts[randomIndex]

	fmt.Println("random posts", post)

	// get unique bot profiles
	m, botProfilesHi := mongo.GetBotProfilesIds("hi")
	n, botProfilesTe := mongo.GetBotProfilesIds("te")
	n = m + n
	botProfilesIds := append(botProfilesHi, botProfilesTe...)
	// shuffle profiles
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(n, func(i, j int) { botProfilesIds[i], botProfilesIds[j] = botProfilesIds[j], botProfilesIds[i] })

	voteCreatorsList := mongo.GetAllVotedUserIncludingScheduled(bson.M{"resource": post.Id})

	botProfiles := utils.Difference(botProfilesIds, voteCreatorsList)

	// add 1 like
	botProfile := mongo.GetProfileById(bson.ObjectIdHex(botProfiles[0]))

	fmt.Println("bot profiles selected", botProfile)

	err = mongo.AddVote(mongo.VoteModelPost{
		Resource:     mongo.PostModel{Id: post.Id},
		Value:        1,
		User:         botProfile.Id,
		ResourceType: "post",
		Created: mongo.Creator{
			Id:        bson.NewObjectId(),
			ProfileId: botProfile.Id,
			Name:      botProfile.Name,
			Avatar:    botProfile.Avatar,
			UserType:  botProfile.Avatar,
			Date:      time.Now(),
		},
	})

	fmt.Println("added vote with")

	if err != nil {
		return
	}

	// send notification
	templateName := "post_like_one"
	data := i18n.DataModel{
		Name: botProfile.Name,
	}
	language := "hi"
	msgStr := i18n.GetString(language, templateName, data)
	htmlMsgStr := i18n.GetHtmlString(language, templateName, data)
	title := i18n.GetAppTitle(language)

	fmt.Println("message are", msgStr)

	messageMeta := mongo.MessageMeta{
		TemplateName: templateName,
		Template:     i18n.Strings[language][templateName],
		Data:         data,
	}
	// update notification message
	deepLink := "manch://posts/" + post.Id.Hex()

	notification := mongo.NotificationModel{
		Receiver:        bson.ObjectIdHex(profileId),
		Identifier:      post.Id.Hex() + "_inactive_vote",
		Participants:    []bson.ObjectId{bson.ObjectIdHex(profileId)},
		DisplayTemplate: constants.NotificationTemplate["TRANSACTIONAL"],
		EntityGroupId:   post.Id.Hex(),
		ActionId:        post.Id,
		ActionType:      "vote",
		Purpose:         constants.NotificationPurpose["VOTE"],
		NUUID:           "",
		Message:         msgStr,
		MessageMeta:     messageMeta,
		MessageHtml:     htmlMsgStr,
		DeepLink:        deepLink,
	}
	notification = mongo.CreateNotification(notification)

	icon := mongo.ExtractThumbNailFromPost(post)

	if icon == "" {
		icon = botProfile.Avatar
	}

	msg := firebase.ManchMessage{
		Title:    title,
		Message:  msgStr,
		Icon:     icon,
		DeepLink: deepLink,
		Id:       notification.NId,
	}

	fmt.Printf("\nGCM Message %+v\n", msg)
	tokens := mongo.GetTokensByProfiles([]bson.ObjectId{bson.ObjectIdHex(profileId)})
	if tokens != nil {
		for _, token := range tokens {
			fmt.Println("sending notification")
			go firebase.SendMessage(msg, token.Token, notification)
		}
	} else {
		fmt.Printf("No token")
	}
	fmt.Printf("Processed user inactive event on subject %s! with Value %+v\n", subj, u.ProfileId)
}
