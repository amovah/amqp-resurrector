package amqpresurrector

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Channel struct {
	*amqp.Channel
}

type ChannelQoS struct {
	PrefetchCount int
	PrefetchSize  int
	Global        bool
}

func (c *Channel) QoS(qos ChannelQoS) error {
	return c.Channel.Qos(qos.PrefetchCount, qos.PrefetchSize, qos.Global)
}
