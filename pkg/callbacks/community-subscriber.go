package callbacks

import (
	"fmt"
	"notification-service/pkg/subscribers"
	"time"
	"notification-service/pkg/mongo"
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
		whatsappSchedule := mongo.CreateWhatsAppSchedule(user, scheduleTime)
		fmt.Printf("whatsapp schedule \n%+v\n\n", whatsappSchedule)
		mongo.AddWhatsAppSchedule(whatsappSchedule)
	}

	fmt.Printf("Processed a New Community on subject %s! with Id %s\n", subj, C.Id)

}
