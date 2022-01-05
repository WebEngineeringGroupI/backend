package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"xorm.io/xorm"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/click"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/postgres/model"
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

type DBSession struct {
	session *xorm.Session
}

func (d *DBSession) Append(ctx context.Context, identity string, events ...event.Event) error {
	panic("implement me")
}

func (d *DBSession) Load(ctx context.Context, identity string) (*event.Stream, error) {
	panic("implement me")
}

func (d *DBSession) FindLoadBalancedURLByHash(ctx context.Context, hash string) (*url.LoadBalancedURL, error) {
	var result model.LoadBalancedUrlList
	err := d.session.Context(ctx).Find(&result, &model.LoadBalancedUrl{Hash: hash})
	if len(result) == 0 {
		return nil, url.ErrValidURLNotFound // FIXME(fede): Should we use another kind of error here?
	}
	if err != nil {
		return nil, fmt.Errorf("unknown error retrieving short url: %w", err)
	}

	return model.LoadBalancedURLToDomain(result), nil
}

func (d *DBSession) SaveLoadBalancedURL(ctx context.Context, aURL *url.LoadBalancedURL) error {
	dbURL := model.LoadBalancedURLFromDomain(aURL)
	_, err := d.session.Context(ctx).Insert(&dbURL)
	if d.isDuplicateError(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("unable to save load-balanced URL: %w", err)
	}
	return nil
}

var (
	errDuplicateConstraintViolation pq.ErrorCode = "23505"
)

func (d *DBSession) SaveShortURL(ctx context.Context, url *url.ShortURL) error {
	shortURL := model.ShortURLFromDomain(url)
	_, err := d.session.Context(ctx).Insert(&shortURL)
	if d.isDuplicateError(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("unable to save short URL: %w", err)
	}
	return nil
}

func (d *DBSession) FindShortURLByHash(ctx context.Context, hash string) (*url.ShortURL, error) {
	shortURL := model.Shorturl{Hash: hash}
	exists, err := d.session.Context(ctx).Get(&shortURL)
	if !exists {
		return nil, url.ErrShortURLNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("unknown error retrieving short url: %w", err)
	}

	return model.ShortURLToDomain(shortURL), nil
}

func (d *DBSession) SaveClick(ctx context.Context, click *click.Details) error {
	clickModel := model.ClickDetailsFromDomain(click)
	_, err := d.session.Context(ctx).Insert(&clickModel)
	if err != nil {
		return fmt.Errorf("unknow error saving click: %w", err)
	}
	return nil
}

func (d *DBSession) FindClicksByHash(ctx context.Context, hash string) ([]*click.Details, error) {
	var clicksModel []*model.Clickdetails
	err := d.session.Context(ctx).Find(&clicksModel, model.Clickdetails{Hash: hash})
	if err != nil {
		return nil, fmt.Errorf("unknow error finding clicks by hash: %w", err)
	}
	var clicks []*click.Details
	for _, clickModel := range clicksModel {
		clicks = append(clicks, model.ClickDetailsToDomain(clickModel))
	}
	return clicks, nil
}

func (d *DBSession) Close() error {
	return d.session.Close()
}

func (d *DBSession) isDuplicateError(err error) bool {
	var pqError *pq.Error
	if errors.As(err, &pqError) {
		if pqError.Code == errDuplicateConstraintViolation {
			return true
		}
	}
	return false
}

func newDBSession(session *xorm.Session) *DBSession {
	return &DBSession{session: session}
}

type DB struct {
	engine *xorm.Engine
}

func NewDB(connectionDetails ConnectionDetails) (*DB, error) {
	engine, err := xorm.NewEngine("postgres", connectionDetails.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("unable to create connection to database: %w", err)
	}

	return &DB{
		engine: engine,
	}, nil
}

func (t *DB) Transactional(f func(*DBSession) (interface{}, error)) (interface{}, error) {
	session := t.engine.NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return nil, err
	}

	result, err := f(newDBSession(session))
	if err != nil {
		return result, err
	}

	if err := session.Commit(); err != nil {
		return result, err
	}

	return result, nil
}

func (t *DB) Session() *DBSession {
	return newDBSession(t.engine.NewSession())
}
