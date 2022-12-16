package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	amqpresurrector "github.com/amovah/amqp-resurrector"
	"github.com/rabbitmq/amqp091-go"
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

	count := 1
	for {
		err := ch.PublishWithContext(
			context.TODO(),
			"",
			"hello",
			false,
			false,
			amqp091.Publishing{
				DeliveryMode: amqp091.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(strconv.Itoa(count)),
			},
		)
		if err != nil {
			fmt.Println("error in publishing")
		} else {
			fmt.Println("published", count)
			count++
		}

		time.Sleep(1 * time.Second)
	}
}
