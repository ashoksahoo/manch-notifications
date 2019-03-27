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

func RepostSubscriberCB(subj, reply string, p *subscribers.Post) {
	fmt.Printf("Received a post on subject %s! with Post %+v\n", subj, p)

	err, post := mongo.GetPostById(p.Id)
	if err != nil {
		return
	}

	m, botProfilesHi := mongo.GetBotProfilesIds("hi")
	n, botProfilesTe := mongo.GetBotProfilesIds("te")
	n = m + n
	botProfilesIds := append(botProfilesHi, botProfilesTe...)
	// shuffle profiles
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(n, func(i, j int) { botProfilesIds[i], botProfilesIds[j] = botProfilesIds[j], botProfilesIds[i] })

	if post.PostLevel == "2" && post.UpVotes < 25 {
		// schedule votes to reach 30 votes
		fmt.Println("General user....")
		diff := 25 - post.UpVotes
		fmt.Println("total difference is", diff)
		j := 0
		// schedule 5 votes in 30 minutes
		total_votes := diff
		if diff > 0 {
			diff -= 5
			t := utils.SplitTimeInRange(1*60, 30*60, 5, time.Second)
			for k := 0; j < 5; j, k = j+1, k+1 {
				vote := mongo.CreateVotesSchedulePost(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[j]))
				mongo.AddVoteSchedule(vote)
			}
		}
		// schedule rest votes in next 1.30 hours
		if diff > 0 {
			t := utils.SplitTimeInRange(30, 120, diff, time.Minute)
			for k := 0; j < total_votes; j, k = j+1, k+1 {
				vote := mongo.CreateVotesSchedulePost(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[j]))
				mongo.AddVoteSchedule(vote)
			}
		}
		fmt.Println("total votes added: ", j)

	} else if post.PostLevel == "1" && post.UpVotes < 20 {
		// schedule votes to reach 30 votes
		fmt.Println("General user....")
		diff := 20 - post.UpVotes
		fmt.Println("total difference is", diff)
		j := 0
		// schedule 5 votes in 30 minutes
		total_votes := diff
		if diff > 0 {
			diff -= 5
			t := utils.SplitTimeInRange(1*60, 30*60, 5, time.Second)
			for k := 0; j < 5; j, k = j+1, k+1 {
				vote := mongo.CreateVotesSchedulePost(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[j]))
				mongo.AddVoteSchedule(vote)
			}
		}
		// schedule rest votes in next 1.30 hours
		if diff > 0 {
			t := utils.SplitTimeInRange(30, 120, diff, time.Minute)
			for k := 0; j < total_votes; j, k = j+1, k+1 {
				vote := mongo.CreateVotesSchedulePost(t[k], bson.ObjectIdHex(p.Id), bson.ObjectIdHex(botProfilesIds[j]))
				mongo.AddVoteSchedule(vote)
			}
		}
		fmt.Println("total votes added: ", j)

	}
	fmt.Printf("Processed a post on subject %s! with Post Id %s\n", subj, p.Id)

}
