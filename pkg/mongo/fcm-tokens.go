package mongo

import (
	"github.com/globalsign/mgo/bson"
)

type TokenModel struct {
	Id      bson.ObjectId `json:"_id" bson:"_id"`
	Title   string        `json:"title"`
	Profile bson.ObjectId `json:"profile_id" bson:"profile_id"`
}

func GetTokensByQuery(Query bson.M) []TokenModel {
	session := Session()
	posts := session.DB("manch").C("fcm_tokens")
	var tokens []TokenModel
	posts.Find(Query).All(&tokens)
	return tokens
}

func GetTokensByProfiles(profiles []string) []TokenModel {
	return GetTokensByQuery(bson.M{"profile_id": bson.M{"$in": profiles}, "deleted": false})
}
