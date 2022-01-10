package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"xorm.io/xorm"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/postgres/serializer"
)

type ConnectionDetails struct {
	User     string
	Pass     string
	Host     string
	Port     int
	Database string
	SSLMode  string
}

func (d *ConnectionDetails) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User,
		d.Pass,
		d.Host,
		d.Port,
		d.Database,
		d.SSLMode)
}

type DB struct {
	engine     *xorm.Engine
	serializer serializer.Serializer
}

type DomainEvent struct {
	ID      string `xorm:"'id'"`
	Version int    `xorm:"'version'"`
	Payload []byte `xorm:"'payload'"`
}

type DomainEventOutbox struct {
	ID      int    `xorm:"'id' autoincr"`
	Payload []byte `xorm:"'payload'"`
}

func (d *DB) Append(ctx context.Context, identity string, events ...event.Event) error {
	serializedEvents := make([]interface{}, 0, len(events))
	outboxEvents := make([]interface{}, 0, len(events))

	for _, event := range events {
		payload, err := d.serializer.MarshalEvent(event)
		if err != nil {
			return fmt.Errorf("unable to save event in the database: %w", err)
		}
		serializedEvents = append(serializedEvents, DomainEvent{
			ID:      identity,
			Version: event.EventVersion(),
			Payload: payload,
		})
		outboxEvents = append(outboxEvents, DomainEventOutbox{
			Payload: payload,
		})
	}

	_, err := d.engine.Transaction(func(session *xorm.Session) (interface{}, error) {
		_, err := session.Context(ctx).Insert(serializedEvents...)
		if isDuplicateError(err) {
			return nil, fmt.Errorf("unable to insert event in database, check the version of the events: %w", err)
		}
		if err != nil {
			return nil, fmt.Errorf("unable to insert events in database: %w", err)
		}

		_, err = session.Context(ctx).Insert(outboxEvents...)
		if err != nil {
			return nil, fmt.Errorf("unable to insert outbox events: %w", err)
		}

		return nil, nil
	})

	return err
}

func (d *DB) Load(ctx context.Context, identity string) (*event.Stream, error) {
	resultInterface, err := d.engine.Transaction(func(session *xorm.Session) (interface{}, error) {
		var result []DomainEvent
		err := session.Context(ctx).Find(&result, &DomainEvent{ID: identity})
		if err != nil {
			return nil, fmt.Errorf("unknown error retrieving short url: %w", err)
		}
		if len(result) == 0 {
			return nil, event.ErrEntityNotFound
		}
		return result, nil
	})
	if err != nil {
		return nil, err
	}

	result, ok := resultInterface.([]DomainEvent)
	if !ok {
		return nil, fmt.Errorf("result from transaction is not slice of domain events")
	}

	events := make([]event.Event, 0, len(result))
	for _, domainEvent := range result {
		event, err := d.serializer.UnmarshalEvent(domainEvent.Payload)
		if err != nil {
			return nil, fmt.Errorf("error retrieving event from database: %w", err)
		}
		events = append(events, event)
	}

	return event.StreamFrom(events), nil
}

func isDuplicateError(err error) bool {
	var pqError *pq.Error
	if errors.As(err, &pqError) {
		if pqError.Code == ("23505") {
			return true
		}
	}
	return false
}

func NewDB(connectionDetails *ConnectionDetails, serializer serializer.Serializer) (*DB, error) {
	engine, err := xorm.NewEngine("postgres", connectionDetails.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("unable to create connection to database: %w", err)
	}

	return &DB{
		engine:     engine,
		serializer: serializer,
	}, nil
}
