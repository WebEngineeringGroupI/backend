package postgres

import (
	"context"

	"github.com/WebEngineeringGroupI/backend/pkg/domain"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/postgres/model"
)

func (d *DBSession) SaveEvent(ctx context.Context, event domain.Event) error {
	_, err := d.session.Context(ctx).InsertOne(model.DomainEventFromDomain(event))
	return err
}
