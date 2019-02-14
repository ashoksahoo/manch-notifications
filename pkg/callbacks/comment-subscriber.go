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

/**
This processes Comments from Posts
1) Get Comment Details and Unique commentator count
2) Validate self comment
3) Get Who created the post -> He gets the notification and we need his current lang
4) Get tokens from the above profile (Supports multiple device tokens.)
5) Create/Update Notification Table which has the meta info for the notificaiotn
6) Construct Data for i18n template
7) Generate template using template data and String Formatter
8) Create push notification
9) Fire the notifications in routines.

*/
func CommentSubscriberCB(subj, reply string, c *subscribers.Comment) {
	//fmt.Printf("\nNats MSG %+v", c)
	fmt.Printf("Received a comment on subject %s! with Comment %+v\n", subj, c)
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in subscribers.CommentSubscriber", r)
		}
	}()
	comment, uniqueCommentator := mongo.GetFullCommentById(c.Id)

	var commentEntity, replyEntity []mongo.Entity
	commentEntity = []mongo.Entity{
		{
			EntityId:   comment.Post.Id,
			EntityType: "post",
		},
		{
			EntityId:   comment.Id,
			EntityType: "comment",
		},
	}
	// get replied on comment
	var replyOnComment mongo.CommentModel
	if len(comment.Parents) >= 2 {
		// this is reply on comments
		replyOnComment = mongo.GetCommentById(comment.CommentId.Hex())
		replyEntity = []mongo.Entity{
			{
				EntityId:   comment.Post.Id,
				EntityType: "post",
			},
			{
				EntityId:   replyOnComment.Id,
				EntityType: "comment",
			},
		}
		/*
			Create notification for Multi reply on comments
		*/
		participants := mongo.GetRepliesByCommentId(comment.Post.Id, replyOnComment.Id, comment.Created.ProfileId)
		for _, participant := range participants {
			if participant == comment.Created.ProfileId || participant == replyOnComment.Created.ProfileId {
				continue
			}

			identity := participant.Hex() + "_" + comment.Post.Id.Hex() + "_" + replyOnComment.Id.Hex() + "_multi_reply"
			notif := mongo.CreateNotification(mongo.NotificationModel{
				Receiver:        participant,
				Identifier:      identity,
				Participants:    []bson.ObjectId{comment.Created.ProfileId},
				DisplayTemplate: constants.NotificationTemplate["TRANSACTIONAL"],
				EntityGroupId:   c.Id,
				ActionId:        comment.Id,
				ActionType:      "comment",
				Purpose:         constants.NotificationPurpose["MULTI_REPLY"],
				Entities:        replyEntity,
				NUUID:           "",
			})
			receiver := mongo.GetProfileById(participant)
			commentTitle := utils.TruncateTitle(replyOnComment.Content, 4)
			tokens := mongo.GetTokensByProfiles([]bson.ObjectId{participant})
			data := i18n.DataModel{
				Name:    comment.Created.Name,
				Comment: commentTitle,
			}
			var templateName, msgStr string
			templateName = "reply_on_same_comment_one"

			msgStr = i18n.GetString(receiver.Language, templateName, data)
			msgStr = strings.Replace(msgStr, "\"\" ", "", 1)
			title := i18n.GetAppTitle(receiver.Language)

			messageMeta := mongo.MessageMeta{
				Template: templateName,
				Data:     data,
			}
			// update notification message
			mongo.UpdateNotification(bson.M{"_id": notif.Id}, bson.M{
				"message":      msgStr,
				"message_meta": messageMeta,
			})

			msg := firebase.ManchMessage{
				Title:      title,
				Message:    msgStr,
				Icon:       mongo.ExtractThumbNailFromPost(comment.Post),
				DeepLink:   "manch://posts/" + comment.PostId.Hex(),
				BadgeCount: strconv.Itoa(comment.Post.CommentCount),
				Id:         notif.NId,
			}

			fmt.Printf("\nGCM Message %+v\n", msg)
			//firebase.SendMessage(msg, "frgp37gfvFg:APA91bHbnbfoX-bp3M_3k-ceD7E4fZ73fcmVL4b5DGB5cQn-fFEvfbj3aAI9g0wXozyApIb-6wGsJauf67auK1p3Ins5Ff7IXCN161fb5JJ5pfBnTZ4LEcRUatO6wimsbiS7EANoGDr4")
			if tokens != nil {
				for _, token := range tokens {
					go firebase.SendMessage(msg, token.Token, notif)
				}
			} else {
				fmt.Printf("No token")
			}
			fmt.Printf("Notification created for multi reply with id %s", notif.Id.Hex())
		}
	}
	// reply on comment
	if len(comment.Parents) >= 2 && replyOnComment.Created.ProfileId != comment.Created.ProfileId {
		fmt.Println("reply on comments")
		// if replyOnComment.Created.ProfileId == comment.Created.ProfileId {
		// 	// Self Reply on comment
		// 	fmt.Println("Self reply on comments")
		// 	return;
		// }
		fmt.Printf("reply on comments %+v\n", replyOnComment)
		// get replied on comment creator
		replyOnCommentCreator := mongo.GetProfileById(replyOnComment.Created.ProfileId)
		// notification1 := mongo.CreateNotification(replyOnComment.Id, "comment", "comment", comment.Created.ProfileId)
		count := mongo.GetReplierCount(replyOnComment.Id, replyOnComment.Created.ProfileId) - 1
		notification1 := mongo.CreateNotification(mongo.NotificationModel{
			Receiver:        replyOnComment.Created.ProfileId,
			Identifier:      replyOnComment.Id.Hex() + "_reply",
			Participants:    []bson.ObjectId{comment.Created.ProfileId},
			DisplayTemplate: constants.NotificationTemplate["TRANSACTIONAL"],
			EntityGroupId:   replyOnComment.Id.Hex(),
			ActionId:        comment.Id,
			ActionType:      "comment",
			Purpose:         constants.NotificationPurpose["REPLY"],
			Entities:        replyEntity,
			NUUID:           "",
		})
		tokens1 := mongo.GetTokensByProfiles([]bson.ObjectId{replyOnComment.Created.ProfileId})
		// comment title
		// count := len(notification1.Participants) - 1
		fmt.Println("Comment count: ", count)
		commentTitle := utils.TruncateTitle(replyOnComment.Content, 4)
		data1 := i18n.DataModel{
			Name:    comment.Created.Name,
			Comment: commentTitle,
			Count:   count,
		}

		var msgStr1 string
		var templateName string
		if count > 0 {
			templateName = "comment_reply_multi"
		} else {
			templateName = "comment_reply_one"
		}

		msgStr1 = i18n.GetString(replyOnCommentCreator.Language, templateName, data1)
		msgStr1 = strings.Replace(msgStr1, "\"\" ", "", 1)
		title := i18n.GetAppTitle(replyOnCommentCreator.Language)

		messageMeta := mongo.MessageMeta{
			Template: templateName,
			Data:     data1,
		}
		// update notification message
		mongo.UpdateNotification(bson.M{"_id": notification1.Id}, bson.M{
			"message":      msgStr1,
			"message_meta": messageMeta,
		})

		msg := firebase.ManchMessage{
			Title:      title,
			Message:    msgStr1,
			Icon:       mongo.ExtractThumbNailFromPost(comment.Post),
			DeepLink:   "manch://posts/" + comment.PostId.Hex(),
			BadgeCount: strconv.Itoa(replyOnComment.CommentCount),
			Id:         notification1.NId,
		}
		fmt.Printf("\nGCM Message %+v\n", msg)
		if tokens1 != nil {
			for _, token := range tokens1 {
				go firebase.SendMessage(msg, token.Token, notification1)
			}
		} else {
			fmt.Printf("No token")
		}
		fmt.Println("end reply on comments")
	}

	/*
		Comment on same post logic
		1. create unique identifier <profileId>_<post_id>_<purpose>
		2. get all previous commentator and update their participants
	*/
	if len(comment.Parents) < 2 {
		participants := mongo.GetCommentsByPostId(comment.Post.Id, comment.Created.ProfileId)
		for _, participant := range participants {
			if participant == comment.Created.ProfileId || participant == comment.Post.Created.ProfileId {
				continue
			}
			identity := participant.Hex() + "_" + comment.Post.Id.Hex() + "_multi_comment"
			notif := mongo.CreateNotification(mongo.NotificationModel{
				Receiver:        participant,
				Identifier:      identity,
				Participants:    []bson.ObjectId{comment.Created.ProfileId},
				DisplayTemplate: constants.NotificationTemplate["TRANSACTIONAL"],
				EntityGroupId:   c.Id,
				ActionId:        comment.Id,
				ActionType:      "comment",
				Purpose:         constants.NotificationPurpose["MULTI_COMMENT"],
				Entities:        commentEntity,
				NUUID:           "",
			})
			participantCount := len(notif.Participants)

			receiver := mongo.GetProfileById(participant)
			postTitle := utils.TruncateTitle(comment.Post.Title, 4)
			data := i18n.DataModel{
				Name:  comment.Created.Name,
				Post:  postTitle,
				Count: participantCount - 1,
			}
			var templateName, msgStr string
			if participantCount > 1 {
				templateName = "comment_on_same_post_multi"
			} else {
				templateName = "comment_on_same_post_one"
			}

			msgStr = i18n.GetString(receiver.Language, templateName, data)
			msgStr = strings.Replace(msgStr, "\"\" ", "", 1)

			messageMeta := mongo.MessageMeta{
				Template: templateName,
				Data:     data,
			}
			// update notification message
			mongo.UpdateNotification(bson.M{"_id": notif.Id}, bson.M{
				"message":      msgStr,
				"message_meta": messageMeta,
			})

			fmt.Printf("Notification created for multi comments with id %s", notif.Id.Hex())
		}
	}

	// Comment notification
	if comment.Post.Created.ProfileId == comment.Created.ProfileId {
		//Self comment
		fmt.Println("Self Comment")
		return
	}
	postCreator := mongo.GetProfileById(comment.Post.Created.ProfileId)
	tokens := mongo.GetTokensByProfiles([]bson.ObjectId{comment.Post.Created.ProfileId})

	// update commentCreator's coin
	mongo.CreateUserCoin(mongo.UserCoinsModel{
		ProfileId: comment.Created.ProfileId,
		CoinsEarned: 5,
		Action: "comment",
	})

	// update postCreator's coin
	mongo.CreateUserCoin(mongo.UserCoinsModel{
		ProfileId: postCreator.Id,
		CoinsEarned: 5,
		Action: "comment",
	})

	// update user score
	mongo.CreateUserScore(mongo.UserScore{
		ProfileId:   comment.Created.ProfileId,
		CommunityId: comment.Post.CommunityIds[0],
		Score:       1,
		UserType:    comment.Created.UserType,
	})

	notification := mongo.CreateNotification(mongo.NotificationModel{
		Receiver:        postCreator.Id,
		Identifier:      c.Id + "_comment",
		Participants:    []bson.ObjectId{comment.Created.ProfileId},
		DisplayTemplate: constants.NotificationTemplate["TRANSACTIONAL"],
		EntityGroupId:   c.Id,
		ActionId:        comment.Id,
		ActionType:      "comment",
		Purpose:         constants.NotificationPurpose["COMMENT"],
		Entities:        commentEntity,
		NUUID:           "",
	})

	postTitle := utils.TruncateTitle(comment.Post.Title, 4)
	data := i18n.DataModel{
		Name:  comment.Created.Name,
		Count: uniqueCommentator - 1,
		Post:  postTitle,
	}

	var msgStr string
	var templateName string
	if uniqueCommentator > 1 {
		templateName = "comment_multi"
	} else {
		templateName = "comment_one"
	}

	msgStr = i18n.GetString(postCreator.Language, templateName, data)
	if uniqueCommentator > 25 {
		msgStr = "👏 " + msgStr
	}
	msgStr = strings.Replace(msgStr, "\"\" ", "", 1)
	title := i18n.GetAppTitle(postCreator.Language)

	messageMeta := mongo.MessageMeta{
		Template: templateName,
		Data:     data,
	}
	// update notification message
	mongo.UpdateNotification(bson.M{"_id": notification.Id}, bson.M{
		"message":      msgStr,
		"message_meta": messageMeta,
	})

	msg := firebase.ManchMessage{
		Title:      title,
		Message:    msgStr,
		Icon:       mongo.ExtractThumbNailFromPost(comment.Post),
		DeepLink:   "manch://posts/" + comment.PostId.Hex(),
		BadgeCount: strconv.Itoa(comment.Post.CommentCount),
		Id:         notification.NId,
	}

	fmt.Printf("\nGCM Message %+v\n", msg)
	//firebase.SendMessage(msg, "frgp37gfvFg:APA91bHbnbfoX-bp3M_3k-ceD7E4fZ73fcmVL4b5DGB5cQn-fFEvfbj3aAI9g0wXozyApIb-6wGsJauf67auK1p3Ins5Ff7IXCN161fb5JJ5pfBnTZ4LEcRUatO6wimsbiS7EANoGDr4")
	if tokens != nil {
		for _, token := range tokens {
			go firebase.SendMessage(msg, token.Token, notification)
		}
	} else {
		fmt.Printf("No token")
	}
	fmt.Printf("Processed a comment on subject %s! with Comment %s\n", subj, c.Id)

}
