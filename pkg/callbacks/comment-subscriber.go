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

	commentCreator := mongo.GetProfileById(comment.Created.ProfileId)

	// data notification for app ratings popup
	if commentCreator.CommentsCount == 5 {
		entities := []mongo.Entity{
			{
				EntityId:   comment.Id,
				EntityType: "comment",
			},
		}
		notification := mongo.CreateNotification(mongo.NotificationModel{
			Receiver:        commentCreator.Id,
			Identifier:      commentCreator.Id.Hex() + "_user_review",
			Participants:    []bson.ObjectId{commentCreator.Id},
			PlaceHolderIcon: []string{comment.Created.Avatar},
			DisplayTemplate: constants.NotificationTemplate["TRANSACTIONAL"],
			ActionId:        comment.Id,
			ActionType:      "comment",
			Purpose:         constants.NotificationPurpose["USER_REVIEW"],
			Entities:        entities,
			PushType:        "manch:D",
		})

		tokens := mongo.GetTokensByProfiles([]bson.ObjectId{commentCreator.Id})
		msg := firebase.ManchMessage{
			Title:     "",
			Message:   "",
			Namespace: "manch:D",
			Id:        notification.NId,
		}
		if tokens != nil {
			for _, token := range tokens {
				fmt.Println("successfully sent data message")
				go firebase.SendMessage(msg, token.Token, notification)
			}
		} else {
			fmt.Printf("No token")
		}

	}

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

	// create community stats
	community := mongo.GetCommunityById(comment.Post.CommunityIds[0].Hex())
	mongo.CreateCommunityStats(mongo.CommunityStatsModel{
		CommunityId:           comment.Post.CommunityIds[0],
		Action:                "comment",
		EntityId:              comment.Id,
		EntityType:            "comment",
		ProfileId:             comment.Created.ProfileId,
		CommentsCount:         1,
		ActionSource:          comment.Post.SourcedBy,
		CommunityCreatorType:  community.Created.Type,
		ActorType:             comment.Created.UserType,
		ParticipatingEntityId: comment.PostId,
	})

	// get replied on comment
	var replyOnComment = mongo.GetCommentById(comment.CommentId.Hex())

	// notification for reply on a comment
	if len(comment.Parents) >= 2 && replyOnComment.Created.ProfileId != comment.Created.ProfileId && replyOnComment.Created.ProfileId != replyOnComment.Post.Created.ProfileId {
		fmt.Printf("reply on comments %+v\n", replyOnComment)
		// get replied on comment creator
		replyOnCommentCreator := mongo.GetProfileById(replyOnComment.Created.ProfileId)
		// notification1 := mongo.CreateNotification(replyOnComment.Id, "comment", "comment", comment.Created.ProfileId)
		count := mongo.GetReplierCount(replyOnComment.Id, replyOnComment.Created.ProfileId) - 1
		// comment title
		// count := len(notification1.Participants) - 1
		fmt.Println("Comment count: ", count)
		commentTitle := utils.TruncateTitle(replyOnComment.Content, 4)
		data1 := i18n.DataModel{
			Name:    comment.Created.Name,
			Comment: commentTitle,
			Count:   count,
		}

		var msgStr1, htmlMsgStr1 string
		var templateName string
		if count > 0 {
			templateName = "comment_reply_multi"
		} else {
			templateName = "comment_reply_one"
		}

		msgStr1 = i18n.GetString(replyOnCommentCreator.Language, templateName, data1)
		msgStr1 = strings.Replace(msgStr1, "\"\" ", "", 1)
		htmlMsgStr1 = i18n.GetHtmlString(replyOnCommentCreator.Language, templateName, data1)
		title := i18n.GetAppTitle(replyOnCommentCreator.Language)

		messageMeta := mongo.MessageMeta{
			TemplateName: templateName,
			Template:     i18n.Strings[replyOnCommentCreator.Language][templateName],
			Data:         data1,
		}
		deepLink := "manch://posts/" + comment.PostId.Hex()

		notification1 := mongo.CreateNotification(mongo.NotificationModel{
			Receiver:        replyOnComment.Created.ProfileId,
			Identifier:      replyOnComment.Id.Hex() + "_reply",
			Participants:    []bson.ObjectId{comment.Created.ProfileId},
			PlaceHolderIcon: []string{comment.Created.Avatar},
			DisplayTemplate: constants.NotificationTemplate["TRANSACTIONAL"],
			EntityGroupId:   replyOnComment.Id.Hex(),
			ActionId:        comment.Id,
			ActionType:      "comment",
			Purpose:         constants.NotificationPurpose["REPLY"],
			Entities:        replyEntity,
			NUUID:           "",
			Message:         msgStr1,
			MessageMeta:     messageMeta,
			MessageHtml:     htmlMsgStr1,
			DeepLink:        deepLink,
		})

		icon := mongo.ExtractThumbNailFromPost(comment.Post)

		if icon == "" {
			icon = comment.Created.Avatar
		}

		msg := firebase.ManchMessage{
			Title:      title,
			Message:    msgStr1,
			Icon:       icon,
			DeepLink:   deepLink,
			BadgeCount: strconv.Itoa(replyOnComment.CommentCount),
			Id:         notification1.NId,
		}

		commentBatch := []int{1, 2, 3, 5, 25, 50, 75, 100}

		if utils.Contains(commentBatch, count) || ((count % 50) == 0) {
			fmt.Printf("\nGCM Message %+v\n", msg)
			tokens1 := mongo.GetTokensByProfiles([]bson.ObjectId{replyOnComment.Created.ProfileId})
			if tokens1 != nil {
				for _, token := range tokens1 {
					go firebase.SendMessage(msg, token.Token, notification1)
				}
			} else {
				fmt.Printf("No token")
			}
		}

		fmt.Println("end reply on comments")
	}

	// Comment notification
	if comment.Post.Created.ProfileId == comment.Created.ProfileId {
		//Self comment
		fmt.Println("Self Comment")
		return
	}
	postCreator := mongo.GetProfileById(comment.Post.Created.ProfileId)

	// update commentCreator's coin
	mongo.CreateUserCoin(mongo.UserCoinsModel{
		ProfileId:   comment.Created.ProfileId,
		CoinsEarned: 2,
		Action:      "comment",
	})

	// update postCreator's coin
	mongo.CreateUserCoin(mongo.UserCoinsModel{
		ProfileId:   postCreator.Id,
		CoinsEarned: 3,
		Action:      "comment",
	})

	// update user score
	mongo.CreateUserScore(mongo.UserScore{
		ProfileId:   comment.Created.ProfileId,
		CommunityId: comment.Post.CommunityIds[0],
		Score:       1,
		UserType:    comment.Created.UserType,
	})

	postTitle := utils.TruncateTitle(comment.Post.Title, 4)
	data := i18n.DataModel{
		Name:  comment.Created.Name,
		Count: uniqueCommentator - 1,
		Post:  postTitle,
	}

	var msgStr, htmlMsgStr string
	var templateName string
	if uniqueCommentator > 1 {
		templateName = "comment_multi"
	} else {
		templateName = "comment_one"
	}

	msgStr = i18n.GetString(postCreator.Language, templateName, data)
	htmlMsgStr = i18n.GetHtmlString(postCreator.Language, templateName, data)
	if uniqueCommentator > 25 {
		msgStr = "👏 " + msgStr
	}
	msgStr = strings.Replace(msgStr, "\"\" ", "", 1)
	title := i18n.GetAppTitle(postCreator.Language)

	messageMeta := mongo.MessageMeta{
		TemplateName: templateName,
		Template:     i18n.Strings[postCreator.Language][templateName],
		Data:         data,
	}
	deepLink := "manch://posts/" + comment.PostId.Hex()
	notification := mongo.CreateNotification(mongo.NotificationModel{
		Receiver:        postCreator.Id,
		Identifier:      c.Id + "_comment",
		Participants:    []bson.ObjectId{comment.Created.ProfileId},
		PlaceHolderIcon: []string{comment.Created.Avatar},
		DisplayTemplate: constants.NotificationTemplate["TRANSACTIONAL"],
		EntityGroupId:   c.Id,
		ActionId:        comment.Id,
		ActionType:      "comment",
		Purpose:         constants.NotificationPurpose["COMMENT"],
		Entities:        commentEntity,
		NUUID:           "",
		Message:         msgStr,
		MessageMeta:     messageMeta,
		MessageHtml:     htmlMsgStr,
		DeepLink:        deepLink,
	})

	icon := mongo.ExtractThumbNailFromPost(comment.Post)

	if icon == "" {
		icon = comment.Created.Avatar
	}

	msg := firebase.ManchMessage{
		Title:      title,
		Message:    msgStr,
		Icon:       icon,
		DeepLink:   deepLink,
		BadgeCount: strconv.Itoa(comment.Post.CommentCount),
		Id:         notification.NId,
	}

	commentBatch := []int{1, 2, 3, 5, 25, 50, 75, 100}

	if utils.Contains(commentBatch, uniqueCommentator) || ((uniqueCommentator % 50) == 0) {
		fmt.Printf("\nGCM Message %+v\n", msg)
		tokens := mongo.GetTokensByProfiles([]bson.ObjectId{comment.Post.Created.ProfileId})
		if tokens != nil {
			for _, token := range tokens {
				go firebase.SendMessage(msg, token.Token, notification)
			}
		} else {
			fmt.Printf("No token")
		}
	}

	fmt.Printf("Processed a comment on subject %s! with Comment %s\n", subj, c.Id)

}
