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
	Value        int           `json:"vote" bson:"vote"`
	Schedule     Schedule      `json:"schedule" bson:"schedule"`
}

func CreateVotesSchedulePost(scheduleTime time.Time, rId bson.ObjectId, userId bson.ObjectId) VoteScheduleModelPost {

	user := GetUserByProfileId(userId.Hex())
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
