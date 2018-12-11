package mongo

import (
	"github.com/globalsign/mgo/bson"
	"fmt"
	"time"
)

type Schedule struct {
	Id bson.ObjectId `json:"_id" bson:"_id"`
	Scheduletime time.Time `json:"schedule_time" bson:"schedule_time"`
	Created      Creator         `json:"created" bson:"created"`
}

type UserFollowScheduleModel struct {
	Id bson.ObjectId `json:"_id" bson:"_id"`
	ProfileId bson.ObjectId `json:"profile_id" bson:"profile_id"`
	ResourceId     bson.ObjectId `json:"resource_id" bson:"resource_id"`
	ResourceType string `json:"resource_type" bson:"resource_type"`
	Created Creator `json:"created" bson:"created"`
	Schedule Schedule `json:"schedule" bson:"schedule"`
}

func AddFollowSchedule(document UserFollowScheduleModel) {
	s := session.Clone()
	defer s.Close()
	F := s.DB("manch").C("user_follows_scheduleds")
	// fmt.Println("inserting document:", document)
	err := F.Insert(document)
	fmt.Println("err", err)
}


func CreateFollowSchedule(scheduleTime time.Time, pId, rId bson.ObjectId, c Creator) UserFollowScheduleModel {
	
	s := Schedule {
		Id:           bson.NewObjectId(),
		Scheduletime: scheduleTime,
		Created:      c,
	}
	
	return UserFollowScheduleModel{
		Id:           bson.NewObjectId(),
		ProfileId:    pId,
		ResourceId:   rId,
		ResourceType: "user",
		Created: c,
		Schedule:     s,
	}
}
