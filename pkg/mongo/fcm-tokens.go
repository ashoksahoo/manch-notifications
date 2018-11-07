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
	session := Session()
	T := session.DB("manch").C("fcm_tokens")
	var tokens []TokenModel
	count, _ := T.Find(q).Count()
	T.Find(q).All(&tokens)
	fmt.Printf("Count %+v\n", count)
	fmt.Printf("QUERY %+v\n", q)
	fmt.Printf("Tokens %+v\n", tokens)
	if tokens != nil {
		return tokens
	}
	return nil
}

func GetTokensByProfiles(profiles []string) []TokenModel {
	profileIds := make([]bson.ObjectId, len(profiles)-1)
	for _, profile := range profiles {
		profileIds = append(profileIds, bson.ObjectIdHex(profile))
	}
	query := bson.M{
		"profile_id": profileIds[0],
		"deleted":    false,
	}
	fmt.Printf("ProfileIds %s\n", profileIds)
	return GetTokensByQuery(&query)
}

func DeleteToken(token string) {
	session := Session()
	T := session.DB("manch").C("fcm_tokens")
	T.Update(bson.M{"fcm_token": token}, bson.M{"$set": bson.M{"deleted": true}})

}
