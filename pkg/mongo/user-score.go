package mongo

import (
	"fmt"
	"time"

	"github.com/globalsign/mgo/bson"
)

type UserScore struct {
	Id          bson.ObjectId `json:"_id" bson:"_id"`
	ProfileId   bson.ObjectId `json:"profile_id" bson:"profile_id"`
	CommunityId bson.ObjectId `json:"community_id" bson:"community_id"`
	UserType    string        `json:"user_type" bson:"user_type"`
	Score       int           `json:"score" bson:"score"`
	LastUpdated time.Time     `json:"last_updated_at" bson:"last_updated_at"`
	ProfileCreatedAt time.Time `json:"profile_created_at" bson:"profile_created_at"`
}

func CreateUserScore(userScore UserScore) {
	s := session.Clone()
	defer s.Close()
	US := s.DB("manch").C("user_scores")
	userScore.LastUpdated = time.Now()
	userScore.Id = bson.NewObjectId()

	user := GetUserByProfileId(userScore.ProfileId.Hex())
	userScore.ProfileCreatedAt = user.CreatedAt

	count, _ := US.Find(bson.M{"profile_id": userScore.ProfileId, "community_id": userScore.CommunityId}).Count()
	fmt.Println("count", count)	
	if count > 0 {
		// update
		fmt.Printf("updating user score %+v\n", userScore)
		US.Update(bson.M{"profile_id": userScore.ProfileId, "community_id": userScore.CommunityId}, bson.M{
			"$inc": bson.M{"score": userScore.Score},
		})
	} else {
		fmt.Printf("inserting user score %+v\n", userScore)
		US.Insert(userScore)
	}
}

func GetUserScoreById(id bson.ObjectId) UserScore {
	s := session.Clone()
	defer s.Close()
	US := s.DB("manch").C("user_scores")
	userScore := UserScore{}
	US.Find(bson.M{"_id": id}).One(&userScore)
	return userScore
}
