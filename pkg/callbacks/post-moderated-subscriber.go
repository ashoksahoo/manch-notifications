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
	var additionalScore int
	if post.Created.UserType == "bot" {
		additionalScore = 50 * 60
	} else {
		additionalScore = 5 * 60
	}
	elasticsearch.AddTagToIndex(post.Tags, additionalScore, post.TagsPosition)

	// index post
	elasticsearch.CreatePostIndex(post)

	// create or update user hashtags
	mongo.CreateUserTags(post)

	m, botProfilesHi := mongo.GetBotProfilesIds("hi")
	n, botProfilesTe := mongo.GetBotProfilesIds("te")
	n = m + n
	botProfilesIds := append(botProfilesHi, botProfilesTe...)
	// shuffle profiles
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(n, func(i, j int) { botProfilesIds[i], botProfilesIds[j] = botProfilesIds[j], botProfilesIds[i] })

	var noOfVotes int

	scheduledVoteCount := mongo.CountScheduledVotesByQuery(bson.M{"resource": post.Id})
	if scheduledVoteCount > 0 {
		return
	}

	if post.PostLevel == "2" {
		noOfVotes = utils.Random(30, 60)
		totalTimes := utils.Random(90, 120)
		fmt.Println("random vote", noOfVotes)
		fmt.Println("random Time", totalTimes)
		voteIndex := 0
		t := utils.SplitTimeInRange(1, totalTimes, noOfVotes, time.Minute)
		for k := 0; voteIndex < noOfVotes; voteIndex, k = voteIndex+1, k+1 {
			vote := mongo.CreateVotesSchedulePost(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[voteIndex]))
			mongo.AddVoteSchedule(vote)
		}

		randomVote := utils.Random(40, 80)
		noOfVotes += randomVote

		randomTime := utils.Random(totalTimes + 90, totalTimes + 120)
		fmt.Println("random vote", randomVote)
		fmt.Println("random Time", randomTime)
		t = utils.SplitTimeInRange(totalTimes, randomTime, randomVote, time.Minute)
		totalTimes = randomTime
		for k := 0; voteIndex < noOfVotes; voteIndex, k = voteIndex+1, k+1 {
			vote := mongo.CreateVotesSchedulePost(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[voteIndex]))
			mongo.AddVoteSchedule(vote)
		}

		randomVote = utils.Random(40, 80)
		noOfVotes += randomVote
		randomTime = utils.Random(totalTimes + 160, totalTimes + 200)
		fmt.Println("random vote", randomVote)
		fmt.Println("random Time", randomTime)
		t = utils.SplitTimeInRange(totalTimes, randomTime, randomVote, time.Minute)
		totalTimes = randomTime
		for k := 0; voteIndex < noOfVotes; voteIndex, k = voteIndex+1, k+1 {
			vote := mongo.CreateVotesSchedulePost(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[voteIndex]))
			mongo.AddVoteSchedule(vote)
		}

		randomVote = utils.Random(10, 20)
		noOfVotes += randomVote
		randomTime = utils.Random(totalTimes + 160, totalTimes + 200)
		fmt.Println("random vote", randomVote)
		fmt.Println("random Time", randomTime)
		t = utils.SplitTimeInRange(totalTimes, randomTime, randomVote, time.Minute)
		totalTimes = randomTime
		for k := 0; voteIndex < noOfVotes; voteIndex, k = voteIndex+1, k+1 {
			vote := mongo.CreateVotesSchedulePost(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[voteIndex]))
			mongo.AddVoteSchedule(vote)
		}

		randomVote = utils.Random(0, 20)
		noOfVotes += randomVote
		randomTime = utils.Random(totalTimes + 250, totalTimes + 360)
		fmt.Println("random vote", randomVote)
		fmt.Println("random Time", randomTime)
		t = utils.SplitTimeInRange(totalTimes, randomTime, randomVote, time.Minute)
		totalTimes = randomTime
		for k := 0; voteIndex < noOfVotes; voteIndex, k = voteIndex+1, k+1 {
			vote := mongo.CreateVotesSchedulePost(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[voteIndex]))
			mongo.AddVoteSchedule(vote)
		}
	} else if post.PostLevel == "1" {
		noOfVotes = utils.Random(15, 30)
		randomTime := utils.Random(240, 360)
		j := 0
		t := utils.SplitTimeInRange(1, randomTime, noOfVotes, time.Minute)
		for k := 0; j < noOfVotes; j, k = j+1, k+1 {
			vote := mongo.CreateVotesSchedulePost(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[j]))
			mongo.AddVoteSchedule(vote)
		}
	}

	fmt.Println("No. of Vote added", noOfVotes)
	fmt.Printf("Processed a post on subject %s! with Post ID %s\n", subj, p.Id)
}
