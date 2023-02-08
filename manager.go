package amqpresurrector

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type MqConfig struct {
	Host          string
	Port          int
	Username      string
	Password      string
	Vhost         string
	ReconnectRate int
}

func newAmqpConnection(config MqConfig) (*amqp091.Connection, error) {
	return amqp091.Dial(
		fmt.Sprintf(
			"amqp://%s:%s@%s:%d/%s",
			config.Username,
			config.Password,
			config.Host,
			config.Port,
			config.Vhost,
		),
	)
}

func reconnect(ctx context.Context, logger *logrus.Logger, config MqConfig) *amqp091.Connection {
	log := logger
	if log == nil {
		log = logrus.New()
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			amqpConn, err := newAmqpConnection(config)
			if err != nil {
				log.WithError(err).Warn("failed to establish connection to rabbit")
				time.Sleep(time.Duration(config.ReconnectRate) * time.Millisecond)
				break
			}

			log.Info("connection to rabbit has been restored")

			return amqpConn
		}
	}
}

type Manager struct {
	Conn           *amqp091.Connection
	notifier       chan bool
	done           chan bool
	ctx            context.Context
	once           sync.Once
	mutex          sync.Mutex
	reconnectCalls uint64
	connectionId   uuid.UUID
}

func (m *Manager) Reconnect(connectionId uuid.UUID) {
	m.mutex.Lock()

	if m.connectionId != connectionId {
		m.mutex.Unlock()
		return
	}

	m.reconnectCalls += 1

	m.once.Do(func() {
		select {
		case <-m.ctx.Done():
			return
		case m.notifier <- true:
		}
	})

	m.mutex.Unlock()

	<-m.done
}

func (m *Manager) ConnectionId() uuid.UUID {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.connectionId
}

func NewManager(ctx context.Context, config MqConfig, logger *logrus.Logger) (*Manager, error) {
	conn, err := newAmqpConnection(config)
	if err != nil {
		return nil, err
	}

	manager := &Manager{
		Conn:         conn,
		notifier:     make(chan bool),
		ctx:          ctx,
		once:         sync.Once{},
		done:         make(chan bool),
		connectionId: uuid.New(),
	}

	go func() {
		for range manager.notifier {
			manager.Conn.Close()

			conn := reconnect(ctx, logger, config)
			manager.Conn = conn

			manager.mutex.Lock()

			for i := 0; i < int(manager.reconnectCalls); i++ {
				manager.done <- true
			}
			manager.reconnectCalls = 0
			manager.once = sync.Once{}
			manager.connectionId = uuid.New()

			manager.mutex.Unlock()
		}
	}()

	return manager, nil
}
