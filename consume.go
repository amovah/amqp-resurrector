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
	deliveryCh, err := c.Channel.Consume(
		consume.Queue,
		consume.Consumer,
		consume.AutoAck,
		consume.Exclusive,
		consume.NoLocal,
		consume.NoWait,
		consume.Arg,
	)
	if err != nil {
		return nil, err
	}

	durableCh := make(chan amqp.Delivery)
	c.consumes[&consume] = durableCh

	go func() {
		for msg := range deliveryCh {
			durableCh <- msg
		}

		if c.isClosed {
			close(durableCh)
		}
	}()

	return durableCh, nil
}
