package mongo

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
)

type Profile struct {
	Id     bson.ObjectId `json:"_id" bson:"_id"`
	Avatar string        `json:"avatar"`
	Name   string        `json:"name"`
}

type Creator struct {
	Id        bson.ObjectId `json:"_id" bson:"_id"`
	ProfileId bson.ObjectId `json:"_id" bson:"_id"`
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
	session := Session()
	posts := session.DB("manch").C("users")
	user := UserModel{}
	posts.Find(bson.M{"_id": bson.ObjectIdHex(Id)}).One(&user)
	fmt.Printf("Mongo Query return for User %+v\n", user)
	return user
}
