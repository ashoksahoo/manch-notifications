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

	var noOfVotes, noOfShares int

	if post.Created.UserType == "bot" {
		noOfVotes = utils.Random(40, 45)
		noOfShares = int((noOfVotes * utils.Random(50, 80)) / 100)

		// vote schedule
		voteIndex := 0
		t := utils.SplitTimeInRange(1, 30, noOfVotes, time.Minute)
		for k := 0; voteIndex < noOfVotes; voteIndex, k = voteIndex+1, k+1 {
			vote := mongo.CreateVotesSchedulePost(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[voteIndex]))
			mongo.AddVoteSchedule(vote)
		}
		// share schedule
		shareIndex := 0
		t = utils.SplitTimeInRange(1, 30, noOfShares, time.Minute)
		for k := 0; shareIndex < noOfShares; shareIndex, k = shareIndex+1, k+1 {
			share := mongo.CreateShareSchedule(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[shareIndex]))
			mongo.AddShareSchedule(share)
		}
		if post.PostType != "VIDEO" {
			randomVotes := utils.Random(20, 30)
			noOfVotes += randomVotes
			t = utils.SplitTimeInRange(31, 90, randomVotes, time.Minute)
			for k := 0; voteIndex < noOfVotes; voteIndex, k = voteIndex+1, k+1 {
				vote := mongo.CreateVotesSchedulePost(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[voteIndex]))
				mongo.AddVoteSchedule(vote)
			}

			// share schedule
			randomShare := int((randomVotes * utils.Random(50, 80)) / 100)
			noOfShares += randomShare
			t = utils.SplitTimeInRange(31, 90, randomShare, time.Minute)
			for k := 0; shareIndex < noOfShares; shareIndex, k = shareIndex+1, k+1 {
				share := mongo.CreateShareSchedule(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[shareIndex]))
				mongo.AddShareSchedule(share)
			}

		} else {
			randomVotes := utils.Random(25, 36)
			noOfVotes += randomVotes
			t = utils.SplitTimeInRange(31, 120, randomVotes, time.Minute)
			for k := 0; voteIndex < noOfVotes; voteIndex, k = voteIndex+1, k+1 {
				vote := mongo.CreateVotesSchedulePost(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[voteIndex]))
				mongo.AddVoteSchedule(vote)
			}

			// share schedule
			randomShare := int((randomVotes * utils.Random(50, 80)) / 100)
			noOfShares += randomShare
			t = utils.SplitTimeInRange(31, 120, randomShare, time.Minute)
			for k := 0; shareIndex < noOfShares; shareIndex, k = shareIndex+1, k+1 {
				share := mongo.CreateShareSchedule(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[shareIndex]))
				mongo.AddShareSchedule(share)
			}

		}
	} else if post.PostLevel == "2" {
		noOfVotes = utils.Random(40, 45)
		noOfShares = int((noOfVotes * utils.Random(50, 80)) / 100)

		voteIndex := 0
		t := utils.SplitTimeInRange(1, 90, noOfVotes, time.Minute)
		for k := 0; voteIndex < noOfVotes; voteIndex, k = voteIndex+1, k+1 {
			vote := mongo.CreateVotesSchedulePost(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[voteIndex]))
			mongo.AddVoteSchedule(vote)
		}

		shareIndex := 0
		t = utils.SplitTimeInRange(1, 90, noOfShares, time.Minute)
		for k := 0; shareIndex < noOfShares; shareIndex, k = shareIndex+1, k+1 {
			share := mongo.CreateShareSchedule(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[shareIndex]))
			mongo.AddShareSchedule(share)
		}

		if post.PostType != "VIDEO" {
			randomVotes := utils.Random(20, 30)
			noOfVotes += randomVotes
			t = utils.SplitTimeInRange(91, 150, randomVotes, time.Minute)
			for k := 0; voteIndex < noOfVotes; voteIndex, k = voteIndex+1, k+1 {
				vote := mongo.CreateVotesSchedulePost(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[voteIndex]))
				mongo.AddVoteSchedule(vote)
			}

			// share schedule
			randomShare := int((randomVotes * utils.Random(50, 80)) / 100)
			noOfShares += randomShare
			t = utils.SplitTimeInRange(91, 150, randomShare, time.Minute)
			for k := 0; shareIndex < noOfShares; shareIndex, k = shareIndex+1, k+1 {
				share := mongo.CreateShareSchedule(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[shareIndex]))
				mongo.AddShareSchedule(share)
			}
		} else {
			randomVotes := utils.Random(30, 45)
			noOfVotes += randomVotes
			t = utils.SplitTimeInRange(91, 210, randomVotes, time.Minute)
			for k := 0; voteIndex < noOfVotes; voteIndex, k = voteIndex+1, k+1 {
				vote := mongo.CreateVotesSchedulePost(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[voteIndex]))
				mongo.AddVoteSchedule(vote)
			}
			randomShare := int((randomVotes * utils.Random(50, 80)) / 100)
			noOfShares += randomShare
			t = utils.SplitTimeInRange(91, 210, randomShare, time.Minute)
			for k := 0; shareIndex < noOfShares; shareIndex, k = shareIndex+1, k+1 {
				share := mongo.CreateShareSchedule(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[shareIndex]))
				mongo.AddShareSchedule(share)
			}
		}
	} else if post.PostLevel == "1" {
		noOfVotes = utils.Random(5, 10)
		j := 0
		t := utils.SplitTimeInRange(1, 90, noOfVotes, time.Minute)
		for k := 0; j < noOfVotes; j, k = j+1, k+1 {
			vote := mongo.CreateVotesSchedulePost(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[j]))
			mongo.AddVoteSchedule(vote)
		}
	} else if post.PostLevel == "0" {
		noOfVotes = utils.Random(0, 3)
		j := 0
		t := utils.SplitTimeInRange(1, 90, noOfVotes, time.Minute)
		for k := 0; j < noOfVotes; j, k = j+1, k+1 {
			vote := mongo.CreateVotesSchedulePost(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[j]))
			mongo.AddVoteSchedule(vote)
		}
	}

	fmt.Println("No. of Vote added", noOfVotes)
	fmt.Println("No. of Share added", noOfShares)
	fmt.Printf("Processed a post on subject %s! with Post ID %s\n", subj, p.Id)
}
