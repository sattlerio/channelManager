package main

import (
	"github.com/jinzhu/gorm"
	"github.com/gorilla/mux"
	"os"
	"log"
	"net/http"
	_ "github.com/bmizerany/pq"
	"channelManager.sattler.io/helpers"
	"channelManager.sattler.io/models"
	"channelManager.sattler.io/controllers"
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
		helpers.Error.Println("not possible to connedssdzkdlksact to db, going to die now.... UAAAAAAAAAAAAH!!!!")
		panic("failed to connect database")
	}

	defer db.Close()

	db.AutoMigrate(models.Channels{})

	controllers.DbConn = db

	router.HandleFunc("/ping", controllers.PingController).Methods("GET")
	router.HandleFunc("/types/{company_id}", controllers.FetchAllChannels).Methods("GET")
	router.HandleFunc("/create/{type}/{channel}/{company_id}", controllers.CreateNewChannel).Methods("POST")
	router.HandleFunc("/fetch/{type}/{company_id}", controllers.FetchAllChannelsFromCompany).Methods("GET")
	router.HandleFunc("/delete/{type}/{company_id}/{channel_uuid}", controllers.DeleteChannelById).Methods("DELETE")
	router.HandleFunc("/get/{type}/{company_id}/{channel_uuid}", controllers.FetchChannelById).Methods("GET")
	router.HandleFunc("/edit/{type}/{channel}/{company_id}/{channel_uuid}", controllers.EditChannel).Methods("PUT")

	log.Fatal(http.ListenAndServe(":9000", router))
}