package callbacks

import (
	"fmt"
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
		mongo.RemoveNotificationUser(comment.Id, "like", vote.Created.ProfileId)
		//Do not process downvotes and unvote
		return
	}

	post := mongo.GetPostById(comment.PostId.Hex())
	commentCreator := mongo.GetProfileById(comment.Created.ProfileId)
	// notification := mongo.CreateNotification(comment.Id, "like", "comment", vote.Created.ProfileId)
	notification := mongo.CreateNotification(mongo.NotificationModel{
		Receiver:        commentCreator.Id,
		Identifier:      comment.Id.Hex() + "_vote",
		Participants:    []bson.ObjectId{vote.Created.ProfileId},
		DisplayTemplate: "transactional",
		EntityGroupId:   comment.Id.Hex(),
		ActionId:        vote.Id,
		ActionType:      "vote",
		Purpose:         "vote",
		Entities:        []string{"post", "comment", "vote"},
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
	if comment.UpVotes > 1 {
		msgStr = i18n.GetString(commentCreator.Language, "comment_like_multi", data)
	} else {
		msgStr = i18n.GetString(commentCreator.Language, "comment_like_one", data)
	}
	msgStr = strings.Replace(msgStr, "\"\" ", "", 1)
		
	// update notification message
	mongo.UpdateNotificationMessage(notification.Id, msgStr)

	title := i18n.GetAppTitle(commentCreator.Language)
	msg := firebase.ManchMessage{
		Title:    title,
		Message:  msgStr,
		Icon:     mongo.ExtractThumbNailFromPost(post),
		DeepLink: "manch://posts/" + comment.PostId.Hex(),
		Id:       notification.Identifier,
	}

	fmt.Printf("\nGCM Message %+v\n", msg)
	if tokens != nil {
		for _, token := range tokens {
			go firebase.SendMessage(msg, token.Token)
		}
	} else {
		fmt.Printf("No token\n")
	}
	fmt.Printf("Processed a Vote on subject %s! with Vote Id %s\n", subj, v.Id)

}
