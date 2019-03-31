package callbacks

import (
	"fmt"
	"notification-service/pkg/constants"
	"notification-service/pkg/firebase"
	"notification-service/pkg/i18n"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"
	"time"

	"github.com/globalsign/mgo/bson"
)

func MileStoneSubscriberCB(subj, reply string, m *subscribers.MileStone) {
	fmt.Printf("Received a new milestone on subject %s! with Value %+v\n", subj, m)

	profile := mongo.GetProfileById(bson.ObjectIdHex(m.ProfileId))
	if m.MileStone == constants.MileStones["100_COIN_MILESTONE"] {
		// send notification for this milestone

		tokens := mongo.GetTokensByProfiles([]bson.ObjectId{profile.Id})
		var msgStrTitle, msgStrText string
		var templateTitle, templateText string
		templateTitle = "100_coin_milestone_title"
		templateText = "100_coin_milestone_text"
		data := i18n.DataModel{
			Name: profile.Name,
		}
		msgStrTitle = i18n.GetString(profile.Language, templateTitle, data)
		msgStrText = i18n.GetString(profile.Language, templateText, data)
		htmlMsgStr := i18n.GetHtmlString(profile.Language, templateTitle, data)

		bigPicture := i18n.GetString(profile.Language, "100_coin_milestone_image", data)

		messageMeta := mongo.MessageMeta{
			TemplateName: templateTitle,
			Template:     i18n.Strings[profile.Language][templateTitle],
			Data:         data,
		}
		deepLink := "manch://posts/"

		notification := mongo.CreateNotification(mongo.NotificationModel{
			Receiver:        profile.Id,
			Identifier:      profile.Id.Hex() + "_100_milestone",
			Participants:    []bson.ObjectId{profile.Id},
			DisplayTemplate: constants.NotificationTemplate["TRANSACTIONAL"],
			EntityGroupId:   profile.Id.Hex(),
			ActionId:        profile.Id,
			ActionType:      "milestone",
			Purpose:         constants.NotificationPurpose["100_COIN_MILESTONE"],
			Message:         msgStrTitle,
			MessageMeta:     messageMeta,
			MessageHtml:     htmlMsgStr,
			DeepLink:        deepLink,
		})

		msg := firebase.ManchMessage{
			Title:      msgStrTitle,
			Message:    msgStrText,
			DeepLink:   deepLink,
			BigPicture: bigPicture,
			Id:         notification.NId,
		}

		fmt.Printf("\nGCM Message %+v\n", msg)
		if tokens != nil {
			for _, token := range tokens {
				go firebase.SendMessage(msg, token.Token, notification)
			}
		} else {
			fmt.Printf("No token")
		}

		// update referrer's coin
		err, referralData := mongo.GetReferralsByQuery(bson.M{
			"profile_id":       m.ProfileId,
			"referring_params": bson.M{"$exists": true},
		})
		if err != nil {
			fmt.Println("error", err)
		} else {
			fmt.Printf("referal\n%+v", referralData)
			referredBy := referralData.ReferringParams["profile_id"].(string)

			// update coin
			mongo.UpdateProfileById(bson.ObjectIdHex(referredBy), bson.M{
				"$inc": bson.M{"profiles.$.total_coins": 100},
			})

			// send notification

			referrer := mongo.GetProfileById(bson.ObjectIdHex(referredBy))
			tokens := mongo.GetTokensByProfiles([]bson.ObjectId{referrer.Id})
			var msgStrTitle, msgStrText string
			var templateTitle, templateText string
			templateTitle = "100_coin_referral_title"
			templateText = "100_coin_referral_text"
			data := i18n.DataModel{
				Name:  referrer.Name,
				Name2: profile.Name,
			}
			msgStrTitle = i18n.GetString(referrer.Language, templateTitle, data)
			msgStrText = i18n.GetString(referrer.Language, templateText, data)
			htmlMsgStr := i18n.GetHtmlString(referrer.Language, templateTitle, data)
			bigPicture := i18n.GetString(profile.Language, "100_coin_referral_image", data)

			messageMeta := mongo.MessageMeta{
				TemplateName: templateTitle,
				Template:     i18n.Strings[referrer.Language][templateTitle],
				Data:         data,
			}
			deepLink := "manch://posts/"

			notification := mongo.CreateNotification(mongo.NotificationModel{
				Receiver:        referrer.Id,
				Identifier:      referrer.Id.Hex() + profile.Id.Hex() + "_100_referral",
				Participants:    []bson.ObjectId{referrer.Id},
				DisplayTemplate: constants.NotificationTemplate["TRANSACTIONAL"],
				EntityGroupId:   referrer.Id.Hex(),
				ActionId:        referrer.Id,
				ActionType:      "referral",
				Purpose:         constants.NotificationPurpose["100_COIN_REFERRAL"],
				Message:         msgStrTitle,
				MessageMeta:     messageMeta,
				MessageHtml:     htmlMsgStr,
				DeepLink:        deepLink,
			})

			msg := firebase.ManchMessage{
				Title:      msgStrTitle,
				Message:    msgStrText,
				DeepLink:   deepLink,
				BigPicture: bigPicture,
				Id:         notification.NId,
			}

			fmt.Printf("\nGCM Message %+v\n", msg)
			if tokens != nil {
				for _, token := range tokens {
					go firebase.SendMessage(msg, token.Token, notification)
				}
			} else {
				fmt.Printf("No token")
			}

		}

	}

	if m.MileStone == constants.MileStones["500_COIN_MILESTONE"] {
		badge := mongo.Badge{
			ResourceName: "ic_manch_member",
			Icon:         "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/badges/ic_manch_member.png",
		}
		currentMilestoneID:= "2"
		milestone := mongo.Milestone{
			Id:          bson.NewObjectId(),
			MileStoneId: currentMilestoneID,
			Name:        "Manch Member",
			Badge:       badge,
			Value:       500,
			Type:        "coin",
			Date:        time.Now(),
		}

		query := bson.M{
			"profiles._id": profile.Id,
			"profiles.achieved_milestones.milestone_id": bson.M{"$ne": currentMilestoneID},
		}

		update := bson.M{
			"$set":  bson.M{"profiles.$.current_badge": badge,  "profiles.$.current_milestone_id":  currentMilestoneID},
			"$push": bson.M{"profiles.$.achieved_milestones": milestone},
		}
		// Update current badge and achieved milestones
		err := mongo.UpdateUser(query, update)
		if err == nil {
			// update post and comment of this profile
			mongo.UpdatePostsByQuery(bson.M{
				"created.profile_id": profile.Id,
			}, bson.M{
				"$set": bson.M{"created.current_badge": badge},
			})

			mongo.UpdateCommentsByQuery(bson.M{
				"created.profile_id": profile.Id,
			}, bson.M{
				"$set": bson.M{"created.current_badge": badge},
			})
		}

	}

	if m.MileStone == constants.MileStones["10000_COIN_MILESTONE"] {
		badge := mongo.Badge{
			ResourceName: "ic_super_user",
			Icon:         "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/badges/ic_super_user.png",
		}
		currentMilestoneID := "4"
		milestone := mongo.Milestone{
			Id:          bson.NewObjectId(),
			MileStoneId: currentMilestoneID,
			Name:        "Super User",
			Badge:       badge,
			Value:       10000,
			Type:        "coin",
			Date:        time.Now(),
		}

		query := bson.M{
			"profiles._id": profile.Id,
			"profiles.achieved_milestones.milestone_id": bson.M{"$ne": currentMilestoneID},
		}

		update := bson.M{
			"$set":  bson.M{"profiles.$.current_badge": badge, "profiles.$.current_milestone_id":  currentMilestoneID},
			"$push": bson.M{"profiles.$.achieved_milestones": milestone},
		}
		// Update current badge and achieved milestones
		err := mongo.UpdateUser(query, update)
		if err == nil {
			// update post and comment of this profile
			mongo.UpdatePostsByQuery(bson.M{
				"created.profile_id": profile.Id,
			}, bson.M{
				"$set": bson.M{"created.current_badge": badge},
			})

			mongo.UpdateCommentsByQuery(bson.M{
				"created.profile_id": profile.Id,
			}, bson.M{
				"$set": bson.M{"created.current_badge": badge},
			})
		}

	}

	if m.MileStone == constants.MileStones["25000_COIN_MILESTONE"] {
		badge := mongo.Badge{
			ResourceName: "ic_manch_creator",
			Icon:         "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/badges/ic_manch_creator.png",
		}
		currentMilestoneID := "6"
		milestone := mongo.Milestone{
			Id:          bson.NewObjectId(),
			MileStoneId: currentMilestoneID,
			Name:        "Manch Creator",
			Badge:       badge,
			Value:       25000,
			Type:        "coin",
			Date:        time.Now(),
		}

		query := bson.M{
			"profiles._id": profile.Id,
			"profiles.achieved_milestones.milestone_id": bson.M{"$ne": currentMilestoneID},
		}

		update := bson.M{
			"$set":  bson.M{"profiles.$.current_badge": badge, "profiles.$.current_milestone_id":  currentMilestoneID},
			"$push": bson.M{"profiles.$.achieved_milestones": milestone},
		}
		// Update current badge and achieved milestones
		err := mongo.UpdateUser(query, update)
		if err == nil {
			// update post and comment of this profile
			mongo.UpdatePostsByQuery(bson.M{
				"created.profile_id": profile.Id,
			}, bson.M{
				"$set": bson.M{"created.current_badge": badge},
			})

			mongo.UpdateCommentsByQuery(bson.M{
				"created.profile_id": profile.Id,
			}, bson.M{
				"$set": bson.M{"created.current_badge": badge},
			})
		}

	}

	fmt.Printf("Processed a new milestone on subject %s! with Id %s\n", subj, m.ProfileId)
}
