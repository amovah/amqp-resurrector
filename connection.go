package amqpresurrector

import (
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	*amqp.Connection
	closedManually bool
	mutex          sync.Mutex
}

func Dial(url string, reconnectTickRate time.Duration) (*Connection, error) {
	amqpConn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	conn := &Connection{
		Connection: amqpConn,
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
		Channel: amqpCh,
	}

	return ch, nil
}

func (c *Connection) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.closedManually = true
	return c.Connection.Close()
}
