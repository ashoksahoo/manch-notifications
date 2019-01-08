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

	post := mongo.GetPostById(p.Id)

	// i represents no of profiles
	i, botProfilesIds := mongo.GetBotProfilesIds()

	// shuffle profiles
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(i, func(i, j int) { botProfilesIds[i], botProfilesIds[j] = botProfilesIds[j], botProfilesIds[i] })

	if post.Created.UserType == "general" && post.UpVotes < 30 {
		// schedule votes to reach 30 votes
		fmt.Println("General user....")
		diff := 30 - post.UpVotes
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

	} else if post.Created.UserType == "bot" && post.UpVotes < 45 {
		fmt.Println("Bot user....")

		// schedule votes to reach 45 votes
		diff := 45 - post.UpVotes
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