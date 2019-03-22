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

	if uf.ResourceType != "user" {
		fmt.Println("Not a user resource follows")
		if uf.ResourceType == "community" {
			// create community stats
			mongo.CreateCommunityStats(mongo.CommunityStatsModel{
				CommunityId: bson.ObjectIdHex(uf.Resource),
				Action:      "community-follow",
				EntityId:    bson.ObjectIdHex(uf.ProfileId),
				EntityType:  "user",
				ProfileId:   bson.ObjectIdHex(uf.ProfileId),
			})

			// send join manch notification to admins for closed m-manch

			community := mongo.GetCommunityById(uf.Resource)

			if community.Type == "m_manch" && community.Visibility == "protected" {
				admins := community.Admins
				userRequested := mongo.GetProfileById(bson.ObjectIdHex(uf.ProfileId))

				adminProfilesIds := []string{}

				for _, admin := range admins {
					adminProfilesIds = append(adminProfilesIds, admin.ProfileId.Hex())
				}

				adminProfiles := mongo.GetProfilesByIds(adminProfilesIds)

				entities := []mongo.Entity{
					{
						EntityId:   community.Id,
						EntityType: "community",
					},
				}

				data := i18n.DataModel{
					Name:      userRequested.Name,
					Community: community.Name,
				}
				var templateName string
				if community.Visibility == "protected" {
					templateName = "join_manch_request_private"
				} else {
					templateName = "join_manch_request_public"
				}
				deepLink := ""
				for _, adminProfile := range adminProfiles {
					msgStr := i18n.GetString(adminProfile.Language, templateName, data)
					htmlMsgStr := i18n.GetHtmlString(adminProfile.Language, templateName, data)
					msgStr = strings.Replace(msgStr, "\"\" ", "", 1)
					title := i18n.GetAppTitle(adminProfile.Language)

					messageMeta := mongo.MessageMeta{
						TemplateName: templateName,
						Template:     i18n.Strings[adminProfile.Language][templateName],
						Data:         data,
					}
					purpose := constants.NotificationPurpose["JOIN_MANCH_REQUEST"]
					notification := mongo.CreateNotification(mongo.NotificationModel{
						Receiver:        adminProfile.Id,
						Identifier:      adminProfile.Id.Hex() + purpose,
						Participants:    []bson.ObjectId{adminProfile.Id},
						DisplayTemplate: constants.NotificationTemplate["TRANSACTIONAL"],
						EntityGroupId:   community.Id.Hex(),
						ActionId:        community.Id,
						ActionType:      "userfollow",
						Purpose:         purpose,
						Entities:        entities,
						Message:         msgStr,
						MessageHtml:     htmlMsgStr,
						DeepLink:        deepLink,
						MessageMeta:     messageMeta,
					})

					tokens := mongo.GetTokensByProfiles([]bson.ObjectId{adminProfile.Id})

					msg := firebase.ManchMessage{
						Title:    title,
						Message:  msgStr,
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

				}
			}

		}
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
			EntityId:   userFollow.Id,
			EntityType: "user_follow",
		},
	}
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
	})
	count := followsTo.FollowerCount
	data := i18n.DataModel{
		Name:  follower.Name,
		Count: count,
	}
	var msgStr string
	var templateName string
	if count == 9 {
		// 10th follower notification
		data.Name = followsTo.Name
		data.Count = count + 1
		notifImages := []string{"tenth_follower_image_1", "tenth_follower_image_2", "tenth_follower_image_3"}
		notifText := []string{"tenth_follower_text_1", "tenth_follower_text_2", "tenth_follower_text_3"}
		randomIndex := utils.Random(0, 3)
		notifTitleTemplate := "tenth_follower_title"
		notifTextTemplate := notifText[randomIndex]
		notifImageTemplate := notifImages[randomIndex]
		pnTitle := i18n.GetString(followsTo.Language, notifTitleTemplate, data)
		pnText := i18n.GetString(followsTo.Language, notifTextTemplate, data)
		pnBigImage := i18n.GetString(followsTo.Language, notifImageTemplate, data)

		htmlMsgText := i18n.GetHtmlString(followsTo.Language, notifTitleTemplate, data)

		messageMeta := mongo.MessageMeta{
			TemplateName: "tenth_follower_title," + notifTextTemplate + "," + notifImageTemplate,
			Template:     i18n.Strings[followsTo.Language][notifTitleTemplate],
			Data:         data,
		}
		deepLink := "manch://profile/" + followsTo.Id.Hex()
		// update notification message
		mongo.UpdateNotification(bson.M{"_id": notification.Id}, bson.M{
			"message":      pnTitle,
			"message_meta": messageMeta,
			"message_html": htmlMsgText,
			"deep_link":    deepLink,
		})

		msg := firebase.ManchMessage{
			Title:      pnTitle,
			Message:    pnText,
			BigPicture: pnBigImage,
			DeepLink:   deepLink,
			Id:         notification.NId,
		}
		fmt.Printf("\nGCM Message %+v\n", msg)
		if tokens != nil {
			for _, token := range tokens {
				go firebase.SendMessage(msg, token.Token, notification)
			}
		} else {
			fmt.Printf("No token\n")
		}
		fmt.Printf("Processed a User follow on subject %s! with user follow ID %s\n", subj, uf.Id)

		return

	} else if count > 1 {
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
		DeepLink: deepLink,
		Id:       notification.NId,
	}

	followNumbers := []int{1, 2, 3, 5, 20, 30, 40, 50, 60, 70, 80, 90, 100, 150, 200, 300, 350, 400, 450, 500, 550, 600, 650, 700, 800, 850, 900, 950}
	if utils.Contains(followNumbers, count+1) || ((count+1)%250 == 0) {
		//firebase.SendMessage(msg, "frgp37gfvFg:APA91bHbnbfoX-bp3M_3k-ceD7E4fZ73fcmVL4b5DGB5cQn-fFEvfbj3aAI9g0wXozyApIb-6wGsJauf67auK1p3Ins5Ff7IXCN161fb5JJ5pfBnTZ4LEcRUatO6wimsbiS7EANoGDr4")
		fmt.Printf("\nGCM Message %+v\n", msg)
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
