package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/nats-io/go-nats"
	"net/http"
	"notification-service/pkg/mongo"
	"notification-service/pkg/subscribers"
)

func main() {
	r := chi.NewRouter()
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	nc, _ := nats.Connect(nats.DefaultURL)
	c, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)

	subscribers.PostSubscriber(c, func(subj, reply string, p *subscribers.Post) {
		fmt.Printf("Received a post on subject %s! with Post ID %s\n", subj, p.Id)
		newPost := mongo.GetPostById(p.Id)
		fmt.Printf("Mongo Query return for Post %+v\n", newPost)

	})

	subscribers.CommentSubscriber(c, func(subj, reply string, cmt *subscribers.Comment) {
		fmt.Printf("Received a comment on subject %s! with Comment ID %s\n", subj, cmt.Id)
		comment := mongo.GetCommentById(cmt.Id)
		fmt.Printf("Mongo Query return for comment %+v\n", comment)

	})

	http.ListenAndServe(":5000", r)
	defer c.Close()
}
