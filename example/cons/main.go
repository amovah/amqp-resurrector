package main

import (
	"fmt"
	"log"
	"time"

	amqpresurrector "github.com/amovah/amqp-resurrector"
)

func main() {
	conn, err := amqpresurrector.Dial("amqp://guest:guest@localhost:5672", 2*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	q, err := ch.QueueDeclare(amqpresurrector.Queue{
		Name:    "hello",
		Durable: true,
	})
	fmt.Println(q.Name, err)

	time.AfterFunc(5*time.Second, func() {
		conn.Close()
		fmt.Println("conn closed")
	})

	msgs, err := ch.Consume(amqpresurrector.Consume{
		Queue: "hello",
	})
	if err != nil {
		log.Fatal(err)
	}

	for msg := range msgs {
		fmt.Println(string(msg.Body))
		msg.Ack(false)
	}
}
