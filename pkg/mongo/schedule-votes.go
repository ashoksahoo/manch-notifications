package mongo

import (
	"fmt"
	"time"

	"github.com/globalsign/mgo/bson"
)

type VoteScheduleModelPost struct {
	Id           bson.ObjectId `json:"_id" bson:"_id"`
	ResourceType string        `json:"resource_type" bson:"resource_type"`
	Resource     bson.ObjectId `json:"resource_id" bson:"resource_id"`
	Created      Creator       `json:"created" bson:"created"`
	User         bson.ObjectId `json:"user" bson:"user"`
	Value        int           `json:"vote" bson:"vote"`
	Schedule     Schedule      `json:"schedule" bson:"schedule"`
	CreatedAt    time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt" bson:"updatedAt"`
	__v          int           `json:"__v" bson:"__v"`
}

func CreateVotesSchedulePost(scheduleTime time.Time, rId bson.ObjectId, userId bson.ObjectId) VoteScheduleModelPost {

	user := GetUserByProfileId(userId.Hex())
	currentTime := time.Now()
	//creator
	c := Creator{
		Id:        bson.NewObjectId(),
		ProfileId: user.Id,
		Name:      user.Name,
		Avatar:    user.Profiles[0].Avatar,
		UserType:  user.UserType,
	}

	s := Schedule{
		Id:           bson.NewObjectId(),
		Scheduletime: scheduleTime,
		Created:      c,
	}

	return VoteScheduleModelPost{
		Id:           bson.NewObjectId(),
		ResourceType: "post",
		Resource:     rId,
		Created:      c,
		Schedule:     s,
		Value:        1,
		User:         user.Id,
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
		__v:          0,
	}
}

func AddVoteSchedule(document VoteScheduleModelPost) {
	s := session.Clone()
	defer s.Close()
	F := s.DB("manch").C("votes_scheduleds")
	err := F.Insert(document)
	if err == nil {
		fmt.Println("inserted: ", document)
	} else {
		fmt.Println("unable to add vote schedule:", document.Resource.Hex())
	}
}
