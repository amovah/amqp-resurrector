package amqpresurrector

import amqp "github.com/rabbitmq/amqp091-go"

type Exchange struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Arg        amqp.Table
}

func (c *Channel) ExchangeDeclare(exchange Exchange) error {
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

	c.exchanges = append(c.exchanges, exchange)

	return nil
}
