package mongo

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
)

type TokenModel struct {
	Id        bson.ObjectId `json:"_id" bson:"_id"`
	ProfileId bson.ObjectId `json:"profile_id" bson:"profile_id"`
	Token     string        `json:"fcm_token" bson:"fcm_token"`
	Created   Creator       `json:"created"`
}

func GetTokensByQuery(q bson.M) []TokenModel {
	session := Session()
	posts := session.DB("manch").C("fcm_tokens")
	var tokens []TokenModel
	posts.Find(q).All(&tokens)
	if tokens != nil {
		return tokens
	}
	return nil
}

func GetTokensByProfiles(profiles []string) []TokenModel {
	profileIds := make([]bson.ObjectId, 100)
	for _, profile := range profiles {
		profileIds = append(profileIds, bson.ObjectIdHex(profile))
	}
	query := bson.M{
		"profile_id": bson.M{"$in": profiles},
		"deleted":    false,
	}
	fmt.Printf("Query %s\n", query)
	return GetTokensByQuery(query)
}
