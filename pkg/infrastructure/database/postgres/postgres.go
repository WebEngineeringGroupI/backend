package postgres

import (
	`errors`
	`fmt`

	`github.com/lib/pq`
	_ "github.com/lib/pq"
	"xorm.io/xorm"

	`github.com/WebEngineeringGroupI/backend/pkg/domain/url`
	`github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/model`
)

type ConnectionDetails struct {
	User     string
	Pass     string
	Host     string
	Port     int
	Database string
	SSLMode  string
}

type DB struct {
	engine *xorm.Engine
}

var (
	errDuplicateConstraintViolation = "23505"
)

func (d *DB) Save(url *url.ShortURL) error {
	shortURL := model.ShortURLFromDomain(url)
	_, err := d.engine.Insert(&shortURL)

	var pqError pq.Error
	if errors.Is(err, &pqError) {
		if pqError.Code == pq.ErrorCode(errDuplicateConstraintViolation) {
			return nil
		}
		return fmt.Errorf("unable to save short URL: %w", err)
	}
	return nil
}

func (d *DB) FindByHash(hash string) (*url.ShortURL, error) {
	shortUrl := model.Shorturl{Hash: hash}
	exists, err := d.engine.Get(&shortUrl)
	if !exists {
		return nil, url.ErrShortURLNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("unknown error retrieving short url: %w", err)
	}

	return model.ShortURLToDomain(shortUrl), nil
}

func NewDB(connectionDetails ConnectionDetails) (*DB, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		connectionDetails.User,
		connectionDetails.Pass,
		connectionDetails.Host,
		connectionDetails.Port,
		connectionDetails.Database,
		connectionDetails.SSLMode)

	engine, err := xorm.NewEngine("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection to database: %w", err)
	}

	return &DB{
		engine: engine,
	}, nil

}
