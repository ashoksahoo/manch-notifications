package firebase

import (
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

func SendMessage(m messaging.AndroidNotification, token string, id string) {
	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data:         MessageBuilder(m, id),
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

func MessageBuilder(m messaging.AndroidNotification, id string) map[string]string {

	return map[string]string{
		"mnc_ns":    "manch:N",
		"mnc_nt":    m.Title,
		"mnc_nm":    m.Body,
		"mnc_ico":   m.Icon,
		"mnc_dl":    "manch://posts/" + id,
		"mnc_sound": "true",
	}
}

func init() () {
	var filename string
	switch os.Getenv("env") {
	case "staging":
		filename = "./private/fcm_staging.json"
		break
	case "development":
	case "production":
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
