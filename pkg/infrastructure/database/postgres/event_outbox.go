package postgres

import (
	"context"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/postgres/model"
)

func (d *DBSession) SaveEvent(ctx context.Context, event event.Event) error {
	_, err := d.session.Context(ctx).InsertOne(model.DomainEventFromDomain(event))
	if d.isDuplicateError(err) {
		return nil
	}
	return err
}
