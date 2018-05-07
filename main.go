package main

import (
	"github.com/jinzhu/gorm"
	"github.com/gorilla/mux"
	"os"
	"log"
	"net/http"
	_ "github.com/bmizerany/pq"
	"channelManager/helpers"
	"channelManager/models"
	"channelManager/controllers"
)

var db *gorm.DB
var err error

func main() {
	helpers.Init(os.Stderr, os.Stdout, os.Stdout, os.Stderr)
	helpers.Info.Println("starting up server...")

	router := mux.NewRouter()

	db, err = gorm.Open(
		"postgres",
		"host="+os.Getenv("PSQL_HOST")+" user="+os.Getenv("PSQL_USER")+
			" dbname="+os.Getenv("PSQL_DBNAME")+" sslmode=disable password="+
			os.Getenv("PSQL_PASSWORD"))

	if err != nil {
		helpers.Error.Println(err)
		helpers.Error.Println("not possible to connect to db, going to die now.... UAAAAAAAAAAAAH!!!!")
		panic("failed to connect database")
	}

	defer db.Close()

	db.AutoMigrate(models.Channels{})

	controllers.DbConn = db

	router.HandleFunc("/ping", controllers.PingController).Methods("GET")
	router.HandleFunc("/types/{company_id}", controllers.FetchAllChannels).Methods("GET")

	log.Fatal(http.ListenAndServe(":11000", router))
}