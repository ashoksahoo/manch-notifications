package mongo

import (
	"fmt"
	"notification-service/pkg/i18n"
	"strconv"
	"time"

	"github.com/globalsign/mgo/bson"
)

// type NotificationModel struct {
// 	Id           bson.ObjectId   `json:"_id,omitempty" bson:"_id,omitempty"`
// 	Type         string          `json:"type" bson:"type"`
// 	ResourceType string          `json:"resource_type" bson:"resource_type"`
// 	Resource     bson.ObjectId   `json:"resource_id" bson:"resource_id"`
// 	UniqueUsers  []bson.ObjectId `json:"profile_ids" bson:"profile_ids"`
// 	Identifier   string          `json:"identifier" bson:"identifier"`
// }

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
	FailReason string    `json:"failed_reason" bson:"failed_reason"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
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

func RemoveNotificationUser(rId bson.ObjectId, t string, u bson.ObjectId) {
	s := session.Clone()
	defer s.Close()
	N := s.DB("manch").C("notifications")
	query, update := bson.M{"resource_id": rId, "type": t}, bson.M{
		"$pull": bson.M{"profile_ids": u},
	}
	N.Update(query, update)
	fmt.Printf("removed from notifications model profile id %s\n", u.Hex())
}

func CreateNotification(notification NotificationModel) NotificationModel {
	s := session.Clone()
	defer s.Close()
	N := s.DB("manch").C("notifications")
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
		NUUID:           notification.NUUID,
	}

	count, _ := N.Find(bson.M{"identifier": notification.Identifier, "is_read": notification.IsRead}).Count()
	print(count)
	if count > 0 {
		N.Upsert(bson.M{"identifier": notification.Identifier, "is_read": notification.IsRead}, bson.M{
			"$addToSet": bson.M{"participants": notification.Participants[0]},
		})
	} else {
		N.Insert(n)
	}
	return GetNotificationByIdentifier(notification.Identifier)
}

func UpdateNotificationMessage(id bson.ObjectId, message string) {
	s := session.Clone()
	defer s.Close()
	N := s.DB("manch").C("notifications")
	N.Update(bson.M{"_id": id}, bson.M{"$set": bson.M{"message": message}})
}

func UpdateNotification(query, update bson.M) {
	s := session.Clone()
	defer s.Close()
	N := s.DB("manch").C("notifications")
	N.Update(query, bson.M{"$set": update})
}

// func CreateNotification(rId bson.ObjectId, t string, rT string, u bson.ObjectId) NotificationModel {
// 	s := session.Clone()
// 	defer s.Close()
// 	N := s.DB("manch").C("notifications")
// 	n := NotificationModel{
// 		Resource:     rId,
// 		Type:         t,
// 		ResourceType: rT,
// 		UniqueUsers:  []bson.ObjectId{u},
// 		Identifier:   GenerateIdentifier(rId, t),
// 	}
// 	count, _ := N.Find(bson.M{"resource_id": rId, "type": t}).Count()
// 	print(count)
// 	if count > 0 {
// 		N.Upsert(bson.M{"resource_id": rId, "type": t}, bson.M{
// 			"$addToSet": bson.M{"profile_ids": u},
// 		})
// 	} else {
// 		N.Insert(n)
// 	}
// 	return GetNotificationByResource(rId, t)

// }

func GetNotificationByIdentifier(identifier string) NotificationModel {
	s := session.Clone()
	defer s.Close()
	N := s.DB("manch").C("notifications")
	notif := NotificationModel{}
	N.Find(bson.M{"identifier": identifier, "is_read": false}).One(&notif)
	return notif
}

// func GetNotificationByResource(rId bson.ObjectId, t string) NotificationModel {
// 	s := session.Clone()
// 	defer s.Close()
// 	N := s.DB("manch").C("notifications")
// 	notif := NotificationModel{}
// 	N.Find(bson.M{"resource_id": rId, "type": t}).One(&notif)
// 	return notif
// }
