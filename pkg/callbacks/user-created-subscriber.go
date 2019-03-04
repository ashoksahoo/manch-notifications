package callbacks

import (
	"fmt"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"
	"time"
)

func UserCreatedSubscriberCB(subj, reply string, u *subscribers.User) {
	fmt.Printf("Received a New User created on subject %s! with User %+v\n", subj, u)
	// get user from db
	user := mongo.GetUserById(u.Id)

	// schedule welcome message
	scheduleTime := time.Now().Add(time.Duration(10) * time.Minute)
	whatsappSchedule := mongo.CreateWhatsAppSchedule(user, scheduleTime)
	fmt.Printf("whatsapp schedule \n%+v\n\n", whatsappSchedule)
	mongo.AddWhatsAppSchedule(whatsappSchedule)

	fmt.Printf("Processed a New User created on subject %s! with User Id %s\n", subj, u.Id)
}
