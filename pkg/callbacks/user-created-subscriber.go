package callbacks

import (
	"fmt"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"
	"notification-service/pkg/i18n"
	"time"
)

func UserCreatedSubscriberCB(subj, reply string, u *subscribers.User) {
	fmt.Printf("Received a New User created on subject %s! with User %+v\n", subj, u)
	// get user from db
	user := mongo.GetUserById(u.Id)

	// schedule welcome message
	scheduleTime := time.Now().Add(time.Duration(10) * time.Minute)
	userCoinScheduleTime := time.Now().Add(time.Duration(15) * time.Minute)

	data := i18n.DataModel{
		Name: "",
	}
	templateName := "welcome_message"
	message := i18n.GetString(user.Profiles[0].Language, templateName, data)
	whatsappSchedule := mongo.CreateWhatsAppSchedule(user, scheduleTime, message, "TEXT")
	fmt.Printf("whatsapp schedule \n%+v\n\n", whatsappSchedule)
	mongo.AddWhatsAppSchedule(whatsappSchedule)

	mongo.CreateUserCoinSchedule(mongo.UserCoinsModelScheduleModel{
		ProfileId:   user.Profiles[0].Id,
		Action:      "whatsapp_welcome",
		CoinsEarned: 50,
	}, userCoinScheduleTime)
	fmt.Printf("Processed a New User created on subject %s! with User Id %s\n", subj, u.Id)
}
