package api

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"networkCommunicationMin/db"
	"networkCommunicationMin/rest"
	"networkCommunicationMin/secondary_function"
)

func RegisterAPI(dbConnect *sql.DB) *chi.Mux {
	srv := rest.NewService(
		db.NewStorage(dbConnect),
	)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(secondary_function.Logger)

	r.Route("/", func(r chi.Router) {
		r.Post("/create", srv.Create)
		r.Post("/make_friends", srv.MakeFriends)
		r.Delete("/user", srv.Delete)
		r.Get("/friends/{id}", srv.GetFriends)
		r.Put("/user_id/{id}", srv.UpdateUserAge)
		r.Get("/get_all", srv.GetAll)
		r.Get("/ping", srv.Ping)
		r.Get("/", srv.Ping)
	})

	return r
}
