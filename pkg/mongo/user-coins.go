package mongo

import (
	"fmt"
	"notification-service/pkg/constants"
	"notification-service/pkg/utils"
	"time"

	"github.com/globalsign/mgo/bson"
)

var (
	USER_COINS_MODEL          = constants.ModelNames["USER_COINS"]
	USER_LEADERBOARDS_MODEL   = constants.ModelNames["USER_LEADERBOARDS"]
	USER_COINS_SCHEDULE_MODEL = constants.ModelNames["USER_COINS_SCHEDULE_MODEL"]
)

var ManchProfiles = []string{
	"5bc498ad9ca925186ffb64b0",
	"5b5449a1e24cb4255b00ffb8",
	"5c3c3bfd89ac4a794d45b14d",
	"5c1c92c8eda9bd1771bcf0a7",
	"5c2f90a86aeb6b6c7345dc33",
	"5c14bb2f5837a165cdead3a1",
}

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
	Id               bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	ProfileId        bson.ObjectId `json:"profile_id" bson:"profile_id"`
	Key              string        `json:"key" bson:"key"`
	Granularity      int           `json:"granularity" bson:"granularity"`
	GranularityStart time.Time     `json:"granularity_start" bson:"granularity_start"`
	GranularityEnd   time.Time     `json:"granularity_end" bson:"granularity_end"`
	Coins            int           `json:"coins" bson:"coins"`
	CreatedAt        time.Time     `json:"createdAt" bson:"updatedAt"`
	UpdatedAt        time.Time     `json:"updatedAt" bson:"updatedAt"`
}

type UserCoinsModelScheduleModel struct {
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
	Schedule    Schedule      `json:"schedule" bson:"schedule"`
	Created     Creator       `json:"created" bson:"created"`
}

func CreateUserCoin(userCoins UserCoinsModel) {

	if utils.ContainsStr(ManchProfiles, userCoins.ProfileId.Hex()) {
		return
	}

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

	loc, _ := time.LoadLocation("Asia/Kolkata")
	granularityDayStart := utils.GetStartOfDay(time.Now().In(loc))
	granularityDayEnd := utils.GetEndOfDay(time.Now().In(loc))
	granularityWeekStart, granularityWeekEnd := utils.WeekRange()
	granularityMonthStart := utils.GetStartOfMonth()
	granularityMonthEnd := utils.GetEndOfMonth()
	granularityYearStart := utils.GetStartOfYear()
	granularityYearEnd := utils.GetEndOfYear()

	CreateUserLeaderBoard(key1, userCoins.ProfileId, userCoins.DayKey, userCoins.CoinsEarned, granularityDayStart, granularityDayEnd)
	CreateUserLeaderBoard(key2, userCoins.ProfileId, userCoins.WeekKey, userCoins.CoinsEarned, granularityWeekStart, granularityWeekEnd)
	CreateUserLeaderBoard(key3, userCoins.ProfileId, userCoins.MonthKey, userCoins.CoinsEarned, granularityMonthStart, granularityMonthEnd)
	CreateUserLeaderBoard(key4, userCoins.ProfileId, userCoins.YearKey, userCoins.CoinsEarned, granularityYearStart, granularityYearEnd)

}

func CreateUserCoinSchedule(userCoinsSchedule UserCoinsModelScheduleModel, scheduleTime time.Time) {

	if utils.ContainsStr(ManchProfiles, userCoinsSchedule.ProfileId.Hex()) {
		return
	}

	profile := GetProfileById(userCoinsSchedule.ProfileId)
	c := Creator{
		Id:        bson.NewObjectId(),
		ProfileId: profile.Id,
		Name:      profile.Name,
		Avatar:    profile.Avatar,
		UserType:  profile.Type,
	}

	schedule := Schedule{
		Id:           bson.NewObjectId(),
		Scheduletime: scheduleTime,
		Created:      c,
	}

	dayKey, weekKey, monthKey, yearKey := utils.GetCurrentDateKeys()
	userCoinsSchedule.DayKey = dayKey
	userCoinsSchedule.WeekKey = weekKey
	userCoinsSchedule.MonthKey = monthKey
	userCoinsSchedule.YearKey = yearKey
	userCoinsSchedule.CreatedAt = time.Now()
	userCoinsSchedule.UpdatedAt = time.Now()

	userCoinsSchedule.Schedule = schedule
	userCoinsSchedule.Created = c

	s := session.Clone()
	defer s.Close()
	C := s.DB("manch").C(USER_COINS_SCHEDULE_MODEL)
	err := C.Insert(userCoinsSchedule)
	if err != nil {
		fmt.Println("Error while inserting userCoins")
	}
}

func CreateUserLeaderBoard(key string, profileId bson.ObjectId, granularity string, coins int, granularityStart, granularityEnd time.Time) {
	s := session.Clone()
	defer s.Close()
	C := s.DB("manch").C(USER_LEADERBOARDS_MODEL)

	info, err := C.Upsert(bson.M{"key": key}, bson.M{
		"$inc": bson.M{"coins": coins},
		"$set": bson.M{
			"updatedAt":         time.Now(),
			"granularity_start": granularityStart,
			"granularity_end":   granularityEnd,
		},
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
