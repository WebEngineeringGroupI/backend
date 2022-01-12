package main

import (
	"context"
	"log"
	"time"

	"github.com/streadway/amqp"

	"github.com/WebEngineeringGroupI/backend/internal/app"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/event/redirector"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/event/serializer/json"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/broker/rabbitmq"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/postgres"
)

type Factory struct {
	rabbitMQSingleton redirector.ExternalBrokerSender
}

func (f *Factory) NewRedirector(ctx context.Context) *redirector.Redirector {
	return redirector.NewRedirector(f.newPostgresqlDB(), f.newRabbitMQ(ctx), 5*time.Second)
}

func (f *Factory) connectionDetails() *postgres.ConnectionDetails {
	return app.PostgresConnectionDetails()
}

func (f *Factory) newPostgresqlDB() redirector.OutboxSource {
	db, err := postgres.NewDB(f.connectionDetails(), json.NewSerializer())
	check(err, "unable to create database connection: %s\n", err)
	return db
}

func (f *Factory) newRabbitMQ(ctx context.Context) redirector.ExternalBrokerSender {
	if f.rabbitMQSingleton != nil {
		return f.rabbitMQSingleton
	}

	connection := f.newRabbitMQConnection()
	ch, err := connection.Channel()
	check(err, "unable to open channel to RabbitMQ: %s", err)
	defer ch.Close()

	declareExchangesToSendTo(ch)

	f.rabbitMQSingleton, err = rabbitmq.NewSenderClient(ctx, f.rabbitMQConnectionString(), "urlshortener", "event")
	check(err, "unable to create rabbitmq sender: %s", err)

	return f.rabbitMQSingleton
}

func (f *Factory) newRabbitMQConnection() *amqp.Connection {
	sendConnection, err := amqp.Dial(f.rabbitMQConnectionString())
	check(err, "unable to establish sending connection to RabbitMQ: %s", err)
	return sendConnection
}

//
//// FIXME(fede): Move these to the validator receiver
//func bindQueuesToListenExchanges(ch *amqp.Channel) {
//	err := ch.QueueBind("urlshortener_validator", "#", "urlshortener_validator", false, nil)
//	check(err, "unable to bind queue to exchange to receive events from: %s", err)
//}
//
//func declareQueuesToListenFrom(ch *amqp.Channel) {
//	_, err := ch.QueueDeclare("urlshortener_validator", true, false, false, false, nil)
//	check(err, "unable to declare queue to receive messages from: %s", err)
//}
//
//func declareExchangesToReceiveFrom(ch *amqp.Channel) {
//	err := ch.ExchangeDeclare("urlshortener_validator", "topic", true, false, false, false, nil)
//	check(err, "unable to declare exchange to receive messages from: %s", err)
//}

func declareExchangesToSendTo(ch *amqp.Channel) {
	err := ch.ExchangeDeclare("urlshortener", "topic", true, false, false, false, nil)
	check(err, "unable to declare exchange to send messages: %s", err)
}

func (f *Factory) rabbitMQConnectionString() string {
	return app.RabbitMQConnectionString()
}

func check(err error, msg string, args ...interface{}) {
	if err != nil {
		log.Fatalf(msg, args...)
	}
}
