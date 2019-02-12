package mongo

import (
	"notification-service/pkg/constants"
	"fmt"
	"time"

	"github.com/globalsign/mgo/bson"
)

var (
	USER_FOLLOWS_SCHEDULEDS_MODEL = constants.ModelNames["USER_FOLLOWS_SCHEDULEDS"]
)

type Schedule struct {
	Id           bson.ObjectId `json:"_id" bson:"_id"`
	Scheduletime time.Time     `json:"schedule_time" bson:"schedule_time"`
	Created      Creator       `json:"created" bson:"created"`
}

type UserFollowScheduleModel struct {
	Id           bson.ObjectId `json:"_id" bson:"_id"`
	ProfileId    bson.ObjectId `json:"profile_id" bson:"profile_id"`
	ResourceId   bson.ObjectId `json:"resource_id" bson:"resource_id"`
	ResourceType string        `json:"resource_type" bson:"resource_type"`
	Created      Creator       `json:"created" bson:"created"`
	Schedule     Schedule      `json:"schedule" bson:"schedule"`
}

func AddFollowSchedule(document UserFollowScheduleModel) {
	s := session.Clone()
	defer s.Close()
	F := s.DB("manch").C(USER_FOLLOWS_SCHEDULEDS_MODEL)
	// fmt.Println("inserting document:", document)
	err := F.Insert(document)
	// fmt.Println("err", err)
	if err == nil {
		// fmt.Println("profile added successfully:", document.ProfileId.Hex())
		fmt.Printf("inserted follow schedule: %+v\n", document)
	} else {
		fmt.Println("unable to add profile:", document.ProfileId.Hex())
	}
}

func CreateFollowSchedule(scheduleTime time.Time, pId, rId bson.ObjectId) UserFollowScheduleModel {

	user := GetUserByProfileId(pId.Hex())
	//creator
	c := Creator{
		Id:        bson.NewObjectId(),
		ProfileId: pId,
		Name:      user.Name,
		Avatar:    user.Profiles[0].Avatar,
		UserType:  user.UserType,
	}

	s := Schedule{
		Id:           bson.NewObjectId(),
		Scheduletime: scheduleTime,
		Created:      c,
	}

	return UserFollowScheduleModel{
		Id:           bson.NewObjectId(),
		ProfileId:    pId,
		ResourceId:   rId,
		ResourceType: "user",
		Created:      c,
		Schedule:     s,
	}
}
