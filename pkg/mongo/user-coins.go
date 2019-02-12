package mongo

import (
	"fmt"
	"notification-service/pkg/constants"
	"notification-service/pkg/utils"
	"time"

	"github.com/globalsign/mgo/bson"
)

var (
	USER_COINS_MODEL = constants.ModelNames["USER_COINS_MODEL"]
)

type UserCoinsModel struct {
	Id          bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	ProfileId   bson.ObjectId `json:"profile_id" bson:"profile_id"`
	Action      string        `json:"action" bson:"action"`
	CoinsEarned int           `json:"coins_earned" bson:"coins_earned"`
	DayKey      int           `json:"day_key" bson:"day_key"`
	WeekKey     int           `json:"week_key" bson:"week_key"`
	MonthKey    int           `json:"month_key" bson:"month_key"`
	YearKey     int           `json:"year_key" bson:"year_key"`
	CreatedAt   time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt" bson:"updatedAt"`
}

func CreateUserCoin(userCoins UserCoinsModel) {
	dayKey, weekKey, monthKey, yearKey := utils.GetCurrentDateKeys()
	userCoins.DayKey = dayKey
	userCoins.WeekKey = weekKey
	userCoins.MonthKey = monthKey
	userCoins.YearKey = yearKey
	userCoins.CreatedAt = time.Now()
	userCoins.UpdatedAt = time.Now()

	s := session.Clone()
	defer s.Close()
	C := s.DB("manch").C(USER_COINS_MODEL)
	err := C.Insert(userCoins)
	if err != nil {
		fmt.Println("Error while inserting userCoins")
	}
}

func GetUserCoinsByQuery(query bson.M) (error, UserCoinsModel) {
	s := session.Clone()
	defer s.Close()
	C := s.DB("manch").C(USER_COINS_MODEL)
	data := UserCoinsModel{}
	err := C.Find(query).One(&data)
	return err, data
}
