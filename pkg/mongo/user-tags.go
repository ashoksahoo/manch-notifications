package mongo

import (
	"fmt"
	"notification-service/pkg/constants"
	"time"

	"github.com/globalsign/mgo/bson"
)

var (
	USER_TAGS_MODEL = constants.ModelNames["USER_TAGS"]
)

type UserTags struct {
	ID        bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	ProfileID bson.ObjectId `json:"profile_id" bson:"profile_id"`
	UserType  string        `json:"user_type" bson:"user_type"`
	Tag       string        `json:"tag" bson:"tag"`
	Score     int           `json:"score" bson:"score"`
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt" bson:"updatedAt"`
}

func CreateUserTags(post PostModel) {
	s := session.Clone()
	defer s.Close()
	UT := s.DB("manch").C(USER_TAGS_MODEL)
	profileId := post.Created.ProfileId
	currentTime := time.Now()
	userTag := UserTags{
		ProfileID: profileId,
		UserType:  post.Created.UserType,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Score:     1,
	}
	for _, tag := range post.Tags {
		userTag.Tag = tag
		count, _ := UT.Find(bson.M{"profile_id": profileId, "tag": tag}).Count()
		fmt.Println("count", count)
		if count > 0 {
			// update
			UT.Update(bson.M{"profile_id": profileId, "tag": tag}, bson.M{
				"$inc": bson.M{"score": 1},
				"$set": bson.M{"updatedAt": time.Now()},
			})
		} else {
			fmt.Printf("inserting user score %+v\n", userTag)
			UT.Insert(userTag)
		}
	}
}

func GetUserTagsById(id bson.ObjectId) UserTags {
	s := session.Clone()
	defer s.Close()
	US := s.DB("manch").C(USER_TAGS_MODEL)
	userTag := UserTags{}
	US.Find(bson.M{"_id": id}).One(&userTag)
	return userTag
}
