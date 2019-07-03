package mongo

import (
	"fmt"
	"notification-service/pkg/constants"
	"time"

	"github.com/globalsign/mgo/bson"
)

var (
	FCM_TOKEN_MODEL = constants.ModelNames["FCM_TOKENS"]
)

type TokenModel struct {
	Id                 bson.ObjectId `json:"_id" bson:"_id"`
	ProfileId          bson.ObjectId `json:"profile_id" bson:"profile_id"`
	Token              string        `json:"fcm_token" bson:"fcm_token"`
	LastVoteNotifiedAt time.Time     `json:"last_vote_notified_at" bson:"last_vote_notified_at"`
	Created            Creator       `json:"created" bson:"created"`
}

func GetTokensByQuery(q *bson.M) []TokenModel {
	s := session.Clone()
	defer s.Close()
	T := s.DB("manch").C(FCM_TOKEN_MODEL)
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
	return tokens
}

func DeleteToken(token string) {
	s := session.Clone()
	defer s.Close()
	T := s.DB("manch").C(FCM_TOKEN_MODEL)
	T.UpdateAll(bson.M{"fcm_token": token}, bson.M{"$set": bson.M{"deleted": true}})
	fmt.Printf("deleted token")
}

func UpdateFCMTokenByQuery(query bson.M, update bson.M) {
	s := session.Clone()
	defer s.Close()
	T := s.DB("manch").C(FCM_TOKEN_MODEL)
	info, err := T.UpdateAll(query, update)
	if err != nil {
		fmt.Println("Error while updating fcm records")
	} else {
		fmt.Println("Fcm record update info", info)
	}
}