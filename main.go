package main

import (
	"firebase.google.com/go/messaging"
	"fmt"
	"github.com/globalsign/mgo/bson"
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
		newPost := mongo.GetPostById(bson.ObjectIdHex(p.Id))
		fmt.Printf("Mongo return for Post %+v\n", newPost)

	})
	subscribers.CommentSubscriber(func(subj, reply string, cmt *subscribers.Comment) {
		comment := mongo.GetCommentById(cmt.Id)
		fmt.Printf("Created By %+v\n", comment.Created)

		profiles := []bson.ObjectId{comment.Created.ProfileId}
		tokens := mongo.GetTokensByProfiles(profiles)
		msg := messaging.AndroidNotification{
			Title:        "New Comment",
			Body:         comment.Created.Name + " commented on Your Post",
			Icon:         comment.Created.Avatar,
			Color:        "",
			Sound:        "true",
			Tag:          "",
			ClickAction:  "",
			BodyLocKey:   "",
			BodyLocArgs:  nil,
			TitleLocKey:  "",
			TitleLocArgs: nil,
		}
		firebase.SendMessage(msg, "c3H8bqZDN2M:APA91bGR8azrCNJwgygmVKb42kC_4_PlVq28IeI5i5217vHEKNIWd3AMfYojERdgvkHvQxTU3VGfPmpJoM4e7u_HXEUyf6fB0Nfc1Ey-he20uVrOyzv4cefIVjbeC02co3zM4FUFaKUj")
		if tokens != nil {
			for _, token := range tokens {
				go firebase.SendMessage(msg, token.Token)
			}
		} else {
			fmt.Printf("No token")
		}
	})

	http.ListenAndServe(":5000", r)
}

func init() {
	fmt.Println("Initializing Notification Service.")
}
