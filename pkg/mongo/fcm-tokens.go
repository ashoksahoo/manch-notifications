package mongo

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
)

type TokenModel struct {
	Id        bson.ObjectId `json:"_id" bson:"_id"`
	ProfileId bson.ObjectId `json:"profile_id" bson:"profile_id"`
	Token     string        `json:"fcm_token" bson:"fcm_token"`
	Creator
}

func GetTokensByQuery(q *bson.M) []TokenModel {
	s := session.Clone()
	defer s.Close()
	T := s.DB("manch").C("fcm_tokens")
	var tokens []TokenModel
	T.Find(q).All(&tokens)
	if tokens != nil {
		return tokens
	}
	return nil
}

func GetTokensByProfiles(profiles []bson.ObjectId) []TokenModel {
	query := bson.M{
		"profile_id": profiles[0],
		"deleted":    false,
	}
	return GetTokensByQuery(&query)
}

func DeleteToken(token string) {
	s := session.Clone()
	defer s.Close()
	T := s.DB("manch").C("fcm_tokens")
	T.Update(bson.M{"fcm_token": token}, bson.M{"$set": bson.M{"deleted": true}})
	fmt.Printf("deleted token")

}
