package mongo

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
)

type NotificationModel struct {
	Id       bson.ObjectId `json:"_id" bson:"_id"`
	Type     string        `json:"type" bson:"type"`
	Resource bson.ObjectId `json:"resource_id" bson:"resource_id"`
}

func CreateNotification(n NotificationModel) error {
	session := Session()
	N := session.DB("manch").C("notifications")
	count, err := N.Find(bson.M{"_id": n.Id}).Limit(1).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("resource %s already exists", n.Id)
	}
	return N.Insert(n)

}

func GetNotificationByResource(Id string) NotificationModel {
	session := Session()
	N := session.DB("manch").C("notifications")
	notif := NotificationModel{}
	N.Find(bson.M{"_id": bson.ObjectIdHex(Id)}).One(&notif)
	return notif
}
