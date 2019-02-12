package api

import (
	"github.com/go-chi/chi"
)

// with to use inline middleware

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/notifications", GetAllNotification)
	router.Get("/notifications/{id}", GetNotificationById)
	router.Post("/notifications/{id}", UpdateNotificationById)
	return router
}
