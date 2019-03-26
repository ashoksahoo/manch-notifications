package callbacks

import (
	"fmt"
	"math/rand"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"
	"notification-service/pkg/utils"
	"time"

	"github.com/globalsign/mgo/bson"
)

func PostSubscriberCB(subj, reply string, p *subscribers.Post) {
	fmt.Printf("Received a post on subject %s! with Post %+v\n", subj, p)

	m, botProfilesHi := mongo.GetBotProfilesIds("hi")
	n, botProfilesTe := mongo.GetBotProfilesIds("te")
	n = m + n
	botProfilesIds := append(botProfilesHi, botProfilesTe...)
	// shuffle profiles
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(n, func(i, j int) { botProfilesIds[i], botProfilesIds[j] = botProfilesIds[j], botProfilesIds[i] })

	// update user score for new post
	_, post := mongo.GetPostById(p.Id)

	// create community stats
	mongo.CreateCommunityStats(mongo.CommunityStatsModel{
		CommunityId: post.CommunityIds[0],
		Action:      "post",
		EntityId:    post.Id,
		EntityType:  "post",
		ProfileId:   post.Created.ProfileId,
		PostsCount:  1,
	})

	mongo.CreateUserScore(mongo.UserScore{
		ProfileId:   post.Created.ProfileId,
		CommunityId: post.CommunityIds[0],
		Score:       1,
		UserType:    post.Created.UserType,
	})

	var no_of_votes int
	if p.CreatorType == "bot" {
		no_of_votes = utils.Random(20, 25)
	} else if p.CreatorType == "super_level_1" {
		no_of_votes = utils.Random(20, 25)
	} else if p.CreatorType == "verified_level_1" {
		no_of_votes = utils.Random(20, 25)
	} else {
		no_of_votes = utils.Random(5, 10)
	}

	j := 0
	fmt.Println("no_of_votes: ", no_of_votes)
	t := utils.SplitTimeInRange(1*60, 30*60, no_of_votes, time.Second)
	for k := 0; j < no_of_votes; j, k = j+1, k+1 {
		vote := mongo.CreateVotesSchedulePost(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[j]))
		mongo.AddVoteSchedule(vote)
	}

	randomVotes := utils.Random(5, 20)
	no_of_votes += randomVotes
	fmt.Println("no_of_votes: ", no_of_votes)
	t = utils.SplitTimeInRange(30, 2*24*60, randomVotes, time.Minute)
	for k := 0; j < no_of_votes; j, k = j+1, k+1 {
		vote := mongo.CreateVotesSchedulePost(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[j]))
		mongo.AddVoteSchedule(vote)
	}

	// schedule shares on posts
	var no_of_shares int
	no_of_shares = utils.Random(15, 40)

	j = 0
	fmt.Println("no_of_shares: ", no_of_shares)
	t = utils.SplitTimeInRange(1, 240, no_of_shares, time.Minute)
	for k := 0; j < no_of_shares; j, k = j+1, k+1 {
		share := mongo.CreateShareSchedule(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[j]))
		mongo.AddShareSchedule(share)
	}

	fmt.Printf("Add %d share scheduleds\n", no_of_shares)

	fmt.Printf("Processed a post on subject %s! with Post Id%s\n", subj, p.Id)
}
