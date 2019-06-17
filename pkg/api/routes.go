package api

import (
	"github.com/go-chi/chi"
)

// with to use inline middleware

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/search", SearchHashTags)
	router.Get("/search/{id}", GetTagByID)
	router.Get("/users/search", HandleUserSearch)
	router.Get("/posts/search", HandlePostSearch)
	router.Post("/hashtag-weight/{id}", UpdateHashtagWeight)
	router.Get("/hashtag/image/{id}", GetHashTagImageById)
	router.Post("/hashtag/image/{id}", UpdateHashtagImage)
	router.Post("/hashtag/title/{id}", UpdateHashtagTitle)
	router.Get("/notifications", GetAllNotification)
	router.Get("/notifications/{id}", GetNotificationById)
	router.Post("/notifications/{id}", UpdateNotificationById)
	return router
}
