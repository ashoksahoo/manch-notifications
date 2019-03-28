package mongo

import (
	"fmt"
	"notification-service/pkg/constants"

	"github.com/globalsign/mgo/bson"
)

var (
	REFERRALS = constants.ModelNames["REFERRALS"]
)

type ReferralModel struct {
	Referrer                   string                 `json:"referrer" bson:"referrer"`
	Referree                   string                 `json:"referree" bson:"referree"`
	ProfileId                  string                 `json:"profile_id" bson:"profile_id"`
	FirebaseIntanceId          string                 `json:"firabase_instance_id" bson:"firabase_instance_id"`
	ReferringParams            map[string]interface{} `json:"referring_params" bson:"referring_params"`
	LatestReferringParams      map[string]interface{} `json:"latest_referring_params" bson:"latest_referring_params"`
	BranchFirstReferringParams map[string]interface{} `json:"branch_first_referring_params" bson:"branch_first_referring_params"`
}

func GetReferralsByProfileId(profileId string) (error, ReferralModel) {
	s := session.Clone()
	defer s.Close()
	referralData := ReferralModel{}
	R := s.DB("manch").C(REFERRALS)
	err := R.Find(bson.M{"profile_id": profileId}).One(&referralData)
	fmt.Println("error", err)
	return err, referralData
}

func GetReferralsByQuery(query bson.M) (error, ReferralModel) {
	s := session.Clone()
	defer s.Close()
	referralData := ReferralModel{}
	R := s.DB("manch").C(REFERRALS)
	err := R.Find(query).One(&referralData)
	fmt.Println("error", err)
	return err, referralData
}
