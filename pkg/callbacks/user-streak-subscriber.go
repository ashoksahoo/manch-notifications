package callbacks

import (
	"fmt"
	"notification-service/pkg/constants"
	"notification-service/pkg/firebase"
	"notification-service/pkg/i18n"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"
	"notification-service/pkg/utils"
	"time"

	"github.com/globalsign/mgo/bson"
)

func UserStreakCB(subj, reply string, userStreak *subscribers.UserStreak) {
	fmt.Printf("Received a user streak on subject %s! with Value %+v\n", subj, userStreak)

	profile := mongo.GetProfileById(bson.ObjectIdHex(userStreak.ProfileId))

	if (subscribers.Streak{}) == userStreak.CurrentStreak {
		fmt.Println("reset streak")
		// update last and longest streak on profile
		query := bson.M{"profiles._id": profile.Id}

		update := bson.M{
			"$set":  bson.M{"profiles.$.last_streak": userStreak.LastStreak, "profiles.$.longest_streak": userStreak.LongestStreak},
		}
		err := mongo.UpdateUser(query, update)
		if err != nil {
			fmt.Println("error while updating last and longest streak", err)
		} else {
			fmt.Println("update last and longest streak on profile", profile.Id.Hex())
		}
		return
	}

	// update current streak on profile
	err := mongo.UpdateUser(bson.M{
		"profiles._id": profile.Id,
	}, bson.M{
		"$set":  bson.M{"profiles.$.current_streak": userStreak.CurrentStreak},
	})

	if err != nil {
		fmt.Println("error while updating last and current streak", err)
	} else {
		fmt.Println("update current streak on profile", profile.Id.Hex())
	}

	if utils.IncludesInt([]int{1, 7, 30, 100}, userStreak.CurrentStreak.StreakLength) {

		var resourceName, resourceIcon, milestoneId, milestoneName, bigPictureTemplateName, notifIdentifierText, notifPurpose string
		var milestoneValue int
		if userStreak.CurrentStreak.StreakLength == 1 {
			resourceName = "ic_milestone_1_day_steak"
			resourceIcon = "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/badges/ic_milestone_1_day_steak.png"
			milestoneId = "0"
			milestoneName = "1 day streak"
			milestoneValue = 1
			bigPictureTemplateName = "streak_milestone_image_1"
			notifIdentifierText = "milestone_streak_1"
			notifPurpose = constants.NotificationPurpose["1_STREAK_MILESTONE"]
		} else if userStreak.CurrentStreak.StreakLength == 7 {
			resourceName = "ic_milestone_7_day_steak"
			resourceIcon = "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/badges/ic_milestone_7_day_steak.png"
			milestoneId = "1"
			milestoneName = "7 day streak"
			milestoneValue = 7
			bigPictureTemplateName = "streak_milestone_image_7"
			notifIdentifierText = "milestone_streak_7"
			notifPurpose = constants.NotificationPurpose["7_STREAK_MILESTONE"]
		} else if userStreak.CurrentStreak.StreakLength == 30 {
			resourceName = "ic_milestone_30_day_steak"
			resourceIcon = "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/badges/ic_milestone_30_day_steak.png"
			milestoneId = "3"
			milestoneName = "30 day streak"
			milestoneValue = 30
			bigPictureTemplateName = "streak_milestone_image_30"
			notifIdentifierText = "milestone_streak_30"
			notifPurpose = constants.NotificationPurpose["30_STREAK_MILESTONE"]
		} else if userStreak.CurrentStreak.StreakLength == 100 {
			resourceName = "ic_milestone_100_day_steak"
			resourceIcon = "https://s3.ap-south-1.amazonaws.com/manch-dev/notifications/badges/ic_milestone_100_day_steak.png"
			milestoneId = "5"
			milestoneName = "100 day streak"
			milestoneValue = 100
			bigPictureTemplateName = "streak_milestone_image_100"
			notifIdentifierText = "milestone_streak_100"
			notifPurpose = constants.NotificationPurpose["100_STREAK_MILESTONE"]
		}

		badge := mongo.Badge{
			ResourceName: resourceName,
			Icon:         resourceIcon,
		}
		milestone := mongo.Milestone{
			Id:          bson.NewObjectId(),
			MileStoneId: milestoneId,
			Name:        milestoneName,
			Badge:       badge,
			Value:       milestoneValue,
			Type:        "streak",
			Date:        time.Now(),
		}

		query := bson.M{
			"profiles": bson.M{
				"$elemMatch": bson.M{
					"_id":                              profile.Id,
					"achieved_milestones.milestone_id": bson.M{"$ne": milestoneId},
				},
			},
		}

		update := bson.M{
			"$set":  bson.M{"profiles.$.current_badge": badge, "profiles.$.current_milestone_id": milestoneId},
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

			// send notification
			templateName := "streak_milestone"
			data := i18n.DataModel{
				Name:  profile.Name,
				Name2: milestoneName, // milestone name
			}
			msgStr := i18n.GetString(profile.Language, templateName, data)
			htmlMsgStr := i18n.GetString(profile.Language, templateName, data)
			title := i18n.GetAppTitle(profile.Language)
			bigPicture := i18n.GetString(profile.Language, bigPictureTemplateName, data)

			messageMeta := mongo.MessageMeta{
				TemplateName: templateName,
				Template:     i18n.Strings[profile.Language][templateName],
				Data:         data,
			}
			deepLink := "manch://profile/" + profile.Id.Hex()
			notification := mongo.CreateNotification(mongo.NotificationModel{
				Receiver:        profile.Id,
				Identifier:      profile.Id.Hex() + notifIdentifierText,
				Participants:    []bson.ObjectId{profile.Id},
				DisplayTemplate: constants.NotificationTemplate["TRANSACTIONAL"],
				EntityGroupId:   profile.Id.Hex(),
				ActionId:        profile.Id,
				ActionType:      "streak_milestone",
				Purpose:         notifPurpose,
				Message:         msgStr,
				MessageMeta:     messageMeta,
				MessageHtml:     htmlMsgStr,
				DeepLink:        deepLink,
			})

			msg := firebase.ManchMessage{
				Title:      title,
				Message:    msgStr,
				DeepLink:   deepLink,
				BigPicture: bigPicture,
				Id:         notification.NId,
			}
			tokens := mongo.GetTokensByProfiles([]bson.ObjectId{profile.Id})
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

	fmt.Printf("Processed a user streak on subject %s! with profile Id %s\n", subj, userStreak.ProfileId)
}
