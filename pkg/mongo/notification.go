package mongo

import (
	"fmt"
	"notification-service/pkg/constants"
	"notification-service/pkg/i18n"
	"strconv"
	"time"

	"github.com/gofrs/uuid"

	"github.com/globalsign/mgo/bson"
)

var (
	NOTIFICATION_V2_MODEL = constants.ModelNames["NOTIFICATION_V2"]
)

type Entity struct {
	EntityId   bson.ObjectId `json:"entity_id" bson:"entity_id"`
	EntityType string        `json:"entity_type" bson:"entity_type"`
}

type MessageMeta struct {
	Template string         `json:"template_name" bson:"template_name"`
	Data     i18n.DataModel `json:"data" bson:"data"`
}

type PushMeta struct {
	Status     string    `json:"status" bson:"status"`
	PushId     string    `json:"push_id" bson:"push_id"`
	FailReason string    `json:"error" bson:"error"`
	CreatedAt  time.Time `json:"date" bson:"date"`
}

type NotificationModel struct {
	Id              bson.ObjectId   `json:"_id,omitempty" bson:"_id,omitempty"`
	Receiver        bson.ObjectId   `json:"receiver" bson:"receiver"`
	Identifier      string          `json:"identifier" bson:"identifier"`
	Message         string          `json:"message" bson:"message"`
	IsRead          bool            `json:"is_read" bson:"is_read"`
	Participants    []bson.ObjectId `json:"participants" bson:"participants"`
	DisplayTemplate string          `json:"display_template" bson:"display_template"`
	EntityGroupId   string          `json:"entity_group_id" bson:"entity_group_id"`
	ActionType      string          `json:"action_type" bson:"action_type"`
	ActionId        bson.ObjectId   `json:"action_id" bson:"action_id"`
	NId             string          `json:"n_id" bson:"n_id"`
	NUUID           string          `json:"nuuid" bson:"nuuid"`
	Purpose         string          `json:"purpose" bson:"purpose"`
	Entities        []Entity        `json:"entities" bson:"entities"`
	MessageMeta     MessageMeta     `json:"message_meta" bson:"message_meta"`
	Push            PushMeta        `json:"push" bson:"push"`
	CreatedAt       time.Time       `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt" bson:"updatedAt"`
	Delivered       bool            `json:"delivered" bson:"delivered"`
}

func GenerateIdentifier(Id bson.ObjectId, t string) string {
	ts := Id.Time().Unix()
	ts2018 := ts - 1514764800
	var identifier string
	switch t {
	case "like":
		identifier = "1"
	case "comment":
		identifier = "2"
	}
	return strconv.FormatInt(ts2018, 10) + identifier
}

func RemoveParticipants(identifier string, isRead bool, participant bson.ObjectId) {
	s := session.Clone()
	defer s.Close()
	N := s.DB("manch").C(NOTIFICATION_V2_MODEL)
	query, update := bson.M{"identifier": identifier, "is_read": isRead},
		bson.M{"$pull": bson.M{"participants": participant}}
	N.Update(query, update)
	fmt.Printf("removed from participants with id %s\n", participant.Hex())
}

func CreateNotification(notification NotificationModel) NotificationModel {
	s := session.Clone()
	defer s.Close()
	N := s.DB("manch").C(NOTIFICATION_V2_MODEL)

	push := PushMeta{
		Status:    constants.NotificationStatus["PENDING"],
		CreatedAt: time.Now(),
	}

	n := NotificationModel{
		Receiver:        notification.Receiver,
		Identifier:      notification.Identifier,
		Message:         notification.Message,
		IsRead:          false,
		Participants:    notification.Participants,
		DisplayTemplate: notification.DisplayTemplate,
		EntityGroupId:   notification.EntityGroupId,
		ActionId:        notification.ActionId,
		ActionType:      notification.ActionType,
		NId:             GenerateIdentifier(notification.ActionId, notification.ActionType),
		Purpose:         notification.Purpose,
		Entities:        notification.Entities,
		Push:            push,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Delivered:       false,
	}

	count, _ := N.Find(bson.M{"identifier": notification.Identifier, "is_read": notification.IsRead}).Count()
	print(count)
	if count > 0 {
		var nuuid string
		value, error := uuid.NewV4()
		if error != nil {
			nuuid = ""
		} else {
			nuuid = value.String()
		}
		N.Upsert(bson.M{"identifier": notification.Identifier, "is_read": notification.IsRead}, bson.M{
			"$set":         bson.M{"updatedAt": time.Now()},
			"$addToSet":    bson.M{"participants": notification.Participants[0]},
			"$setOnInsert": bson.M{"nuuid": nuuid},
		})
	} else {
		if n.NUUID == "" {
			var nuuid string
			value, error := uuid.NewV4()
			if error != nil {
				nuuid = ""
			} else {
				nuuid = value.String()
			}
			n.NUUID = nuuid
		}
		N.Insert(n)
	}
	_, notif := GetNotificationByIdentifier(notification.Identifier)
	return notif
}

func UpdateNotification(query, update bson.M) {
	update["updatedAt"] = time.Now()
	fmt.Println("update is", update)
	s := session.Clone()
	defer s.Close()
	N := s.DB("manch").C(NOTIFICATION_V2_MODEL)
	N.Update(query, bson.M{"$set": update})
}

func GetNotificationByIdentifier(identifier string) (error, NotificationModel) {
	s := session.Clone()
	defer s.Close()
	N := s.DB("manch").C(NOTIFICATION_V2_MODEL)
	notif := NotificationModel{}
	err := N.Find(bson.M{"identifier": identifier, "is_read": false}).One(&notif)
	return err, notif
}

// find, select, sort, skip, limit
func GetNotificationByQuery(query bson.M) (error, []NotificationModel)  {
	s := session.Clone()
	defer s.Close()
	fmt.Println("Query is ", query)
	N := s.DB("manch").C(NOTIFICATION_V2_MODEL)
	notifications := []NotificationModel{}
	var err error
	if limit, ok := query["limit"]; ok {
		fmt.Println("ok")
		delete(query, "limit")
		skip := query["skip"]
		delete(query, "skip")
		err = N.Find(query).Skip(skip.(int)).Limit(limit.(int)).All(&notifications)	
	} else {
		err = N.Find(query).All(&notifications)
	}
	return err, notifications
}

func GetNotificationById(id bson.ObjectId) (error, NotificationModel) {
	s := session.Clone()
	defer s.Close()
	N := s.DB("manch").C(NOTIFICATION_V2_MODEL)
	notif := NotificationModel{}
	err := N.Find(bson.M{"_id": id}).One(&notif)
	return err, notif
}
