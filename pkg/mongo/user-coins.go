package mongo

import (
	"fmt"
	"notification-service/pkg/constants"
	"notification-service/pkg/utils"
	"time"

	"github.com/globalsign/mgo/bson"
)

var (
	USER_COINS_MODEL        = constants.ModelNames["USER_COINS"]
	USER_LEADERBOARDS_MODEL = constants.ModelNames["USER_LEADERBOARDS"]
)

type UserCoinsModel struct {
	Id          bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	ProfileId   bson.ObjectId `json:"profile_id" bson:"profile_id"`
	Action      string        `json:"action" bson:"action"`
	CoinsEarned int           `json:"coins_earned" bson:"coins_earned"`
	DayKey      string        `json:"day_key" bson:"day_key"`
	WeekKey     string        `json:"week_key" bson:"week_key"`
	MonthKey    string        `json:"month_key" bson:"month_key"`
	YearKey     string        `json:"year_key" bson:"year_key"`
	CreatedAt   time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt" bson:"updatedAt"`
}

type UserLeaderBoard struct {
	Id          bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	ProfileId   bson.ObjectId `json:"profile_id" bson:"profile_id"`
	Key         string        `json:"key" bson:"key"`
	Granularity int           `json:"granularity" bson:"granularity"`
	Coins       int           `json:"coins" bson:"coins"`
	CreatedAt   time.Time     `json:"createdAt" bson:"updatedAt"`
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
		return
	}

	// update User profiles coin
	UpdateProfileById(userCoins.ProfileId, bson.M{
		"$inc": bson.M{"profiles.$.total_coins": userCoins.CoinsEarned},
	})

	key1 := userCoins.ProfileId.Hex() + "_" + userCoins.DayKey
	key2 := userCoins.ProfileId.Hex() + "_" + userCoins.WeekKey
	key3 := userCoins.ProfileId.Hex() + "_" + userCoins.MonthKey
	key4 := userCoins.ProfileId.Hex() + "_" + userCoins.YearKey

	CreateUserLeaderBoard(key1, userCoins.ProfileId, userCoins.DayKey, userCoins.CoinsEarned)
	CreateUserLeaderBoard(key2, userCoins.ProfileId, userCoins.WeekKey, userCoins.CoinsEarned)
	CreateUserLeaderBoard(key3, userCoins.ProfileId, userCoins.MonthKey, userCoins.CoinsEarned)
	CreateUserLeaderBoard(key4, userCoins.ProfileId, userCoins.YearKey, userCoins.CoinsEarned)
}

func CreateUserLeaderBoard(key string, profileId bson.ObjectId, granularity string, coins int) {
	s := session.Clone()
	defer s.Close()
	C := s.DB("manch").C(USER_LEADERBOARDS_MODEL)

	info, err := C.Upsert(bson.M{"key": key}, bson.M{
		"$inc": bson.M{"coins": coins},
		"$set": bson.M{"updatedAt": time.Now()},
		"$setOnInsert": bson.M{
			"key":         key,
			"profile_id":  profileId,
			"granularity": granularity,
			"createdAt":   time.Now(),
		},
	})
	if err != nil {
		fmt.Println("Error while upsert CreateUserLeaderBoard", err)
	} else {
		fmt.Println("Upsert Info CreateUserLeaderBoard", info)
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
