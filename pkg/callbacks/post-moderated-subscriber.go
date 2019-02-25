package callbacks

import (
	"fmt"
	"math"
	"math/rand"
	"notification-service/pkg/firebase"
	"notification-service/pkg/i18n"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"
	"notification-service/pkg/utils"
	"strconv"
	"time"

	"github.com/globalsign/mgo/bson"
)


func PostModeratedSubscriberCB(subj, reply string, p *subscribers.Post) {
	fmt.Printf("Received a post on subject %s! with Post %+v\n", subj, p)
	err, post := mongo.GetPostByQuery(bson.M{"_id": bson.ObjectIdHex(p.Id)})
	if err != nil {
		return
	}

	// i represents no of profiles
	n, botProfilesIds := mongo.GetBotProfilesIds(post.Language)
	randomBotIndex := utils.Random(0, n)

	postCreator := mongo.GetProfileById(post.Created.ProfileId)

	// schedule auto comment on post if it is good
	if post.PostLevel == "2" || post.PostLevel == "1" {
		if post.Language == "te" {
			return
		}
		var dbCommentKeys []string
		// get comment string from db
		err, commentString := mongo.GetCommentStringsByProfileId(postCreator.Id)
		if err != nil {
			dbCommentKeys = []string{}
		} else {
			dbCommentKeys = commentString.CommentStringIds
		}
		if len(dbCommentKeys) >= 5 {
			return
		}
		// get unique auto comments for this postCreator
		commentkeys := make([]string, 0, len(i18n.CommentStrings[postCreator.Language]))
		for k := range i18n.CommentStrings[postCreator.Language] {
			commentkeys = append(commentkeys, k)
		}

		// get set difference of commentKeys and dbCommentKeys
		keys := utils.Difference(commentkeys, dbCommentKeys)
		// no key unique left
		if len(keys) == 0 {
			return
		}
		randomCommentKeyIndex := utils.Random(0, len(keys))
		comment := i18n.CommentStrings[postCreator.Language][keys[randomCommentKeyIndex]]
		profileId := botProfilesIds[randomBotIndex]
		commentator := mongo.GetProfileById(bson.ObjectIdHex(profileId))
		commentCreator := mongo.Creator{
			Id:        bson.NewObjectId(),
			ProfileId: commentator.Id,
			Name:      commentator.Name,
			Avatar:    commentator.Avatar,
			UserType:  commentator.Type,
		}
		randomMinute := utils.Random(15, 30)
		fmt.Println("random Minute", randomMinute)
		scheduleTime := time.Now().Local().Add(time.Minute * time.Duration(randomMinute))
		fmt.Println("schedule time", scheduleTime)
		// schedule comments
		mongo.CreateCommentSchedule(comment, post.Id, commentCreator, scheduleTime)
		mongo.AddCommentStringToProfileId(postCreator.Id, keys[randomCommentKeyIndex])
		// schedule comments in 20-30 minutes random
	}

	// send notificaiton on block and warned
	send_notification := false
	entities := []mongo.Entity{
		{
			EntityId:   post.Id,
			EntityType: "post",
		},
	}
	notification := mongo.CreateNotification(mongo.NotificationModel{
		Receiver:        postCreator.Id,
		Identifier:      post.Id.Hex() + "_user_blocked",
		Participants:    []bson.ObjectId{postCreator.Id},
		DisplayTemplate: "transactional",
		ActionId:        post.Id,
		ActionType:      "post",
		Purpose:         "user.blocked",
		Entities:        entities,
		NUUID:           "",
	})

	blockedStatus := map[string]string{}
	var reason string
	// warning for 1st and 2nd delete post
	// block on every 3rd delete post
	if post.PostLevel == "-1000" {
		// delete post callback
		PostDeletedSubscriberCB(subj, reply, p)

		query := bson.M{"created.profile_id": postCreator.Id, "deleted": true}
		deleteCount := mongo.GetPostCountByQuery(query)
		if deleteCount == 1 || deleteCount == 2 {
			// Warn the user
			reason = post.Reason.IgnoreFeedReason
			if reason == "" {
				reason = post.IgnoreReason
			}
			mongo.UpdateUser(bson.M{
				"profiles._id": postCreator.Id,
			}, bson.M{
				"$set": bson.M{
					"blacklist.status":         "warning",
					"blacklist.last_warned_on": time.Now(),
					"blacklist.reason":         reason,
				},
				"$inc": bson.M{"blacklist.warn_count": 1},
			})
			notification.Purpose = "user.warned"
			notification.Identifier = post.Id.Hex() + "_user_warned"
			send_notification = true
		} else if deleteCount%3 == 0 {
			// block for 2 ^ deleteCount/3 days
			// manch:D, namespace & purpose
			days := deleteCount / 3
			blockForDays := math.Pow(float64(2), float64(days))
			blockTill := time.Now().Local().Add(time.Hour * 24 * time.Duration(int64(blockForDays)))

			reason = post.Reason.DeleteReason

			if reason == "" {
				reason = post.Reason.IgnoreFeedReason
			}

			mongo.UpdateUser(bson.M{
				"profiles._id": postCreator.Id,
			}, bson.M{
				"$set": bson.M{
					"blacklist.status":       "blocked",
					"blacklist.blocked_on":   time.Now(),
					"blacklist.blocked_till": blockTill,
					"blacklist.reason":       reason,
				},
			})
			blockTillString := strconv.FormatInt(blockTill.Unix(), 10)
			blockOnString := strconv.FormatInt(time.Now().Unix(), 10)

			blockedStatus["status"] = "blocked"
			blockedStatus["blocked_on"] = blockOnString
			blockedStatus["blocked_till"] = blockTillString
			send_notification = true
		}
	}

	// warning for 3rd & 4th ignore post
	// block on every 5th ignore for 2^i days

	if post.PostLevel == "-2" {
		// ignore from feed callback
		PostRemovedSubscriberCB(subj, reply, p)

		query := bson.M{"created.profile_id": postCreator.Id, "ignore_from_feed": true, "deleted": false}
		ignoreCount := mongo.GetPostCountByQuery(query)
		if ignoreCount == 3 || ignoreCount == 4 {
			reason = post.Reason.IgnoreFeedReason
			if reason == "" {
				reason = post.Reason.DeleteReason
			}
			// Warn the user
			mongo.UpdateUser(bson.M{
				"profiles._id": postCreator.Id,
			}, bson.M{
				"$set": bson.M{
					"blacklist.status":         "warning",
					"blacklist.last_warned_on": time.Now(),
					"blacklist.reason":         reason,
				},
				"$inc": bson.M{"blacklist.warn_count": 1},
			})
			notification.Purpose = "user.warned"
			notification.Identifier = post.Id.Hex() + ""
			send_notification = true
		} else if ignoreCount%5 == 0 {
			// block for 2 ^ ignoreCount / 5 days
			days := ignoreCount / 5
			blockForDays := math.Pow(float64(2), float64(days))
			blockTill := time.Now().Local().Add(time.Hour * 24 * time.Duration(blockForDays))
			reason = post.Reason.DeleteReason
			if reason == "" {
				reason = post.Reason.IgnoreFeedReason
			}
			mongo.UpdateUser(bson.M{
				"profiles._id": postCreator.Id,
			}, bson.M{
				"$set": bson.M{
					"blacklist.status":       "blocked",
					"blacklist.blocked_on":   time.Now(),
					"blacklist.blocked_till": blockTill,
					"blacklist.reason":       reason,
				},
			})

			blockTillString := strconv.FormatInt(blockTill.Unix(), 10)
			blockOnString := strconv.FormatInt(time.Now().Unix(), 10)

			blockedStatus["status"] = "blocked"
			blockedStatus["blocked_on"] = blockOnString
			blockedStatus["blocked_till"] = blockTillString
			send_notification = true
		}
	}

	if send_notification {
		tokens := mongo.GetTokensByProfiles([]bson.ObjectId{post.Created.ProfileId})
		msg := firebase.ManchMessage{
			Title:     "",
			Message:   "",
			Namespace: "manch:D",
			Id:        notification.NId,
			Reason:    reason,
		}

		if _, ok := blockedStatus["status"]; ok {
			msg.Status = blockedStatus["status"]
			msg.BlockedTill = blockedStatus["blocked_till"]
			msg.BlockedOn = blockedStatus["blocked_on"]
		}
		if tokens != nil {
			for _, token := range tokens {
				go firebase.SendMessage(msg, token.Token, notification)
			}
		} else {
			fmt.Printf("No token")
		}
	}
	fmt.Printf("Processed a post on subject %s! with Post ID %s\n", subj, p.Id)
}
