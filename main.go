package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"projekat/handler"
	"projekat/model"
	"projekat/service"
	"projekat/store"

	"github.com/gorilla/mux"
)

func main() {

	endpoints := []string{"localhost:2379"}
	timeout := 5 * time.Second

	db, err := store.New(endpoints, timeout)

	if err != nil {
		log.Fatal(err)
	}

	appservice := service.NewApplicationService(db)
	nsService := service.NewNamespaceService(db)

	app1, _ := appservice.RunApplication("app1", "ns1", 25)

	app2, _ := appservice.RunApplication("app2", "ns2", 100)

	appservice.CreateDataItem(app1.ApplicationId, &model.DataSpaceItem{Path: "app1/Root", Name: "fajl", SizeKB: 1}, "nekasema", false)

	appservice.CreateSoftlink(app1, app2)

	//appservice.ChangeDateSpaceState(*app1, model.Closed)

	app, err := db.GetApp("ns1", "app1")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(app.ApplicationId)
	items := nsService.RunDataDiscovery("ns1")
	for _, item := range items {
		fmt.Println(item)
	}

	appHandler := handler.NewAppHandler(appservice, nsService)

	nsService.DeleteAppDefault("ns1", "app1")

	r := mux.NewRouter()

	r.HandleFunc("/runApp", appHandler.RunApp).Methods("POST")
	r.HandleFunc("/dataDiscovery/{nsId}", appHandler.RunDataDiscovery).Methods("GET")
	r.HandleFunc("/addDataItem", appHandler.AddDataItem).Methods("POST")
	r.HandleFunc("/deleteApp", appHandler.DeleteApp).Methods("DELETE")
	srv := &http.Server{
		Handler: r,
		Addr:    ":8001",
	}
	log.Fatal(srv.ListenAndServe())

}
