package main

import (
	"fmt"
	"net/http"
	"notification-service/pkg/callbacks"
	"notification-service/pkg/subscribers"
	"notification-service/pkg/utils"

	"notification-service/pkg/api"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON), // set content type header as json
		middleware.Logger,          // Log API request calls
		middleware.DefaultCompress, // Compress results
		middleware.RedirectSlashes, // redirect slashes to no slashes
		middleware.Recoverer,       // Recover from panics without crashing server
		middleware.URLFormat,
	)

	router.Route("/", func(r chi.Router) {
		r.Mount("/", api.Routes())
	})

	return router
}

func main() {
	router := Routes()
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
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

	// received post deleted event
	subscribers.PostDeletedSubscriber(callbacks.PostDeletedSubscriberCB)

	// received post moderation event
	subscribers.PostModeratedSubscriber(callbacks.PostModeratedSubscriberCB)

	// received a Share event
	subscribers.ShareSubscriber(callbacks.ShareSubscriberCB)

	// listen on http server 5000
	http.ListenAndServe(":5000", router)
}

func init() {
	fmt.Println("Initializing Notification Service.")
}
