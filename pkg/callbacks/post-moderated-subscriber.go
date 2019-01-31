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
	"time"

	"github.com/globalsign/mgo/bson"
)

const (
	MANCH_OFFICIAL_TE = "5c3c3bfd89ac4a794d45b14d"
	MANCH_OFFICIAL_HE = "5c1c92c8eda9bd1771bcf0a7"
)

func PostModeratedSubscriberCB(subj, reply string, p *subscribers.Post) {
	fmt.Printf("Received a post on subject %s! with Post %+v\n", subj, p)
	err, post := mongo.GetPostByQuery(bson.M{"_id": bson.ObjectIdHex(p.Id)})
	if err != nil {
		return
	}
	// get all bot users
	botUsers := mongo.GetBotUsers()
	// array of bot profiles ids
	var botProfilesIds [100]string
	// no. of profiles counter
	i := 0
	for _, botUser := range botUsers {
		profiles := botUser.Profiles
		for _, profile := range profiles {
			if i == 100 {
				break
			}
			if profile.Id.Hex() == MANCH_OFFICIAL_HE || profile.Id.Hex() == MANCH_OFFICIAL_TE {
				continue
			}
			botProfilesIds[i] = profile.Id.Hex()
			i++
		}
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(i, func(i, j int) { botProfilesIds[i], botProfilesIds[j] = botProfilesIds[j], botProfilesIds[i] })

	randomBotIndex := utils.Random(0, i)

	postCreator := mongo.GetProfileById(post.Created.ProfileId)


	// schedule auto comment on post if it is good
	if post.PostLevel == "2" || post.PostLevel == "1" {
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
		scheduleTime := time.Now().Local().Add(time.Minute * time.Duration(randomMinute))
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
	// warning for 1st and 2nd delete post
	// block on every 3rd delete post
	if post.PostLevel == "-1000" {
		// delete post callback
		PostDeletedSubscriberCB(subj, reply, p)

		query := bson.M{"created.profile_id": postCreator.Id, "deleted": true}
		deleteCount := mongo.GetPostCountByQuery(query)
		if deleteCount == 1 || deleteCount == 2 {
			// Warn the user
			mongo.UpdateUser(bson.M{
				"profiles._id": postCreator.Id,
			}, bson.M{
				"$set": bson.M{
					"blacklist.status":         "warning",
					"blacklist.last_warned_on": time.Now(),
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
			mongo.UpdateUser(bson.M{
				"profiles._id": postCreator.Id,
			}, bson.M{
				"$set": bson.M{
					"blacklist.status":       "blocked",
					"blacklist.blocked_on":   time.Now(),
					"blacklist.blocked_till": blockTill,
				},
			})
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
			// Warn the user
			mongo.UpdateUser(bson.M{
				"profiles._id": postCreator.Id,
			}, bson.M{
				"$set": bson.M{
					"blacklist.status":         "warning",
					"blacklist.last_warned_on": time.Now(),
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
			mongo.UpdateUser(bson.M{
				"profiles._id": postCreator.Id,
			}, bson.M{
				"$set": bson.M{
					"blacklist.status":       "blocked",
					"blacklist.blocked_on":   time.Now(),
					"blacklist.blocked_till": blockTill,
				},
			})
			send_notification = true
		}
	}

	if send_notification {
		tokens := mongo.GetTokensByProfiles([]bson.ObjectId{post.Created.ProfileId})
		msg := firebase.ManchMessage{
			Title:   "",
			Message: "",
			Namespace: "manch:D",
			Id:      notification.NId,
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
