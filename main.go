package main

import (
	"fmt"
	"net/http"
	"notification-service/pkg/callbacks"
	"notification-service/pkg/subscribers"

	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	subscribers.PostSubscriber(callbacks.PostSubscriberCB)

	subscribers.CommentSubscriber(callbacks.CommentSubscriberCB)

	subscribers.VotePostSubscriber(callbacks.VotePostSubscriberCB)

	subscribers.VoteCommentSubscriber(callbacks.VoteCommentSubscriberCB)

	subscribers.UserSubscriber(callbacks.UserSubscriberCB)

	subscribers.UserFollowSubscriber(callbacks.UserFollowSubscriberCB)

	subscribers.PostRemovedSubscriber(callbacks.PostRemovedSubscriberCB)

	subscribers.UserFollowRemovedSubscriber(callbacks.UserFollowRemovedSubscriberCB)

	http.ListenAndServe(":5000", r)
}

func init() {
	fmt.Println("Initializing Notification Service.")
}
