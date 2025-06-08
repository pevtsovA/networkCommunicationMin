package main

import (
	"flag"
	"net/http"
	"networkCommunicationMin/api"
	"networkCommunicationMin/db"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	port := flag.String("port", "8080", "Listen server port")
	password := flag.String("dbpassword", "", "Database password")
	isDocker := flag.Bool("runviadocker", false, "Run the application")
	flag.Parse()

	dbConnect, closeConnect := db.ConnectToBD(*password, *isDocker)
	defer closeConnect()

	r := api.RegisterAPI(dbConnect)

	log.Info("listening localhost: ", *port)
	if err := http.ListenAndServe(":"+*port, r); err != nil {
		log.Fatalln("listenAndServe server: ", err)
	}
}
