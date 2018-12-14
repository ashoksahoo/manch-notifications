package firebase

import (
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

 */
type ManchMessage struct {
	Id          string `json:"mnc_nid,omitempty"`
	Namespace   string `json:"mnc_ns,omitempty"`
	Title       string `json:"mnc_nt,omitempty"`
	Message     string `json:"mnc_nm,omitempty"`
	Icon        string `json:"mnc_ico,omitempty"`
	DeepLink    string `json:"mnc_dl,omitempty"`
	Sound       string `json:"mnc_sound,omitempty"`
	BigPicture  string `json:"mnc_bp,omitempty"`
	BadgeIcon   string `json:"mnc_bi,omitempty"`
	BadgeCount  string `json:"mnc_bc,omitempty"`
	ChannelId   string `json:"mnc_cid,omitempty"`
	CollapseKey string `json:"mnc_ck,omitempty"`
	Priority    string `json:"mnc_pr,omitempty"`
	Actions     string `json:"mnc_acts,omitempty"`
	Silent      string `json:"mns_sn,omitempty"`
	MessageType string `json:"manch_message_type,omitempty"`
}

func MessageBuilder(m ManchMessage) map[string]string {
	var inInterface map[string]string
	inrec, _ := json.Marshal(m)
	json.Unmarshal(inrec, &inInterface)
	return inInterface
}

func SendMessage(m ManchMessage, token string) {
	// See documentation on defining a message payload.
	m.Namespace = "manch:N"
	if m.Icon == "" {
		m.Icon = "https://manch.app/img/new-logo.png"
	}
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
	if err != nil {
		//log.Fatalln(err)
		fmt.Println("Error:", err, token)
		mongo.DeleteToken(token)
	} else {
		// Response is a message ID string.
		fmt.Println("Successfully sent message:", response, token)
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
