package mongo

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
)

type PostModel struct {
	Id    string `json:"_id"`
	Title string `json:"title"`
}

func GetPostById(Id string) PostModel {
	session := Session()
	posts := session.DB("manch").C("posts")
	newPost := PostModel{}
	posts.Find(bson.M{"_id": bson.ObjectIdHex(Id)}).One(&newPost)
	fmt.Printf("Mongo Query return for Post %+v\n", newPost)
	return newPost
}
