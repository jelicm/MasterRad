package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type EventDTO struct {
	EventTopic string `json:"event"`
	AddEvent   bool   `json:"add"`
}

type EventHandler struct {
	EventTopics []string
	//Conn        *nats.Conn
}

func (handler *EventHandler) helloPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Post iz triggera!\n")
	var eventDTO EventDTO
	err := json.NewDecoder(r.Body).Decode(&eventDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rez := eventDTO.EventTopic
	if eventDTO.AddEvent {
		handler.EventTopics = append(handler.EventTopics, rez)
	} else {
		fmt.Println("remove")
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(rez))
}

/*func (handler *EventHandler) trigger() {
	time.Sleep(2 * time.Second)
	for _, topic := range handler.EventTopics {
		err := handler.Conn.Publish(topic, []byte("aaaajo"))
		if err != nil {
			log.Fatal(err)
		}
	}
	err := handler.Conn.Publish("mojtopic", []byte("aaaajo"))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("poslata porukica")
}*/

func main() {
	//handler := EventHandler{EventTopics: []string{}, Conn: Conn()}
	handler := EventHandler{EventTopics: []string{}}
	//handler.trigger()
	r := mux.NewRouter()
	r.HandleFunc("/", handler.helloPost).Methods("POST")
	fmt.Println("Listening on :8000...")
	srv := &http.Server{
		Handler: r,
		Addr:    ":8000",
	}
	log.Fatal(srv.ListenAndServe())
}

/*func Conn() *nats.Conn {
	conn, err := nats.Connect("trigger-nats.internal:4222")
	if err != nil {
		log.Fatal(err)
	}
	return conn
}*/
