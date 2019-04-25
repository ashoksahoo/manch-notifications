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
	TemplateName string         `json:"template_name" bson:"template_name"`
	Template     string         `json:"template" bson:"template"`
	Data         i18n.DataModel `json:"data" bson:"data"`
}

type PushMeta struct {
	Status     string    `json:"status" bson:"status"`
	PushId     string    `json:"push_id" bson:"push_id"`
	FailReason string    `json:"error" bson:"error"`
	CreatedAt  time.Time `json:"date" bson:"date"`
}

type PurposeIcon struct {
	ResourceId string `json:"resource_id" bson:"resource_id"`
	IconUrl    string `json:"icon_url" bson:"icon_url"`
}

type NotificationModel struct {
	Id                    bson.ObjectId   `json:"_id,omitempty" bson:"_id,omitempty"`
	Receiver              bson.ObjectId   `json:"receiver" bson:"receiver"`
	Identifier            string          `json:"identifier" bson:"identifier"`
	Message               string          `json:"message" bson:"message"`
	IsRead                bool            `json:"is_read" bson:"is_read"`
	Participants          []bson.ObjectId `json:"participants" bson:"participants"`
	DisplayTemplate       string          `json:"display_template" bson:"display_template"`
	EntityGroupId         string          `json:"entity_group_id" bson:"entity_group_id"`
	ActionType            string          `json:"action_type" bson:"action_type"`
	ActionId              bson.ObjectId   `json:"action_id" bson:"action_id"`
	NId                   string          `json:"n_id" bson:"n_id"`
	NUUID                 string          `json:"nuuid" bson:"nuuid"`
	Purpose               string          `json:"purpose" bson:"purpose"`
	Entities              []Entity        `json:"entities" bson:"entities"`
	MessageMeta           MessageMeta     `json:"message_meta" bson:"message_meta"`
	Push                  PushMeta        `json:"push" bson:"push"`
	CreatedAt             time.Time       `json:"createdAt" bson:"createdAt"`
	UpdatedAt             time.Time       `json:"updatedAt" bson:"updatedAt"`
	Delivered             bool            `json:"delivered" bson:"delivered"`
	DeepLink              string          `json:"deep_link" bson:"deep_link"`
	MessageHtml           string          `json:"message_html" bson:"message_html"`
	PlaceHolderIcon       []string        `json:"place_holder_icons" bson:"place_holder_icons"`
	PurposeIcon           PurposeIcon     `json:"purpose_icon" bson:"purpose_icon"`
	PushType              string          `json:"push_type" bson:"push_type"`
	NotificationUpdatedAt time.Time       `json:"notification_updated_at" bson:"notification_updated_at"`
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
		Receiver:              notification.Receiver,
		Identifier:            notification.Identifier,
		Message:               notification.Message,
		IsRead:                false,
		Participants:          notification.Participants,
		DisplayTemplate:       notification.DisplayTemplate,
		EntityGroupId:         notification.EntityGroupId,
		ActionId:              notification.ActionId,
		ActionType:            notification.ActionType,
		NId:                   GenerateIdentifier(notification.ActionId, notification.ActionType),
		Purpose:               notification.Purpose,
		Entities:              notification.Entities,
		MessageMeta:           notification.MessageMeta,
		MessageHtml:           notification.MessageHtml,
		DeepLink:              notification.DeepLink,
		Push:                  push,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		Delivered:             false,
		NotificationUpdatedAt: time.Now(),
	}

	if notification.PushType == "" {
		n.PushType = "manch:N"
	} else {
		n.PushType = notification.PushType
	}

	// get participants avatar
	placeHolderIcons := []string{}
	participantsIdsString := []string{}
	// convert objectid to string for participants
	for _, participantId := range notification.Participants {
		participantsIdsString = append(participantsIdsString, participantId.Hex())
	}

	participants := GetProfilesByIds(participantsIdsString)
	for _, participant := range participants {
		placeHolderIcons = append(placeHolderIcons, participant.Avatar)
	}

	n.PlaceHolderIcon = placeHolderIcons
	// notification purpose icon
	n.PurposeIcon = PurposeIcon{
		ResourceId: constants.NotificationPurposeResource[n.Purpose],
		IconUrl:    constants.NotificationPurposeIcon[n.Purpose],
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
			"$set":         bson.M{"updatedAt": time.Now(), "notification_updated_at": time.Now(), "place_holder_icons": n.PlaceHolderIcon},
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
	fmt.Printf("notification query is %+v", query)
	fmt.Printf("update is %+v", update)

	s := session.Clone()
	defer s.Close()
	N := s.DB("manch").C(NOTIFICATION_V2_MODEL)
	err := N.Update(query, bson.M{"$set": update})
	if err != nil {
		fmt.Println("Notification update error", err)
	}
}

func GetNotificationByIdentifier(identifier string) (error, NotificationModel) {
	s := session.Clone()
	defer s.Close()
	N := s.DB("manch").C(NOTIFICATION_V2_MODEL)
	notif := NotificationModel{}
	err := N.Find(bson.M{"identifier": identifier, "is_read": false}).One(&notif)
	return err, notif
}

func GetNotificationByQuery(query bson.M) (error, []NotificationModel) {
	s := session.Clone()
	defer s.Close()
	fmt.Println("Query is ", query)
	N := s.DB("manch").C(NOTIFICATION_V2_MODEL)
	notifications := []NotificationModel{}
	limit, limitOk := query["limit"]
	delete(query, "limit")
	sort, sortOk := query["sort"]
	delete(query, "sort")
	skip, skipOk := query["skip"]
	delete(query, "skip")
	queryResult := N.Find(query)
	fmt.Println("sort is ", sort)
	if sortOk {
		queryResult.Sort(sort.([]string)...)
	}
	if skipOk {
		queryResult.Skip(skip.(int))
	}
	if limitOk {
		queryResult.Limit(limit.(int))
	}
	err := queryResult.Select(bson.M{
		"participants":    0,
		"entity_group_id": 0,
		"entities":        0,
	}).All(&notifications)
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
