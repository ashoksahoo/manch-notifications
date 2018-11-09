package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"notification-service/pkg/firebase"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"
)

func main() {
	r := chi.NewRouter()
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	subscribers.PostSubscriber(func(subj, reply string, p *subscribers.Post) {
		fmt.Printf("Received a post on subject %s! with Post ID %s\n", subj, p.Id)
		newPost := mongo.GetPostById(p.Id)
		fmt.Printf("Mongo Query return for Post %+v\n", newPost)

	})

	subscribers.CommentSubscriber(func(subj, reply string, cmt *subscribers.Comment) {
		fmt.Printf("Received a comment on subject %s! with Comment ID %s\n", subj, cmt.Id)
		comment := mongo.GetCommentById(cmt.Id)
		fmt.Printf("Mongo Query return for comment %+v\n", comment)

		userAshok := []string{"5b470fa4c15b665c96950f19"}
		tokens := mongo.GetTokensByProfiles(userAshok)
		if tokens != nil {
			for _, token := range tokens {
				go firebase.SendMessage(token.Token)
			}
		} else {
			fmt.Printf("No token")
		}
	})
	firebase.SendMessage("fpxicDsbUko:APA91bFnp_B4hu8wdEMO2Ad4LXhLk7hts9CDUWfCIwdTJ63XpG0QEqlCpvXAeZQJQ7VG-9Wucf38p75H1qRILUzCZkiaMDR30SJoKgDJS8Ggm3Y7qywGvsR2EXwPC1SQ2vaVRx9SHyId")
	firebase.SendMessage("c3H8bqZDN2M:APA91bGR8azrCNJwgygmVKb42kC_4_PlVq28IeI5i5217vHEKNIWd3AMfYojERdgvkHvQxTU3VGfPmpJoM4e7u_HXEUyf6fB0Nfc1Ey-he20uVrOyzv4cefIVjbeC02co3zM4FUFaKUj")
	firebase.SendMessage("eow3qWbmKlc:APA91bE_9vQzdZCViUk-6DS-1QEOIH64J78swU3VGihzOSz1pr7RUV_5sQW1ETyfFjOBR9OnIDrMf0pv4IOCnPCMxCyCBGxk7p8PMBQJ9GTR55vvPYdoFAXEwB_7FYyRznAjLwx35a9v")
	firebase.SendMessage("eow3qWbmKlc:APA91bE_9vQzdZCViUk-6DS-1QEOIH64J78swU3VGihzOSz1pr7RUV_5sQW1ETyfFjOBR9OnIDrMf0pv4IOCnPCMxCyCBGxk7p8PMBQJ9GTR55vvPYdoFAXEwB_7FYyRznAjLwx35a9v")

	http.ListenAndServe(":5000", r)
}

func init() {
	fmt.Println("Initializing Notification Service.")
}
