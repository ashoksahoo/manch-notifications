package mongo

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"strconv"
)

type NotificationModel struct {
	Id           bson.ObjectId   `json:"_id,omitempty" bson:"_id,omitempty"`
	Type         string          `json:"type" bson:"type"`
	ResourceType string          `json:"resource_type" bson:"resource_type"`
	Resource     bson.ObjectId   `json:"resource_id" bson:"resource_id"`
	UniqueUsers  []bson.ObjectId `json:"profile_ids" bson:"profile_ids"`
	Identifier   string          `json:"identifier" bson:"identifier"`
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
	N := s.DB("manch").C("notification")
	N.Upsert(bson.M{"resource_id": rId, "type": t}, bson.M{
		"$pull": bson.M{"profile_ids": u},
	})
	fmt.Println("removed")
}

func CreateNotification(rId bson.ObjectId, t string, rT string, u bson.ObjectId) NotificationModel {
	s := session.Clone()
	defer s.Close()
	N := s.DB("manch").C("notifications")
	n := NotificationModel{
		Resource:     rId,
		Type:         t,
		ResourceType: rT,
		UniqueUsers:  []bson.ObjectId{u},
		Identifier:   GenerateIdentifier(rId, t),
	}
	count, _ := N.Find(bson.M{"resource_id": rId, "type": t}).Count()
	print(count)
	if count > 0 {
		N.Upsert(bson.M{"resource_id": rId, "type": t}, bson.M{
			"$addToSet": bson.M{"profile_ids": u},
		})
	} else {
		N.Insert(n)
	}
	return GetNotificationByResource(rId, t)

}

func GetNotificationByResource(rId bson.ObjectId, t string) NotificationModel {
	s := session.Clone()
	defer s.Close()
	N := s.DB("manch").C("notifications")
	notif := NotificationModel{}
	N.Find(bson.M{"resource_id": rId, "type": t}).One(&notif)
	return notif
}
