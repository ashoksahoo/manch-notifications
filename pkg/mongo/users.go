package mongo

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
)

type Profile struct {
	Id       bson.ObjectId `json:"_id" bson:"_id"`
	Avatar   string        `json:"avatar"`
	Name     string        `json:"name"`
	Language string        `json:"language"`
}

type Creator struct {
	Id        bson.ObjectId `json:"_id" bson:"_id"`
	ProfileId bson.ObjectId `json:"profile_id" bson:"profile_id"`
	Name      string        `json:"name" bson:"name"`
	Avatar    string        `json:"avatar" bson:"avatar"`
}

type UserModel struct {
	Id       bson.ObjectId `json:"_id" bson:"_id"`
	Name     string        `json:"name"`
	Phone    string        `json:"phone"`
	Profiles []Profile     `json:"profiles"`
}

func GetUserById(Id string) UserModel {
	s := session.Clone()
	defer s.Close()
	posts := s.DB("manch").C("users")
	user := UserModel{}
	posts.Find(bson.M{"_id": bson.ObjectIdHex(Id)}).One(&user)
	fmt.Printf("Mongo Query return for User %+v\n", user)
	return user
}

func GetProfileById(Id bson.ObjectId) Profile {
	s := session.Clone()
	defer s.Close()
	posts := s.DB("manch").C("users")
	user := UserModel{}
	posts.Find(bson.M{"profiles._id": Id}).Select(bson.M{"email": 1, "profiles.$": 1}).One(&user)
	fmt.Printf("Mongo Query return for Profile %+v\n", user.Profiles)
	return user.Profiles[0]
}
