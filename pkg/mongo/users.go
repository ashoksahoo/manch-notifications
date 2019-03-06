package mongo

import (
	"fmt"
	"notification-service/pkg/constants"
	"notification-service/pkg/utils"
	"os"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/mitchellh/mapstructure"
)

type Stats struct {
	PosLevel1 int `json:"pos_level_1_posts" bson:"pos_level_1_posts"`
	PosLevel2 int `json:"pos_level_2_posts" bson:"pos_level_2_posts"`
	NegLevel1 int `json:"neg_level_1_posts" bson:"neg_level_1_posts"`
	NegLevel2 int `json:"neg_level_2_posts" bson:"neg_level_2_posts"`
}

type BlackList struct {
	Status       string    `json:"status" bson:"status"`
	WarnCount    string    `json:"warn_count" bson:"warn_count"`
	LastWarnedOn time.Time `json:"last_warned_on" bson:"last_warned_on"`
	BlockedOn    time.Time `json:"blocked_on" bson:"blocked_on"`
	BlockedTill  time.Time `json:"blocked_till" bson:"blocked_till"`
	Reason       string    `json:"reason" bson:"reason"`
}

var (
	USERS_MODEL = constants.ModelNames["USERS"]
)

type Profile struct {
	Id            bson.ObjectId `json:"_id" bson:"_id"`
	Avatar        string        `json:"avatar"`
	Name          string        `json:"name"`
	Language      string        `json:"language"`
	FollowerCount int           `json:"no_of_followers" bson:"no_of_followers"`
	Type          string        `json:"type" bson:"type"`
	RandomName    bool          `json:"random_name" bson:"random_name"`
}

type Creator struct {
	Id        bson.ObjectId `json:"_id" bson:"_id"`
	ProfileId bson.ObjectId `json:"profile_id" bson:"profile_id"`
	Name      string        `json:"name" bson:"name"`
	Avatar    string        `json:"avatar" bson:"avatar"`
	UserType  string        `json:"type" bson:"type"`
}

type UserModel struct {
	Id        bson.ObjectId `json:"_id" bson:"_id"`
	Name      string        `json:"name" bson:"name"`
	Phone     string        `json:"phone" bson:"phone"`
	Profiles  []Profile     `json:"profiles" bson:"profiles"`
	UserType  string        `json:"type" bson:"type"`
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
	Language  string        `json:"language" bson:"language"`
}

func GetUserById(Id string) UserModel {
	s := session.Clone()
	defer s.Close()
	users := s.DB("manch").C(USERS_MODEL)
	user := UserModel{}
	users.Find(bson.M{"_id": bson.ObjectIdHex(Id)}).One(&user)
	//fmt.Printf("Mongo Query return for User %+v\n", user)
	return user
}

func GetUserByProfileId(ProfileId string) UserModel {
	s := session.Clone()
	defer s.Close()
	users := s.DB("manch").C(USERS_MODEL)
	user := UserModel{}
	users.Find(bson.M{"profiles._id": bson.ObjectIdHex(ProfileId)}).One(&user)
	return user
}

func GetBotUsers(language string) []UserModel {
	s := session.Clone()
	defer s.Close()
	users := s.DB("manch").C(USERS_MODEL)
	allUsers := []UserModel{}
	users.Find(bson.M{"type": "bot", "language": language}).All(&allUsers)
	// fmt.Println("allusers: ", allUsers)
	fmt.Println("total bot users", len(allUsers))
	return allUsers
}

// func GetProfile(profileId string) Profile {
// 	s := session.Clone()
// 	defer s.Close()
// 	users := s.DB("manch").C("users")
// }

func GetProfileById(Id bson.ObjectId) Profile {
	s := session.Clone()
	defer s.Close()
	users := s.DB("manch").C(USERS_MODEL)
	user := UserModel{}
	users.Find(bson.M{"profiles._id": Id}).Select(bson.M{"email": 1, "profiles.$": 1}).One(&user)
	fmt.Printf("Mongo Query return for Profile %+v\n", user.Profiles)
	return user.Profiles[0]
}

func GetProfilesByIds(Ids []string) []Profile {

	profileIds := []bson.ObjectId{}

	for _, id := range Ids {
		profileIds = append(profileIds, bson.ObjectIdHex(id))
	}

	s := session.Clone()
	defer s.Close()
	usersCollection := s.DB("manch").C(USERS_MODEL)
	pipe := usersCollection.Pipe([]bson.M{
		{"$match": bson.M{"profiles._id": bson.M{"$in": profileIds}}},
		{"$unwind": "$profiles"},
		{"$match": bson.M{"profiles._id": bson.M{"$in": profileIds}}},
		{"$project": bson.M{"profiles": 1}},
	})
	resp := []bson.M{}
	err := pipe.All(&resp)
	if err != nil {
		fmt.Println("Error found while fetching profiles", err)
		return []Profile{}
	}
	profiles := []Profile{}
	for _, res := range resp {
		profile := Profile{}
		id := res["profiles"].(bson.M)["_id"]
		_id, _ := id.(bson.ObjectId)
		profile.Id = _id
		mapstructure.Decode(res["profiles"], &profile)
		profiles = append(profiles, profile)
	}

	return profiles

}

func GetBotProfilesIds(language string) (int, []string) {
	env := os.Getenv("env")
	if env != "production" {
		botUsers := GetBotUsers(language)
		botProfilesIds := make([]string, 0, 1000)
		i := 0
		for _, botUser := range botUsers {
			profiles := botUser.Profiles
			for _, profile := range profiles {
				if profile.Id.Hex() == constants.MANCH_OFFICIAL_PROFILE_HE || profile.Id.Hex() == constants.MANCH_OFFICIAL_PROFILE_TE {
					continue
				}
				botProfilesIds = append(botProfilesIds, profile.Id.Hex())
				i++
			}
		}
		return i, botProfilesIds
	}
	botProfilesIds := utils.BotProfiles[language]
	return len(botProfilesIds), botProfilesIds
}

func UpdateProfileById(profileId bson.ObjectId, update bson.M) {
	s := session.Clone()
	defer s.Close()

	C := s.DB("manch").C(USERS_MODEL)

	fmt.Println("update:", update)
	fmt.Println("profile id is ", bson.M{"profiles._id": profileId})
	err := C.Update(bson.M{"profiles._id": profileId}, update)

	if err != nil {
		fmt.Println("Error while updating profile", err)
	} else {
		fmt.Println("update profiles successfully")
	}

}

func UpdateUser(query, update bson.M) {
	s := session.Clone()
	defer s.Close()
	users := s.DB("manch").C("users")
	err := users.Update(query, update)
	if err != nil {
		fmt.Println(err, query, update)
	}
}
