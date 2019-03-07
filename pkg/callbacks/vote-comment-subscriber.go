package callbacks

import (
	"fmt"
	"notification-service/pkg/constants"
	"notification-service/pkg/firebase"
	"notification-service/pkg/i18n"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"
	"notification-service/pkg/utils"
	"strconv"
	"strings"

	"github.com/globalsign/mgo/bson"
)

func VoteCommentSubscriberCB(subj, reply string, v *subscribers.Vote) {
	fmt.Printf("Received a Vote on subject %s! with Vote %+v\n", subj, v)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in subscribers.VoteCommentSubscriber", r)
		}
	}()
	dir, err := strconv.Atoi(v.Direction)
	if err != nil {
		fmt.Println("Invalid Vote")
		//Do not process downvotes and unvote
		return
	}
	comment := mongo.GetCommentById(v.Resource)
	vote := comment.GetVote(v.Id)
	if vote.Created.ProfileId == comment.Created.ProfileId {
		//Self Vote
		fmt.Println("Self Vote")
		return
	}

	if dir < 1 {
		mongo.RemoveParticipants((comment.Id.Hex() + "_vote"), false, vote.Created.ProfileId)
		//Do not process downvotes and unvote
		return
	}

	err, post := mongo.GetPostById(comment.PostId.Hex())
	if err != nil {
		return
	}
	commentCreator := mongo.GetProfileById(comment.Created.ProfileId)
	// notification := mongo.CreateNotification(comment.Id, "like", "comment", vote.Created.ProfileId)

	entities := []mongo.Entity{
		{
			EntityId:   comment.Post.Id,
			EntityType: "post",
		},
		{
			EntityId:   comment.Id,
			EntityType: "comment",
		},
		{
			EntityId:   vote.Id,
			EntityType: "vote",
		},
	}
	notification := mongo.CreateNotification(mongo.NotificationModel{
		Receiver:        commentCreator.Id,
		Identifier:      comment.Id.Hex() + "_vote",
		Participants:    []bson.ObjectId{vote.Created.ProfileId},
		DisplayTemplate: constants.NotificationTemplate["TRANSACTIONAL"],
		EntityGroupId:   comment.Id.Hex(),
		ActionId:        vote.Id,
		ActionType:      "vote",
		Purpose:         constants.NotificationPurpose["VOTE"],
		Entities:        entities,
		NUUID:           "",
	})

	tokens := mongo.GetTokensByProfiles([]bson.ObjectId{comment.Created.ProfileId})

	commentTitle := utils.TruncateTitle(comment.Content, 4)

	data := i18n.DataModel{
		Name:    vote.Created.Name,
		Comment: commentTitle,
		Count:   comment.UpVotes,
	}

	var msgStr string
	var templateName string
	if comment.UpVotes > 1 {
		templateName = "comment_like_multi"
	} else {
		templateName = "comment_like_one"
	}

	msgStr = i18n.GetString(commentCreator.Language, templateName, data)
	htmlMsgStr := i18n.GetHtmlString(commentCreator.Language, templateName, data)
	msgStr = strings.Replace(msgStr, "\"\" ", "", 1)
	title := i18n.GetAppTitle(commentCreator.Language)

	messageMeta := mongo.MessageMeta{
		TemplateName: templateName,
		Template:     i18n.Strings[commentCreator.Language][templateName],
		Data:         data,
	}

	deepLink := "manch://posts/" + comment.PostId.Hex()
	// update notification message
	mongo.UpdateNotification(bson.M{"_id": notification.Id}, bson.M{
		"message":      msgStr,
		"message_meta": messageMeta,
		"message_html": htmlMsgStr,
		"deep_link":    deepLink,
	})

	msg := firebase.ManchMessage{
		Title:    title,
		Message:  msgStr,
		Icon:     mongo.ExtractThumbNailFromPost(post),
		DeepLink: deepLink,
		Id:       notification.NId,
	}

	fmt.Printf("\nGCM Message %+v\n", msg)
	if tokens != nil {
		for _, token := range tokens {
			go firebase.SendMessage(msg, token.Token, notification)
		}
	} else {
		fmt.Printf("No token\n")
	}
	fmt.Printf("Processed a Vote on subject %s! with Vote Id %s\n", subj, v.Id)

}
