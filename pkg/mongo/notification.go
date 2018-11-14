package mongo

import (
	"github.com/globalsign/mgo/bson"
)

type NotificationModel struct {
	Id           bson.ObjectId   `json:"_id,omitempty" bson:"_id,omitempty"`
	Type         string          `json:"type" bson:"type"`
	ResourceType string          `json:"resource_type" bson:"resource_type"`
	Resource     bson.ObjectId   `json:"resource_id" bson:"resource_id"`
	UniqueUsers  []bson.ObjectId `json:"profile_ids" bson:"profile_ids"`
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
	}
	count, _ := N.Find(bson.M{"resource_id": rId}).Limit(1).Count()
	if count > 0 {
		N.Upsert(bson.M{"resource_id": rId, "type": t}, bson.M{
			"$addToSet": bson.M{"profile_ids": u},
		})
	}
	N.Insert(n)
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
