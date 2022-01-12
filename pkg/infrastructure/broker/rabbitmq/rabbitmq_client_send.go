package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/event/redirector"
)

type SenderClient struct {
	sendExchangeName string
	sendRoutingKey   string
	sendConnection   *amqp.Connection
	connectionMutex  sync.RWMutex
	connectionString string
}

var _ redirector.ExternalBrokerSender = (*SenderClient)(nil)

func NewSenderClient(ctx context.Context, connectionString string, sendExchangeName string, sendRoutingKey string) (*SenderClient, error) {
	connection, err := amqp.Dial(connectionString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to rabbitmq: %w", err)
	}

	senderClient := &SenderClient{
		sendExchangeName: sendExchangeName,
		sendRoutingKey:   sendRoutingKey,
		sendConnection:   connection,
		connectionString: connectionString,
	}
	go senderClient.reconnectOnError(ctx)
	return senderClient, nil
}

func (c *SenderClient) SendEvents(ctx context.Context, eventData ...[]byte) error {
	c.connectionMutex.RLock()
	defer c.connectionMutex.RUnlock()

	ch, err := c.sendConnection.Channel()
	if err != nil {
		return fmt.Errorf("unable to create channel to send events: %w", err)
	}
	defer ch.Close()

	if err := ch.Tx(); err != nil {
		return fmt.Errorf("unable to create a transactional channel: %w", err)
	}

	for _, data := range eventData {
		err := ch.Publish(
			c.sendExchangeName,
			c.sendRoutingKey,
			false,
			false,
			amqp.Publishing{Body: data})
		if err != nil {
			_ = ch.TxRollback()
			return fmt.Errorf("unable to publish message to channel: %w", err)
		}
	}

	_ = ch.TxCommit()
	return nil
}

func (c *SenderClient) reconnectOnError(ctx context.Context) {
	defer c.sendConnection.Close()

	for {
		notifyClose := c.sendConnection.NotifyClose(make(chan *amqp.Error))
		select {
		case <-notifyClose:
		case <-ctx.Done():
			return
		}

		backoffIdx := 1

	reconnectLoop:
		for {
			select {
			case <-ctx.Done():
				return
			default:

				connection, err := amqp.Dial(c.connectionString)
				if err != nil {
					log.Printf("unable to reconnect to rabbitmq: %s", err)
					time.Sleep(time.Duration(backoffIdx) * time.Second)
					if backoffIdx < 15 {
						backoffIdx++
					}
					continue
				}
				c.connectionMutex.Lock()
				log.Println("reconnected to rabbitmq")
				c.sendConnection = connection
				c.connectionMutex.Unlock()
				break reconnectLoop
			}
		}
	}
}
