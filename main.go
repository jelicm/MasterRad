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

	service := service.NewApplicationService(db)

	app1, _ := service.RunApplication("app1", "ns1", 25)

	app2, _ := service.RunApplication("app2", "ns2", 100)

	service.CreateDataItem(app1, &model.DataSpaceItem{Path: "app1/Root", Name: "fajl"})

	service.CreateSoftlink(app1, app2)

	service.ChangeDateSpaceState(*app1, model.Closed)

	app, err := db.GetApp("ns1", "app1")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(app.ApplicationId)

}
