package callbacks

import (
	"github.com/globalsign/mgo/bson"
	"notification-service/pkg/utils"
	"fmt"
	"math/rand"
	"notification-service/pkg/i18n"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"
	"time"
)

const (
	MANCH_OFFICIAL_TE = "5c3c3bfd89ac4a794d45b14d"
	MANCH_OFFICIAL_HE = "5c1c92c8eda9bd1771bcf0a7"
)

func PostModeratedSubscriberCB(subj, reply string, p *subscribers.Post) {
	fmt.Printf("Received a post on subject %s! with Post %+v\n", subj, p)
	err, post := mongo.GetPostById(p.Id)
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
			return;
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
		commentCreator := mongo.Creator {
			Id:        bson.NewObjectId(),
			ProfileId: commentator.Id,
			Name:      commentator.Name,
			Avatar:    commentator.Avatar,
			UserType:  commentator.Type,
		}
		randomMinute := utils.Random(15, 30)
		fmt.Println("random Minute", randomMinute)
		scheduleTime := time.Now().Local().Add(time.Minute*time.Duration(randomMinute))
		fmt.Println("schedule time", scheduleTime)			
		// schedule comments
		mongo.CreateCommentSchedule(comment, post.Id, commentCreator, scheduleTime)
		mongo.AddCommentStringToProfileId(postCreator.Id, keys[randomCommentKeyIndex])
		// schedule comments in 20-30 minutes random
	}
	fmt.Printf("Processed a post on subject %s! with Post ID %s\n", subj, p.Id)
}
