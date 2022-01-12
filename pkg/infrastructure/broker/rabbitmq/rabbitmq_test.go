package rabbitmq_test

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/streadway/amqp"

	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/broker/rabbitmq"
)

var _ = Describe("Infrastructure / Broker / RabbitMQ", func() {
	var (
		sendExchangeName string
		sendRoutingKey   string
		recvQueueName    string
		senderClient     *rabbitmq.SenderClient
		receiverClient   *rabbitmq.ReceiverClient
		ctx              context.Context
		cancel           context.CancelFunc
	)

	BeforeEach(func() {
		sendExchangeName = string(randomPayload())
		recvQueueName = string(randomPayload())
		sendRoutingKey = string(randomPayload())
	})

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		var err error
		senderClient, err = rabbitmq.NewSenderClient(ctx, rabbitMQConnectionString(), sendExchangeName, sendRoutingKey)
		Expect(err).ToNot(HaveOccurred())

		receiverClient, err = rabbitmq.NewReceiverClient(ctx, rabbitMQConnectionString(), recvQueueName)
		Expect(err).ToNot(HaveOccurred())

	})

	BeforeEach(func() {
		sendConn, err := amqp.Dial(rabbitMQConnectionString())
		Expect(err).ToNot(HaveOccurred())
		defer sendConn.Close()

		ch, err := sendConn.Channel()
		Expect(err).ToNot(HaveOccurred())
		defer ch.Close()

		err = ch.ExchangeDeclare(sendExchangeName, "topic", false, true, false, false, nil)
		Expect(err).ToNot(HaveOccurred())

		_, err = ch.QueueDeclare(recvQueueName, false, true, false, false, nil)
		Expect(err).ToNot(HaveOccurred())

		err = ch.QueueBind(recvQueueName, "#", sendExchangeName, false, nil)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		cancel()

		sendConn, err := amqp.Dial(rabbitMQConnectionString())
		Expect(err).ToNot(HaveOccurred())
		defer sendConn.Close()

		ch, err := sendConn.Channel()
		Expect(err).ToNot(HaveOccurred())
		defer ch.Close()

		_, err = ch.QueueDelete(recvQueueName, true, true, false)
		Expect(err).ToNot(HaveOccurred())

		err = ch.ExchangeDelete(sendExchangeName, true, false)
		Expect(err).ToNot(HaveOccurred())
	})

	It("sends the events to the rabbitmq queue", func() {
		someEvent := randomPayload()

		err := senderClient.SendEvents(ctx, someEvent)
		Expect(err).ToNot(HaveOccurred())

		Expect(messageInQueue(recvQueueName).Body).To(Equal(someEvent))
	})

	It("receives the events from the rabbitmq queue", func() {
		someEvent := randomPayload()
		err := senderClient.SendEvents(ctx, someEvent)
		Expect(err).ToNot(HaveOccurred())

		ch, err := receiverClient.ReceiveEvents(ctx)
		Expect(err).ToNot(HaveOccurred())
		Eventually(ch).Should(Receive(Equal(someEvent)))
	})
})

func messageInQueue(queueName string) amqp.Delivery {
	conn, err := amqp.Dial(rabbitMQConnectionString())
	Expect(err).ToNot(HaveOccurred())
	defer conn.Close()

	channel, err := conn.Channel()
	ExpectWithOffset(1, err).ToNot(HaveOccurred())
	msg, ok, err := channel.Get(queueName, true /*autoAck*/)
	ExpectWithOffset(1, err).ToNot(HaveOccurred())
	ExpectWithOffset(1, ok).To(BeTrue())
	ExpectWithOffset(1, msg).ToNot(BeNil())

	return msg
}

func rabbitMQConnectionString() string {
	return "amqp://user:password@localhost:5672/"
}

func randomPayload() []byte {
	rand.Seed(time.Now().UnixNano())
	return []byte(strconv.Itoa(rand.Int()))
}
