package mongo

import (
	// "fmt"
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
	UserType  string        `json:"type" bson:"type"`
}

type UserModel struct {
	Id       bson.ObjectId `json:"_id" bson:"_id"`
	Name     string        `json:"name" bson:"name"`
	Phone    string        `json:"phone" bson:"phone"`
	Profiles []Profile     `json:"profiles" bson:"profiles"`
	UserType string		   `json:"type" bson:"type"`
}

func GetUserById(Id string) UserModel {
	s := session.Clone()
	defer s.Close()
	users := s.DB("manch").C("users")
	user := UserModel{}
	users.Find(bson.M{"_id": bson.ObjectIdHex(Id)}).One(&user)
	//fmt.Printf("Mongo Query return for User %+v\n", user)
	return user
}

func GetBotUsers() []UserModel {
	s := session.Clone()
	defer s.Close()
	users := s.DB("manch").C("users")
	allUsers := []UserModel{}
	users.Find(bson.M{"type": "bot"}).All(&allUsers)
	// fmt.Println("allusers: ", allUsers)
	return allUsers
}

func GetProfileById(Id bson.ObjectId) Profile {
	s := session.Clone()
	defer s.Close()
	posts := s.DB("manch").C("users")
	user := UserModel{}
	posts.Find(bson.M{"profiles._id": Id}).Select(bson.M{"email": 1, "profiles.$": 1}).One(&user)
	//fmt.Printf("Mongo Query return for Profile %+v\n", user.Profiles)
	return user.Profiles[0]
}
