package callbacks

import (
	"fmt"
	"notification-service/pkg/i18n"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"
	"notification-service/pkg/utils"
	"time"

	"github.com/globalsign/mgo/bson"
)

func PostModeratedSubscriberCB(subj, reply string, p *subscribers.Post) {
	fmt.Printf("Received a post on subject %s! with Post %+v\n", subj, p)
	err, post := mongo.GetPostById(p.Id)
	if err != nil {
		return
	}

	// i represents no of profiles
	n, botProfilesIds := mongo.GetBotProfilesIds()
	randomBotIndex := utils.Random(0, n)

	postCreator := mongo.GetProfileById(post.Created.ProfileId)
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
	fmt.Printf("Processed a post on subject %s! with Post ID %s\n", subj, p.Id)
}
