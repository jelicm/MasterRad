package main

import (
	"log"
	"net/http"
	"time"

	"projekat/handler"
	"projekat/service"
	"projekat/store"

	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
)

func main() {

	endpoints := []string{"localhost:2379"}
	timeout := 5 * time.Second

	conn := Conn()
	defer conn.Close()

	db, err := store.New(endpoints, timeout)

	if err != nil {
		log.Fatal(err)
	}

	appservice := service.NewApplicationService(db, conn)
	nsService := service.NewNamespaceService(db)

	appHandler := handler.NewAppHandler(appservice, nsService)

	r := mux.NewRouter()

	r.HandleFunc("/runApp", appHandler.RunApp).Methods("POST")
	r.HandleFunc("/dataDiscovery/{nsId}", appHandler.RunDataDiscovery).Methods("GET")
	r.HandleFunc("/addDataItem", appHandler.AddDataItem).Methods("POST")
	r.HandleFunc("/deleteApp", appHandler.DeleteApp).Methods("DELETE")
	r.HandleFunc("/softlink", appHandler.CreateSoftlink).Methods("POST")
	r.HandleFunc("/changeState", appHandler.ChangeDSIState).Methods("PUT")
	r.HandleFunc("/changePermissions", appHandler.ChangePermissions).Methods("PUT")
	r.HandleFunc("/putScheme", appHandler.PutScheme).Methods("PUT")

	srv := &http.Server{
		Handler: r,
		Addr:    ":8001",
	}
	log.Fatal(srv.ListenAndServe())

}

func Conn() *nats.Conn {
	conn, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
