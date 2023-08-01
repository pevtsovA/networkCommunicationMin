package main

import (
	"database/sql"
	"flag"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"networkCommunicationMin/dataBase"
	"networkCommunicationMin/models"
)

var srv = dataBase.Service{}

func main() {
	port := flag.String("port", "8080", "Listen server port")
	flag.Parse()

	connStr := "user=postgres password=******** dbname=socnetworkdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	srv.DB = db

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(models.Logger)

	r.Route("/", func(r chi.Router) {
		r.Post("/create", srv.Create)
		r.Post("/make_friends", srv.MakeFriends)
		r.Delete("/user", srv.Delete)
		r.Get("/friends/{id}", srv.GetFriends)
		r.Put("/user_id/{id}", srv.UpdateUser)
		r.Get("/get_all", srv.GetAll)
		r.Get("/ping", srv.Ping)
		r.Get("/", srv.Ping)
	})

	log.Println("listening localhost:" + *port)
	if err := http.ListenAndServe("localhost:"+*port, r); err != nil {
		log.Fatalln(err)
	}
}
