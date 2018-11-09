package firebase

import (
	"firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"log"
	"notification-service/pkg/mongo"
)

var client *messaging.Client
var ctx context.Context

func SendMessage(m messaging.AndroidNotification, token string) {
	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data:         MessageBuilder(),
		Notification: nil,
		Android: &messaging.AndroidConfig{
			CollapseKey:           "",
			Priority:              "high",
			TTL:                   nil,
			RestrictedPackageName: "",
			Data:                  nil,
			Notification: &messaging.AndroidNotification{
				Title:        "Hello",
				Body:         "Hello World",
				Icon:         "",
				Color:        "",
				Sound:        "",
				Tag:          "",
				ClickAction:  "",
				BodyLocKey:   "",
				BodyLocArgs:  nil,
				TitleLocKey:  "",
				TitleLocArgs: nil,
			},
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

func MessageBuilder() map[string]string {

	return map[string]string{
		"mnc_ns":    "hello",
		"mnc_nt":    "world",
		"mnc_nm":    "jljsdkalj",
		"mnc_ico":   "kjahdkasjd",
		"mnc_dl":    "manch://posts/${data.post_id}",
		"mnc_sound": "true",
	}
}

func init() () {

	opt := option.WithCredentialsFile("./private/fcm.json")
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
