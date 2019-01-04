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

	// received a post
	subscribers.PostSubscriber(callbacks.PostSubscriberCB)

	// received a comment
	subscribers.CommentSubscriber(callbacks.CommentSubscriberCB)

	// received a vote on post
	subscribers.VotePostSubscriber(callbacks.VotePostSubscriberCB)

	// received a vote on comment
	subscribers.VoteCommentSubscriber(callbacks.VoteCommentSubscriberCB)

	// received a new user
	subscribers.UserSubscriber(callbacks.UserSubscriberCB)

	// received user-follow
	subscribers.UserFollowSubscriber(callbacks.UserFollowSubscriberCB)

	// received post removed (soft delete)
	subscribers.PostRemovedSubscriber(callbacks.PostRemovedSubscriberCB)

	// received user-follow removed
	subscribers.UserFollowRemovedSubscriber(callbacks.UserFollowRemovedSubscriberCB)

	// received re-post event
	subscribers.RepostSubscriber(callbacks.RepostSubscriberCB)

	// received share-post event
	subscribers.SharePostSubscriber(callbacks.SharePostSubscriberCB)


	// listen on http server 5000
	http.ListenAndServe(":5000", r)
}

func init() {
	fmt.Println("Initializing Notification Service.")
}
