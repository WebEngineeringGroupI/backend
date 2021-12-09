package postgres_test

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/postgres"
)

var _ = Describe("Postgres", func() {
	var (
		db  *postgres.DB
		ctx context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()

		var err error
		db, err = postgres.NewDB(connectionDetails())
		Expect(err).ToNot(HaveOccurred())
	})

	It("saves the short URL in the database and retrieves it again", func() {
		hash := randomHash()
		_, err := db.Transactional(func(session *postgres.DBSession) (interface{}, error) {
			_ = session.SaveShortURL(ctx, aShortURLWithHash(hash))

			retrievedShortURL, err := session.FindShortURLByHash(ctx, hash)

			Expect(err).ToNot(HaveOccurred())
			Expect(retrievedShortURL.Hash).To(Equal(hash))
			Expect(retrievedShortURL.OriginalURL).To(Equal(url.OriginalURL{URL: "https://google.com", IsValid: true}))
			return nil, nil
		})
		Expect(err).ToNot(HaveOccurred())
	})

	It("saves the short URL in the database and rolls it back again on error", func() {
		hash := randomHash()
		_, err := db.Transactional(func(session *postgres.DBSession) (interface{}, error) {
			err := session.SaveShortURL(ctx, aShortURLWithHash(hash))
			Expect(err).ToNot(HaveOccurred())

			retrievedShortURL, err := session.FindShortURLByHash(ctx, hash)
			Expect(err).ToNot(HaveOccurred())
			Expect(retrievedShortURL.Hash).To(Equal(hash))
			Expect(retrievedShortURL.OriginalURL).To(Equal(url.OriginalURL{URL: "https://google.com", IsValid: true}))

			return nil, errors.New("an error")
		})
		Expect(err).To(MatchError("an error"))

		retrievedShortURL, err := db.Session().FindShortURLByHash(ctx, hash)
		Expect(err).To(MatchError(url.ErrShortURLNotFound))
		Expect(retrievedShortURL).To(BeNil())
	})
})

func randomHash() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Int())[0:7]
}
