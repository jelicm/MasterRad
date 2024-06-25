package main

import (
	"fmt"
	"log"
	"time"

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

	service.RunApplication("app1", "ns1", 25)

	app, err := db.GetApp("ns1", "app1")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(app.ApplicationId)

}
