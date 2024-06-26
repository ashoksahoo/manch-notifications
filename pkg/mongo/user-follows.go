package mongo

import (
	"notification-service/pkg/constants"
	"github.com/globalsign/mgo/bson"
)

var (
	USER_FOLLOWS_MODEL = constants.ModelNames["USER_FOLLOWS"]
)

type UserFollowModel struct {
	Id           bson.ObjectId `json:"_id" bson:"_id"`
	ProfileId    bson.ObjectId `json:"profile_id" bson:"profile_id"`
	ResourceId   bson.ObjectId `json:"resource_id" bson:"resource_id"`
	ResourceType string        `json:"resource_type" bson:"resource_type"`
	Created      Creator       `json:"created" bson:"created"`
}

func GetUserFollowById(Id string) UserFollowModel {
	s := session.Clone()
	defer s.Close()
	user_follow := UserFollowModel{}
	UF := s.DB("manch").C(USER_FOLLOWS_MODEL)
	UF.Find(bson.M{"_id": bson.ObjectIdHex(Id)}).One(&user_follow)
	return user_follow
}
