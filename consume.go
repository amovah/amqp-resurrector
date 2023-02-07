package amqpresurrector

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consume struct {
	Queue     string
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Arg       amqp.Table
}

func (c *Channel) Consume(consume Consume) (<-chan amqp.Delivery, error) {
	return c.Channel.Consume(
		consume.Queue,
		consume.Consumer,
		consume.AutoAck,
		consume.Exclusive,
		consume.NoLocal,
		consume.NoWait,
		consume.Arg,
	)
}
