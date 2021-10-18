package postgres_test

import (
	. "github.com/onsi/ginkgo"
	. `github.com/onsi/gomega`

	`github.com/WebEngineeringGroupI/backend/pkg/domain/url`
	`github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/postgres`
)

var _ = Describe("Postgres", func() {
	var (
		repository *postgres.DB
	)

	BeforeEach(func() {
		var err error
		repository, err = postgres.NewDB(postgres.ConnectionDetails{
			User:     "postgres",
			Pass:     "root",
			Host:     "localhost",
			Port:     5432,
			Database: "postgres",
			SSLMode:  "disable",
		})

		Expect(err).To(Succeed())
	})

	It("saves the short URL in the database and retrieves it again", func() {
		aShortURL := &url.ShortURL{
			Hash:    "foo",
			LongURL: "https://google.com",
		}

		err := repository.Save(aShortURL)
		Expect(err).To(Succeed())

		retrievedShortURL, err := repository.FindByHash("foo")

		Expect(err).To(Succeed())
		Expect(retrievedShortURL.Hash).To(Equal(aShortURL.Hash))
		Expect(retrievedShortURL.LongURL).To(Equal(aShortURL.LongURL))
	})

	Context("when the short URL already exists in the database", func() {
		It("doesn't return an error", func() {
			aShortURL := &url.ShortURL{
				Hash:    "foo",
				LongURL: "https://google.com",
			}

			err := repository.Save(aShortURL)
			Expect(err).To(Succeed())

			err = repository.Save(aShortURL)
			Expect(err).To(Succeed())
		})
	})

	Context("when the short URL doesn't exist in the database", func() {
		It("returns an error saying not found", func() {
			retrievedShortURL, err := repository.FindByHash("non_exising_hash")

			Expect(err).To(MatchError(url.ErrShortURLNotFound))
			Expect(retrievedShortURL).To(BeNil())
		})
	})
})
