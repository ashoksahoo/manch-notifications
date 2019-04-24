package mongo

import (
	"fmt"
	"notification-service/pkg/constants"
	"time"

	"github.com/globalsign/mgo/bson"
)

var (
	WHATSAPP_SCHEDULE_MODEL = constants.ModelNames["WHATSAPP_SCHEDULEDS"]
)

type WhatsAppMessageBody struct {
	MessagePlatform string `json:"messagePlatform" bson:"messagePlatform"`
	MessageType     string `json:"messsageType" bson:"messsageType"`
	Number          string `json:"number" bson:"number"`
	Message         string `json:"message bson:"message"`
}

type WhatsappSchedule struct {
	Id          bson.ObjectId       `json:"_id" bson:"_id"`
	ProfileId   bson.ObjectId       `json:"profile_id" bson:"profile_id"`
	MessageBody WhatsAppMessageBody `json:"message_body" bson:"message_body"`
	Schedule    Schedule            `json:"schedule" bson:"schedule"`
	CreatedAt   time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt" bson:"updatedAt"`
}

func CreateWhatsAppSchedule(user UserModel, scheduleTime time.Time, message string) WhatsappSchedule {

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

	messageBody := WhatsAppMessageBody{
		MessagePlatform: "WHATSAPP",
		MessageType:     "TEXT",
		Number:          "91" + user.Phone,
		Message:         message,
	}
	return WhatsappSchedule{
		Id:          bson.NewObjectId(),
		ProfileId:   profile.Id,
		MessageBody: messageBody,
		Schedule:    s,
		CreatedAt:   currentTime,
		UpdatedAt:   currentTime,
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
