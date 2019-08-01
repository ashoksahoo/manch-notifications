package mongo

import (
	"notification-service/pkg/constants"
	"time"

	"github.com/globalsign/mgo/bson"
)

var (
	USER_ACTIVE_MODEL = constants.ModelNames["USER_ACTIVE_HOUR"]
)

type UserActiveHour struct {
	Id             string    `json:"_id"`
	ProfileId      string    `json:"profile_id"`
	LastActiveHour time.Time `json:"last_active_hour"`
}

func CountUserActiveHour(query bson.M) int {
	s := session.Clone()
	defer s.Close()
	U := s.DB("manch").C(USER_ACTIVE_MODEL)
	n, _ := U.Find(query).Count()
	return n
}
