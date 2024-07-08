package main

import (
	"fmt"
	"log"
	"time"

	"projekat/model"
	"projekat/service"
	"projekat/store"
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

	appservice.CreateDataItem(app1, &model.DataSpaceItem{Path: "app1/Root", Name: "fajl"})

	appservice.CreateSoftlink(app1, app2)

	appservice.ChangeDateSpaceState(*app1, model.Closed)

	app, err := db.GetApp("ns1", "app1")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(app.ApplicationId)
	items := nsService.RunDataDiscovery("ns2")
	for _, item := range items {
		fmt.Println(item)
	}

}
