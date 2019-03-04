package mongo

import (
	"fmt"
	"notification-service/pkg/constants"
	"notification-service/pkg/i18n"
	"time"

	"github.com/globalsign/mgo/bson"
)

var (
	WHATSAPP_SCHEDULE_MODEL = constants.ModelNames["WHATSAPP_SCHEDULEDS"]
)

type WhatsappSchedule struct {
	Id           bson.ObjectId `json:"_id" bson:"_id"`
	UserId       bson.ObjectId `json:"user_id" bson:"user_id"`
	MobileNumber string        `json:"mobile_number" bson:"mobile_number"`
	Profile      Profile       `json:"profile" bson:"profile"`
	Message      string        `json:"message" bson:"message"`
	Schedule     Schedule      `json:"schedule" bson:"schedule"`
	CreatedAt    time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt" bson:"updatedAt"`
}

func CreateWhatsAppSchedule(user UserModel, scheduleTime time.Time) WhatsappSchedule {

	currentTime := time.Now()
	profile := user.Profiles[0]
	//creator
	c := Creator{
		Id:        bson.NewObjectId(),
		ProfileId: profile.Id,
		Name:      profile.Name,
		Avatar:    profile.Avatar,
		UserType:  user.UserType,
	}

	s := Schedule{
		Id:           bson.NewObjectId(),
		Scheduletime: scheduleTime,
		Created:      c,
	}

	message := i18n.Strings[profile.Language]["welcome_message"]
	return WhatsappSchedule{
		Id:           bson.NewObjectId(),
		UserId:       user.Id,
		MobileNumber: user.Phone,
		Profile:      profile,
		Message:      message,
		Schedule:     s,
		CreatedAt:    currentTime,
		UpdatedAt:    currentTime,
	}
}

func AddWhatsAppSchedule(document WhatsappSchedule) {
	s := session.Clone()
	defer s.Close()
	W := s.DB("manch").C(WHATSAPP_SCHEDULE_MODEL)
	err := W.Insert(document)
	if err == nil {
		fmt.Printf("inserted whatsapp schedule: %+v\n", document)
	} else {
		fmt.Println(err)
		fmt.Println("unable to add whatsapp schedule:", document.Id.Hex())
	}
}
