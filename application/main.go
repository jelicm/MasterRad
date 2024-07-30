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

	subject := "app1+app23/Root/file.txt"

	err := conn.Publish(subject, []byte("aaaajo"))
	if err != nil {
		log.Fatal(err)
	}
}
func Conn() *nats.Conn {
	conn, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
	return conn
}
