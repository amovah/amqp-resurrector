package amqpresurrector

import amqp "github.com/rabbitmq/amqp091-go"

type Queue struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

func (c *Channel) QueueDeclare(queue Queue) (amqp.Queue, error) {
	return c.Channel.QueueDeclare(
		queue.Name,
		queue.Durable,
		queue.AutoDelete,
		queue.Exclusive,
		queue.NoWait,
		queue.Args,
	)
}

type QueueBind struct {
	Name     string
	Key      string
	Exchange string
	NoWait   bool
	Arg      amqp.Table
}

func (c *Channel) QueueBind(qb QueueBind) error {
	return c.Channel.QueueBind(
		qb.Name,
		qb.Key,
		qb.Exchange,
		qb.NoWait,
		qb.Arg,
	)
}

func (c *Channel) QueueUnbind(qb QueueBind) error {
	return c.Channel.QueueUnbind(
		qb.Name,
		qb.Key,
		qb.Exchange,
		qb.Arg,
	)
}
