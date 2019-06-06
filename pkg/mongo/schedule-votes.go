package mongo

import (
	"notification-service/pkg/constants"
	"fmt"
	"time"

	"github.com/globalsign/mgo/bson"
)

var (
	VOTE_SCHEDULEDS_MODEL = constants.ModelNames["VOTE_SCHEDULEDS"]
)

type VoteScheduleModelPost struct {
	Id           bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	ResourceType string        `json:"resource_type" bson:"resource_type"`
	Resource     bson.ObjectId `json:"resource" bson:"resource"`
	Created      Creator       `json:"created" bson:"created"`
	User         bson.ObjectId `json:"user" bson:"user"`
	Value        int           `json:"vote" bson:"vote"`
	Schedule     Schedule      `json:"schedule" bson:"schedule"`
	CreatedAt    time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt" bson:"updatedAt"`
	__v          int           `json:"__v" bson:"__v"`
}

func CreateVotesSchedulePost(scheduleTime time.Time, rId bson.ObjectId, userProfileId bson.ObjectId) VoteScheduleModelPost {

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

	return VoteScheduleModelPost{
		Id:           bson.NewObjectId(),
		ResourceType: "post",
		Resource:     rId,
		Created:      c,
		Schedule:     s,
		Value:        1,
		User:         userProfileId,
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
		__v:          0,
	}
}

func AddVoteSchedule(document VoteScheduleModelPost) {
	s := session.Clone()
	defer s.Close()
	F := s.DB("manch").C(VOTE_SCHEDULEDS_MODEL)
	err := F.Insert(document)
	if err == nil {
		fmt.Printf("inserted vote schedule: %+v\n", document)
	} else {
		fmt.Println(err)
		fmt.Println("unable to add vote schedule:", document.Resource.Hex())
	}
}

func RemoveVoteScheduleByResource(rId bson.ObjectId) {
	s := session.Clone()
	defer s.Close()
	V := s.DB("manch").C(VOTE_SCHEDULEDS_MODEL)
	info, err := V.RemoveAll(bson.M{"resource": rId})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("removed info", info)
	}
}
