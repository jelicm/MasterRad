package main

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	fmt.Println("neka funkcajo")

	conn := Conn()
	defer conn.Close()

	//subject := "app1+novaapp0/Root/folder0"

	/*err := conn.Publish(subject, []byte("aaaajo"))
	if err != nil {
		log.Fatal(err)
	}*/
	_, err := conn.Subscribe("mojtopic", func(message *nats.Msg) {
		fmt.Printf("RECEIVED MESSAGE event %d: %s\n", 0, string(message.Data))
	})
	if err != nil {
		log.Fatal(err)
	}
	select {}
}
func Conn() *nats.Conn {
	conn, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
