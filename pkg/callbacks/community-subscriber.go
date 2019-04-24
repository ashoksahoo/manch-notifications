package callbacks

import (
	"fmt"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"
	"notification-service/pkg/i18n"
	"time"
)

func CommunitySubscriberCB(subj, reply string, C *subscribers.Community) {
	fmt.Printf("Received a New Community on subject %s! with Community %+v\n", subj, C)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in subscribers.CommunityStatusUpdatedCB", r)
		}
	}()
	community := mongo.GetCommunityById(C.Id)
	for _, communityAdmin := range community.Admins {
		user := mongo.GetUserByProfileId(communityAdmin.ProfileId.Hex())
		scheduleTime := time.Now()
		data := i18n.DataModel{
			Community: community.Name,
		}
		templateName := "new_manch_whatsapp"
		message := i18n.GetString(user.Profiles[0].Language, templateName, data)
		whatsappSchedule := mongo.CreateWhatsAppSchedule(user, scheduleTime, message)
		fmt.Printf("whatsapp schedule \n%+v\n\n", whatsappSchedule)
		mongo.AddWhatsAppSchedule(whatsappSchedule)

		scheduleTime = time.Now().Add(30 * time.Minute)

		templateName = "manch_creator_competition"
		message = i18n.GetString(user.Profiles[0].Language, templateName, data)
		whatsappSchedule = mongo.CreateWhatsAppSchedule(user, scheduleTime, message)
		fmt.Printf("whatsapp schedule \n%+v\n\n", whatsappSchedule)
		mongo.AddWhatsAppSchedule(whatsappSchedule)
	}

	fmt.Printf("Processed a New Community on subject %s! with Id %s\n", subj, C.Id)

}
