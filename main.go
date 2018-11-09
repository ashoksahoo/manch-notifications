package main

import (
	"firebase.google.com/go/messaging"
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
	msg := messaging.AndroidNotification{
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
	}

	subscribers.CommentSubscriber(func(subj, reply string, cmt *subscribers.Comment) {
		fmt.Printf("Received a comment on subject %s! with Comment ID %s\n", subj, cmt.Id)
		comment := mongo.GetCommentById(cmt.Id)
		fmt.Printf("Mongo Query return for comment %+v\n", comment)

		userAshok := []string{"5b470fa4c15b665c96950f19"}
		tokens := mongo.GetTokensByProfiles(userAshok)
		if tokens != nil {
			for _, token := range tokens {
				go firebase.SendMessage(msg, token.Token)
			}
		} else {
			fmt.Printf("No token")
		}
	})
	firebase.SendMessage(msg, "c3H8bqZDN2M:APA91bGR8azrCNJwgygmVKb42kC_4_PlVq28IeI5i5217vHEKNIWd3AMfYojERdgvkHvQxTU3VGfPmpJoM4e7u_HXEUyf6fB0Nfc1Ey-he20uVrOyzv4cefIVjbeC02co3zM4FUFaKUj")

	http.ListenAndServe(":5000", r)
}

func init() {
	fmt.Println("Initializing Notification Service.")
}
