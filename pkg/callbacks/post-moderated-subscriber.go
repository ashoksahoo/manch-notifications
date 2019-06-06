package callbacks

import (
	"fmt"
	"math/rand"
	"notification-service/pkg/elasticsearch"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"
	"notification-service/pkg/utils"
	"time"

	"github.com/globalsign/mgo/bson"
)

func PostModeratedSubscriberCB(subj, reply string, p *subscribers.Post) {
	fmt.Printf("Received a post on subject %s! with Post %+v\n", subj, p)
	err, post := mongo.GetPostByQuery(bson.M{"_id": bson.ObjectIdHex(p.Id)})
	if err != nil {
		return
	}

	if post.PostLevel == "-1000" {
		// delete post callback
		PostDeletedSubscriberCB(subj, reply, p)
		return
	}

	if post.PostLevel == "-2" {
		// ignore from feed callback
		PostRemovedSubscriberCB(subj, reply, p)
		return
	}
	if post.PostLevel == "-1" {
		return
	}
	// process hashtags
	image := mongo.ExtractThumbNailFromPost(post)
	elasticsearch.AddTagToIndex(post.Tags, image)

	// create or update user hashtags
	mongo.CreateUserTags(post)

	m, botProfilesHi := mongo.GetBotProfilesIds("hi")
	n, botProfilesTe := mongo.GetBotProfilesIds("te")
	n = m + n
	botProfilesIds := append(botProfilesHi, botProfilesTe...)
	// shuffle profiles
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(n, func(i, j int) { botProfilesIds[i], botProfilesIds[j] = botProfilesIds[j], botProfilesIds[i] })

	var no_of_votes int
	if post.PostLevel == "2" {
		no_of_votes = utils.Random(40, 50)
		if post.PostType == "VIDEO" {
			no_of_votes = utils.Random(80, 100)
		}
	} else if post.PostLevel == "1" {
		no_of_votes = utils.Random(5, 10)
	} else if post.PostLevel == "0" {
		no_of_votes = utils.Random(0, 3)
	}
	j := 0
	fmt.Println("no_of_votes: ", no_of_votes)
	t := utils.SplitTimeInRange(1, 60*24, no_of_votes, time.Minute)
	for k := 0; j < no_of_votes; j, k = j+1, k+1 {
		vote := mongo.CreateVotesSchedulePost(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[j]))
		mongo.AddVoteSchedule(vote)
	}

	// randomVotes := utils.Random(5, 20)
	// no_of_votes += randomVotes
	// fmt.Println("no_of_votes: ", no_of_votes)
	// t = utils.SplitTimeInRange(30, 2*24*60, randomVotes, time.Minute)
	// for k := 0; j < no_of_votes; j, k = j+1, k+1 {
	// 	vote := mongo.CreateVotesSchedulePost(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[j]))
	// 	mongo.AddVoteSchedule(vote)
	// }

	// // schedule shares on posts
	// var no_of_shares int
	// no_of_shares = utils.Random(5, 10)

	// j = 0
	// fmt.Println("no_of_shares: ", no_of_shares)
	// t = utils.SplitTimeInRange(1, 240, no_of_shares, time.Minute)
	// for k := 0; j < no_of_shares; j, k = j+1, k+1 {
	// 	share := mongo.CreateShareSchedule(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[j]))
	// 	mongo.AddShareSchedule(share)
	// }

	fmt.Printf("Processed a post on subject %s! with Post ID %s\n", subj, p.Id)
}
