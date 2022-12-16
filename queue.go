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
	q, err := c.Channel.QueueDeclare(
		queue.Name,
		queue.Durable,
		queue.AutoDelete,
		queue.Exclusive,
		queue.NoWait,
		queue.Args,
	)
	if err != nil {
		return amqp.Queue{}, err
	}

	c.queues = append(c.queues, queue)
	return q, nil
}

type QueueBind struct {
	Name     string
	Key      string
	Exchange string
	NoWait   bool
	Arg      amqp.Table
}

func (c *Channel) QueueBind(qb QueueBind) error {
	err := c.Channel.QueueBind(
		qb.Name,
		qb.Key,
		qb.Exchange,
		qb.NoWait,
		qb.Arg,
	)
	if err != nil {
		return err
	}

	c.queueBinds = append(c.queueBinds, qb)
	return nil
}

func (c *Channel) QueueUnbind(qb QueueBind) error {
	err := c.Channel.QueueUnbind(
		qb.Name,
		qb.Key,
		qb.Exchange,
		qb.Arg,
	)
	if err != nil {
		return err
	}

	for i, queueBind := range c.queueBinds {
		if queueBind.Name == qb.Name && queueBind.Key == qb.Key && queueBind.Exchange == qb.Exchange {
			c.queueBinds = append(c.queueBinds[:i], c.queueBinds[i+1:]...)
		}
	}

	return nil
}
