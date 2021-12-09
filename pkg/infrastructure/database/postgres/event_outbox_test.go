package postgres_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/postgres"
)

var _ = Describe("Infrastructure / Database / Postgres / EventOutbox", func() {
	var (
		repository *postgres.DB
		session    *postgres.DBSession
		ctx        context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		var err error
		repository, err = postgres.NewDB(connectionDetails())
		Expect(err).ToNot(HaveOccurred())
		session = repository.Session()
	})

	AfterEach(func() {
		session.Close()
	})

	It("saves the event in the database", func() {
		hash := randomHash()
		err := session.SaveEvent(ctx, url.NewShortURLCreated(aShortURLWithHash(hash)))
		Expect(err).ToNot(HaveOccurred())
		//
		//retrievedShortURL, err := session.FindShortURLByHash(ctx, hash)
		//
		//Expect(err).ToNot(HaveOccurred())
		//Expect(retrievedShortURL.Hash).To(Equal(hash))
		//Expect(retrievedShortURL.OriginalURL).To(Equal(url.OriginalURL{URL: "https://google.com", IsValid: true}))
	})
})
