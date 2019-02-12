package mongo

import (
	"notification-service/pkg/constants"
	"fmt"
	"time"

	"github.com/globalsign/mgo/bson"
)

var (
	SHARES_SCHEDULEDS_MODEL = constants.ModelNames["SHARES_SCHEDULEDS"]
)
type ShareModel struct {
	Id           bson.ObjectId `json:"_id" bson:"_id"`
	ResourceId   bson.ObjectId `json:"resource_id" bson:"resource_id"`
	ResourceType string        `json:"resource_type" bson:"resource_type"`
	ProfileId    bson.ObjectId `json:"profile_id" bson:"profile_id"`
	Created      Creator       `json:"created" bson:"created"`
}

type ShareScheduleModel struct {
	Id           bson.ObjectId `json:"_id" bson:"_id"`
	ResourceId   bson.ObjectId `json:"resource_id" bson:"resource_id"`
	ResourceType string        `json:"resource_type" bson:"resource_type"`
	ProfileId    bson.ObjectId `json:"profile_id" bson:"profile_id"`
	Created      Creator       `json:"created" bson:"created"`
	Schedule     Schedule      `json:"schedule" bson:"schedule"`
	CreatedAt    time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt" bson:"updatedAt"`
	__v          int           `json:"__v" bson:"__v"`
}

func CreateShareSchedule(scheduleTime time.Time, rId bson.ObjectId, userProfileId bson.ObjectId) ShareScheduleModel {

	user := GetUserByProfileId(userProfileId.Hex())
	currentTime := time.Now()
	//creator
	c := Creator{
		Id:        bson.NewObjectId(),
		ProfileId: userProfileId,
		Name:      user.Name,
		Avatar:    user.Profiles[0].Avatar,
		UserType:  user.UserType,
	}

	s := Schedule{
		Id:           bson.NewObjectId(),
		Scheduletime: scheduleTime,
		Created:      c,
	}

	return ShareScheduleModel{
		Id:           bson.NewObjectId(),
		ResourceId:   rId,
		ResourceType: "post",
		ProfileId:    userProfileId,
		Schedule:     s,
		Created:      c,
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
		__v:          0,
	}
}

func AddShareSchedule(document ShareScheduleModel) {
	s := session.Clone()
	defer s.Close()
	F := s.DB("manch").C(SHARES_SCHEDULEDS_MODEL)
	err := F.Insert(document)
	if err == nil {
		fmt.Printf("inserted share schedule: %+v\n", document)
	} else {
		fmt.Println(err)
		fmt.Println("unable to add share schedule:", document.ResourceId.Hex())
	}
}

func RemoveShareScheduleByResource(rId bson.ObjectId) {
	s := session.Clone()
	defer s.Close()
	F := s.DB("manch").C(SHARES_SCHEDULEDS_MODEL)
	info, err := F.RemoveAll(bson.M{"resource_id": rId})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("removed info", info)
	}
}
