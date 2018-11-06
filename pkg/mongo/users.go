package mongo

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
)

type UserModel struct {
	Id    bson.ObjectId `json:"_id" bson:"_id"`
	Name  string        `json:"name"`
	Phone string        `json:"phone"`
}

func GetUserById(Id string) UserModel {
	session := Session()
	posts := session.DB("manch").C("users")
	user := UserModel{}
	posts.Find(bson.M{"_id": bson.ObjectIdHex(Id)}).One(&user)
	fmt.Printf("Mongo Query return for User %+v\n", user)
	return user
}
