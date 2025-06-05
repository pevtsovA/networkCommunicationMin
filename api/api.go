package api

import (
	"database/sql"
	"networkCommunicationMin/db"
	"networkCommunicationMin/rest"
	secondary "networkCommunicationMin/secondary_function"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func RegisterAPI(dbConnect *sql.DB) *chi.Mux {
	srv := rest.NewService(
		db.NewStorage(dbConnect),
	)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(secondary.Logger)

	r.Route("/", func(r chi.Router) {
		r.Post("/create", srv.Create)
		r.Post("/make_friends", srv.MakeFriends)
		r.Delete("/user/{id}", srv.Delete)
		r.Get("/friends/{id}", srv.GetFriends)
		r.Put("/user/{id}", srv.UpdateUserAge)
		r.Get("/get_all", srv.GetAll)
		r.Get("/ping", srv.Ping)
	})

	return r
}
