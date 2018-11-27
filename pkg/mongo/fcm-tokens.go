package mongo

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
)

type TokenModel struct {
	Id        bson.ObjectId `json:"_id" bson:"_id"`
	ProfileId bson.ObjectId `json:"profile_id" bson:"profile_id"`
	Token     string        `json:"fcm_token" bson:"fcm_token"`
	Created   Creator       `json:"created" bson:"created"`
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

func botFilter(ta []TokenModel) (ret []TokenModel) {
	for _, t := range ta {
		if t.Created.UserType == "bot" {
			ret = append(ret, t)
		}
	}
	return
}

func GetTokensByProfiles(profiles []bson.ObjectId) (tokens []TokenModel) {
	query := bson.M{
		"profile_id": profiles[0],
		"deleted":    false,
	}
	tokens = GetTokensByQuery(&query)
	botTokens := botFilter(tokens)
	return botTokens
}

func DeleteToken(token string) {
	s := session.Clone()
	defer s.Close()
	T := s.DB("manch").C("fcm_tokens")
	T.Update(bson.M{"fcm_token": token}, bson.M{"$set": bson.M{"deleted": true}})
	fmt.Printf("deleted token")
}
