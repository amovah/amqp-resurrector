package amqpresurrector

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	*amqp.Connection
	channels       []*Channel
	closedManually bool
}

func reconnectChannels(conn *Connection) error {
	openedChannels := make([]*amqp.Channel, len(conn.channels), len(conn.channels))
	for i, ch := range conn.channels {
		createdChannel, err := conn.Connection.Channel()
		if err != nil {
			return err
		}

		openedChannels[i] = createdChannel

		tempChan := Channel{
			Channel:   createdChannel,
			consumes:  ch.consumes,
			queues:    ch.queues,
			exchanges: ch.exchanges,
			qos:       ch.qos,
		}

		if err := tempChan.reconnect(); err != nil {
			createdChannel.Close()
			return err
		}
	}

	for i, ch := range conn.channels {
		ch.Channel = openedChannels[i]
	}

	return nil
}

func Dial(url string, reconnectTickRate time.Duration) (*Connection, error) {
	amqpConn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	conn := &Connection{
		Connection: amqpConn,
		channels:   make([]*Channel, 0, 10),
	}

	go func() {
		for {
			<-amqpConn.NotifyClose(make(chan *amqp.Error))

			if conn.closedManually {
				break
			}

			innerConn, err := amqp.Dial(url)
			if err != nil {
				time.Sleep(reconnectTickRate)
				continue
			}

			conn.Connection = innerConn
			amqpConn = innerConn

			channelErr := reconnectChannels(conn)
			for channelErr != nil {
				time.Sleep(reconnectTickRate)
				channelErr = reconnectChannels(conn)
			}
		}

		for _, ch := range conn.channels {
			ch.cleanup()
		}
	}()

	return conn, nil
}

func (c *Connection) Channel() (*Channel, error) {
	amqpCh, err := c.Connection.Channel()
	if err != nil {
		return nil, err
	}

	ch := &Channel{
		Channel:   amqpCh,
		consumes:  make(map[*Consume]chan amqp.Delivery),
		queues:    make([]Queue, 0, 10),
		exchanges: make([]Exchange, 0, 10),
	}

	c.channels = append(c.channels, ch)

	return ch, nil
}

func (c *Connection) Close() error {
	c.closedManually = true
	return c.Connection.Close()
}
