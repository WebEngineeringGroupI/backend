package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/event/redirector"
)

type ReceiverClient struct {
	connectionMutex  sync.RWMutex
	recvConnection   *amqp.Connection
	recvQueueName    string
	connectionString string
}

var _ redirector.ExternalBrokerReceiver = (*ReceiverClient)(nil)

func NewReceiverClient(ctx context.Context, connectionString string, recvQueueName string) (*ReceiverClient, error) {
	connection, err := amqp.Dial(connectionString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to rabbitmq: %w", err)
	}

	receiverClient := &ReceiverClient{
		recvConnection:   connection,
		connectionString: connectionString,
		recvQueueName:    recvQueueName,
	}
	go receiverClient.reconnectOnError(ctx)
	return receiverClient, nil
}

func (c *ReceiverClient) ReceiveEvents(ctx context.Context) (<-chan []byte, error) {
	channel, err := c.recvConnection.Channel()
	if err != nil {
		return nil, fmt.Errorf("unable to create channel to receive events: %w", err)
	}

	consumerID := uuid.New().String()
	consumeCh, err := channel.Consume(
		c.recvQueueName,
		consumerID, // consumer
		false,      // autoAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // args
	)
	if err != nil {
		return nil, fmt.Errorf("unable to start consuming events from channel: %w", err)
	}

	ch := make(chan []byte)

	go func() {
		defer close(ch)
		defer channel.Close()

		for {
			select {
			case <-ctx.Done():
				_ = channel.Cancel(consumerID, false)
				for msg := range consumeCh {
					_ = msg.Nack(false, true)
				}
				return
			case message := <-consumeCh:
				ch <- message.Body
				if err := message.Ack(false); err != nil {
					log.Println("unable to ACK message received")
				}
			}
		}
	}()

	return ch, nil
}

func (c *ReceiverClient) reconnectOnError(ctx context.Context) {
	defer c.recvConnection.Close()

	for {
		notifyClose := c.recvConnection.NotifyClose(make(chan *amqp.Error))
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
				c.recvConnection = connection
				c.connectionMutex.Unlock()
				log.Println("reconnected to rabbitmq")
				break reconnectLoop
			}
		}
	}
}
