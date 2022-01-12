package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/streadway/amqp"

	"github.com/WebEngineeringGroupI/backend/internal/app"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/event/serializer/json"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url/validator"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/broker/rabbitmq"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/clock"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/validator/pipeline"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/validator/reachable"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/validator/safebrowsing"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/validator/schema"
)

type Factory struct {
}

func (f *Factory) NewValidator(ctx context.Context) *validator.Service {
	f.defineRabbitMQRouting()

	return validator.NewService(
		f.brokerReceiver(ctx),
		f.brokerSender(ctx),
		f.urlValidator(),
		json.NewSerializer(&url.ShortURLCreated{}, &url.LoadBalancedURLCreated{}),
		clock.NewFromSystem())
}

func (f *Factory) brokerReceiver(ctx context.Context) *rabbitmq.ReceiverClient {
	receiver, err := rabbitmq.NewReceiverClient(ctx, app.RabbitMQConnectionString(), "urlshortener_to_validator")
	if err != nil {
		log.Fatalf("unable to create receiver client: %s", err)
	}
	return receiver
}

func (f *Factory) brokerSender(ctx context.Context) *rabbitmq.SenderClient {
	sender, err := rabbitmq.NewSenderClient(ctx, app.RabbitMQConnectionString(), "validator", "validator")
	if err != nil {
		log.Fatalf("unable to create rabbitmq sender: %s", err)
	}
	return sender
}

func (f *Factory) defineRabbitMQRouting() {
	conn, err := amqp.Dial(app.RabbitMQConnectionString())
	if err != nil {
		log.Fatalf("unable to connect to rabbitmq to define routings: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("unable to create rabbitmq channel to define routings: %s", err)
	}
	defer ch.Close()

	f.defineRabbitMQRoutingToReceive(ch)
	f.defineRabbitMQRoutingToSend(ch)
}

func (f *Factory) defineRabbitMQRoutingToReceive(ch *amqp.Channel) {
	if err := ch.ExchangeDeclare("urlshortener", "topic", true, false, false, false, nil); err != nil {
		log.Fatalf("unable to declare exchange to receive events from: %s", err)
	}

	queueName := "urlshortener_to_validator"
	if _, err := ch.QueueDeclare(queueName, true, false, false, false, nil); err != nil {
		log.Fatalf("unable to declare rabbitmq queue to receive events from: %s", err)
	}

	if err := ch.QueueBind(queueName, "#", "urlshortener", false, nil); err != nil {
		log.Fatalf("unable to bind queue 'urlshortener' to exchange 'urlshortener_validator': %s", err)
	}
}

func (f *Factory) defineRabbitMQRoutingToSend(ch *amqp.Channel) {
	if err := ch.ExchangeDeclare("validator", "topic", true, false, false, false, nil); err != nil {
		log.Fatalf("unable to declare exchange to receive events from: %s", err)
	}
}

func (f *Factory) urlValidator() url.Validator {
	safebrowsingValidator, err := safebrowsing.NewValidator(app.SafeBrowsingAPIKey())
	if err != nil {
		log.Fatalf("unable to create safebrowsing validator: %s", err)
	}

	return pipeline.NewValidator(
		schema.NewValidator("https", "http"),
		reachable.NewValidator(http.DefaultClient, 5*time.Second),
		safebrowsingValidator,
	)
}
