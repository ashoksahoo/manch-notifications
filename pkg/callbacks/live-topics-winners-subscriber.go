package callbacks

import (
	"fmt"
	"notification-service/pkg/constants"
	"notification-service/pkg/firebase"
	"notification-service/pkg/i18n"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"
	"notification-service/pkg/utils"

	"github.com/globalsign/mgo/bson"
)

func LiveTopicsWinnerSubscriberCB(subj, reply string, W *subscribers.LiveTopicsWinner) {
	fmt.Printf("Received a live topics winners on subject %s! with Value %+v\n", subj, W.Id)
	winners := W.Winners
	participantsIds := W.Participants

	topicTitle := W.Title

	fmt.Println("topic title", topicTitle)
	fmt.Println("winners is", winners)
	fmt.Println("participants is", participantsIds)

	coinsEarned := 500

	winnersProfiles := mongo.GetProfilesByIds(winners)
	participantsProfiles := mongo.GetProfilesByIds(participantsIds)

	for i, winnerId := range winners {
		mongo.CreateUserCoin(mongo.UserCoinsModel{
			ProfileId:   bson.ObjectIdHex(winnerId),
			CoinsEarned: coinsEarned,
			Action:      "live-topics.winner",
		})


		winner := winnersProfiles[i]

		tokens := mongo.GetTokensByProfiles([]bson.ObjectId{bson.ObjectIdHex(winnerId)})


		entities := []mongo.Entity{
			{
				EntityId:   bson.ObjectIdHex(W.Id),
				EntityType: "live_topic",
			},
		}
		notification := mongo.CreateNotification(mongo.NotificationModel{
			Receiver:        bson.ObjectIdHex(winnerId),
			Identifier:      winnerId + "_live_topic_winner",
			Participants:    []bson.ObjectId{bson.ObjectIdHex(winnerId)},
			DisplayTemplate: constants.NotificationTemplate["TRANSACTIONAL"],
			EntityGroupId:   W.Id,
			ActionId:        bson.ObjectIdHex(W.Id),
			ActionType:      "live_topic_winner",
			Purpose:         constants.NotificationPurpose["LIVE_TOPIC_WINNER"],
			Entities:        entities,
			NUUID:           "",
		})

		data := i18n.DataModel{
			Name:  winner.Name,
			Count: coinsEarned,
		}

		notifTitles := []string{"live_topic_winners_title_1", "live_topic_winners_title_1"}
		randomIndex := utils.Random(0, 2)

		notifTitleTemplate := notifTitles[randomIndex]
		msgStrTitle := i18n.GetString(winner.Language, notifTitleTemplate, data)
		msgStrText := i18n.GetString(winner.Language, "live_topics_winner_text", data)

		messageMeta := mongo.MessageMeta{
			Template: notifTitleTemplate,
			Data:     data,
		}
		deepLink := "manch://live/top/" + W.Id
		// update notification message
		mongo.UpdateNotification(bson.M{"_id": notification.Id}, bson.M{
			"message":      msgStrTitle,
			"message_meta": messageMeta,
		})

		msg := firebase.ManchMessage{
			Title:    msgStrTitle,
			Message:  msgStrText,
			DeepLink: deepLink,
			Id:       notification.NId,
		}

		fmt.Printf("\nGCM Message %+v\n", msg)
		if tokens != nil {
			for _, token := range tokens {
				go firebase.SendMessage(msg, token.Token, notification)
			}
		} else {
			fmt.Printf("No token\n")
		}

	}

	for i, participantId := range participantsIds {

		participant := participantsProfiles[i]

		tokens := mongo.GetTokensByProfiles([]bson.ObjectId{bson.ObjectIdHex(participantId)})

		entities := []mongo.Entity{
			{
				EntityId:   bson.ObjectIdHex(W.Id),
				EntityType: "live_topic",
			},
		}
		notification := mongo.CreateNotification(mongo.NotificationModel{
			Receiver:        bson.ObjectIdHex(participantId),
			Identifier:      participantId + "_live_topic_winner",
			Participants:    []bson.ObjectId{bson.ObjectIdHex(participantId)},
			DisplayTemplate: constants.NotificationTemplate["TRANSACTIONAL"],
			EntityGroupId:   W.Id,
			ActionId:        bson.ObjectIdHex(W.Id),
			ActionType:      "live_topic_winner",
			Purpose:         constants.NotificationPurpose["LIVE_TOPIC_WINNER"],
			Entities:        entities,
			NUUID:           "",
		})

		data := i18n.DataModel{
			Name:  participant.Name,
			Count: coinsEarned,
		}

		notifTitles := []string{"live_topic_participants_title_1", "live_topic_participants_title_2"}
		randomIndex := utils.Random(0, 2)

		notifTitleTemplate := notifTitles[randomIndex]
		msgStrTitle := i18n.GetString(participant.Language, notifTitleTemplate, data)

		messageMeta := mongo.MessageMeta{
			Template: notifTitleTemplate,
			Data:     data,
		}
		deepLink := "manch://live/top/" + W.Id
		// update notification message
		mongo.UpdateNotification(bson.M{"_id": notification.Id}, bson.M{
			"message":      msgStrTitle,
			"message_meta": messageMeta,
		})

		msg := firebase.ManchMessage{
			Title:    msgStrTitle,
			Message:  topicTitle,
			DeepLink: deepLink,
			Id:       notification.NId,
		}

		fmt.Printf("\nGCM Message %+v\n", msg)
		if tokens != nil {
			for _, token := range tokens {
				go firebase.SendMessage(msg, token.Token, notification)
			}
		} else {
			fmt.Printf("No token\n")
		}

	}

	fmt.Printf("Processed a live topics winners on subject %s! with Id %s\n", subj, W.Id)
}
