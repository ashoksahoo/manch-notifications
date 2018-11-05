package mongo

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
)

type TokenModel struct {
	Id    string `json:"_id"`
	Title string `json:"title"`
}

func GetTokenByQuery(Query string) TokenModel {
	session := Session()
	posts := session.DB("manch").C("posts")
	token := TokenModel{}
	posts.Find(bson.M{"_id": bson.ObjectIdHex(Query)}).One(&token)
	fmt.Printf("Mongo Query return for Post %+v\n", token)
	return token
}
