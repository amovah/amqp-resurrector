package amqpresurrector

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Channel struct {
	*amqp.Channel
	consumes    map[*Consume]chan amqp.Delivery
	queues      []Queue
	queueBinds  []QueueBind
	exchanges   []Exchange
	qos         *ChannelQoS
	notifyClose func()
	isTx        bool
	isClosed    bool
}

type ChannelQoS struct {
	PrefetchCount int
	PrefetchSize  int
	Global        bool
}

func (c *Channel) QoS(qos ChannelQoS) error {
	err := c.Channel.Qos(qos.PrefetchCount, qos.PrefetchSize, qos.Global)
	if err != nil {
		return err
	}

	c.qos = &qos
	return nil
}

func (c *Channel) reconnect() error {
	if c.qos != nil {
		if err := c.QoS(*c.qos); err != nil {
			return err
		}
	}

	for _, queue := range c.queues {
		_, err := c.Channel.QueueDeclare(
			queue.Name,
			queue.Durable,
			queue.AutoDelete,
			queue.Exclusive,
			queue.NoWait,
			queue.Args,
		)
		if err != nil {
			return err
		}
	}

	for _, exchange := range c.exchanges {
		err := c.Channel.ExchangeDeclare(
			exchange.Name,
			exchange.Kind,
			exchange.Durable,
			exchange.AutoDelete,
			exchange.Internal,
			exchange.NoWait,
			exchange.Arg,
		)
		if err != nil {
			return err
		}
	}

	for _, queueBind := range c.queueBinds {
		err := c.Channel.QueueBind(
			queueBind.Name,
			queueBind.Key,
			queueBind.Exchange,
			queueBind.NoWait,
			queueBind.Arg,
		)
		if err != nil {
			return err
		}
	}

	deliveryMap := make(map[*Consume]<-chan amqp.Delivery)
	for consume := range c.consumes {
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
			return err
		}

		deliveryMap[consume] = deliveryCh
	}

	for consume, durableCh := range c.consumes {
		go func(consume *Consume, durableCh chan amqp.Delivery) {
			for msg := range deliveryMap[consume] {
				durableCh <- msg
			}

			if c.isClosed {
				close(durableCh)
			}
		}(consume, durableCh)
	}

	if c.isTx {
		if err := c.Channel.Tx(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Channel) Tx() error {
	if err := c.Channel.Tx(); err != nil {
		return err
	}

	c.isTx = true
	return nil
}

func (c *Channel) Close() error {
	c.isClosed = true
	if err := c.Channel.Close(); err != nil {
		c.isClosed = false
		return err
	}
	c.notifyClose()
	return nil
}

func (c *Channel) cleanup() {
	for _, durableCh := range c.consumes {
		close(durableCh)
	}
}
