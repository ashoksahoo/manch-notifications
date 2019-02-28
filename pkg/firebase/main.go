package firebase

import (
	"time"
	"notification-service/pkg/constants"
	"github.com/globalsign/mgo/bson"
	"encoding/json"
	"firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"log"
	"notification-service/pkg/mongo"
	"os"
)

var client *messaging.Client
var ctx context.Context

/**
String MSG_NOTIFICATION_ID = "mnc_nid";
String MSG_SILENT = "mnc_sn";
String MSG_NOTIFICATION_SOUND = "mnc_sound";
String MSG_TITLE = "mnc_nt";
String MSG_NOTIFICATION_MESSAGE = "mnc_nm";
String MSG_NOTIFICATION_ICON_PATH = "mnc_ico";
String MSG_BIG_PICTURE = "mnc_bp";
String MSG_BADGE_ICON = "mnc_bi";
String MSG_BADGE_COUNT = "mnc_bc";
String MSG_CHANNEL_ID = "mnc_cid";
String MSG_COLLAPSE_KEY = "mnc_ck";
String MSG_PRIORITY = "mnc_pr";
String MSG_ACTIONS = "mnc_acts";
// name space of the notifiction :  has to be manch:D(for data messages) or manch:N
String MSG_NAME_SPACE = "mnc_ns";
// notification "type" attribute  [ P(promotional) T(Transactional) ]
String MSG_TYPE = "mnc_at";
// notification purpose attribute  [ C(campaign), TC(txn comment), TL(Txn Like) etc.   ]
String MSG_PURPOSE = "mnc_ap";
// notification campaign id
String MSG_CAMPAIGN_ID = "mnc_acid";
// notification AB test variations
String MSG_AB_TEST_VARIATION = "mnc_atv";
// notification date attribute
String MSG_DATE = "mnc_ad";
 */
type ManchMessage struct {
	Id          	string `json:"mnc_nid,omitempty"`
	Namespace   	string `json:"mnc_ns,omitempty"`
	Title       	string `json:"mnc_nt,omitempty"`
	Message     	string `json:"mnc_nm,omitempty"`
	Icon        	string `json:"mnc_ico,omitempty"`
	DeepLink    	string `json:"mnc_dl,omitempty"`
	Sound       	string `json:"mnc_sound,omitempty"`
	BigPicture  	string `json:"mnc_bp,omitempty"`
	BadgeIcon   	string `json:"mnc_bi,omitempty"`
	BadgeCount  	string `json:"mnc_bc,omitempty"`
	ChannelId   	string `json:"mnc_cid,omitempty"`
	CollapseKey 	string `json:"mnc_ck,omitempty"`
	Priority    	string `json:"mnc_pr,omitempty"`
	Actions     	string `json:"mnc_acts,omitempty"`
	Silent      	string `json:"mns_sn,omitempty"`
	MNCID 			string `json:"mnc_id" bson:"mnc_id"`
	MessageType 	string `json:"mnc_at" bson:"mnc_at"`
	Purpose 		string `json:"mnc_ap" bson:"mnc_ap"`
	CampaignId 		string `json:"mnc_acid" bson:"mnc_acid"`
	TestVariation 	string `json:"mnc_atv" bson:"mnc_atv"`
	Date 			string `json:"mnc_ad" bson:"mnc_ad"`
	Status 			string `json:"status" bson:"status"`
	BlockedTill 	string `json:"blocked_till" bson:"blocked_till"`
	BlockedOn 		string `json:"blocked_on" bson:"blocked_on"`
	Reason 			string `json:"reason" bson:"reason"`
	LastWarned 		string `json:"last_warned_on" bson:"last_warned_on"`
}

func MessageBuilder(m ManchMessage) map[string]string {
	var inInterface map[string]string
	inrec, _ := json.Marshal(m)
	json.Unmarshal(inrec, &inInterface)
	return inInterface
}

func SendMessage(m ManchMessage, token string, notification mongo.NotificationModel) {
	// See documentation on defining a message payload.
	if m.Namespace == "" {
		m.Namespace = "manch:N"
	}
	if m.Icon == "" {
		m.Icon = "https://manch.app/img/new-logo.png"
	}

	if notification.DisplayTemplate == constants.NotificationTemplate["TRANSACTIONAL"] {
		m.MessageType = "T"
	}

	m.Purpose = notification.Purpose
	m.Date = time.Now().String()

	m.MNCID = notification.NUUID
	message := &messaging.Message{
		Data:         MessageBuilder(m),
		Notification: nil,
		Android: &messaging.AndroidConfig{
			CollapseKey:           "",
			Priority:              "high",
			TTL:                   nil,
			RestrictedPackageName: "",
			Data:                  nil,
			Notification:          nil,
		},
		Webpush:   nil,
		APNS:      nil,
		Token:     token,
		Topic:     "",
		Condition: "",
	}

	// Send a message to the device corresponding to the provided
	// registration token.
	response, err := client.Send(ctx, message)

	fmt.Println("response from firebase is:", response)
	fmt.Println("error is:", err)

	if err != nil {
		//log.Fatalln(err)
		fmt.Println("Error:", err, token)
		// delete token
		mongo.DeleteToken(token)
		// update push info
		fmt.Println("update")
		mongo.UpdateNotification(bson.M{"_id": notification.Id}, bson.M{
			"push": mongo.PushMeta{
				Status: constants.NotificationStatus["FAILED"],
				FailReason: err.Error(),
				CreatedAt: time.Now(),
			},
		})
	} else {
		// Response is a message ID string.
		fmt.Println("Successfully sent message:", response, token)
		// update push info
		mongo.UpdateNotification(bson.M{"_id": notification.Id}, bson.M{
			"push": mongo.PushMeta{
				Status: constants.NotificationStatus["SENT"],
				PushId: response,
				CreatedAt: time.Now(),
			},
		})
	}

}

func init() () {
	var filename string
	switch os.Getenv("env") {
	case "staging":
		filename = "./private/fcm_staging.json"
	case "development":
		filename = "./private/fcm.json"
	case "production":
		filename = "./private/fcm.json"
	default:
		filename = "./private/fcm.json"
	}
	opt := option.WithCredentialsFile(filename)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	ctx = context.Background()
	if client, err = app.Messaging(ctx); err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	} else {
		fmt.Println("Initialized GCM.")
	}

}
