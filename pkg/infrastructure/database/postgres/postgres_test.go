package postgres_test

import (
	"github.com/WebEngineeringGroupI/backend/pkg/domain/click"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/postgres"
)

var _ = Describe("Postgres", func() {
	var (
		aShortURL  *url.ShortURL
		repository *postgres.DB
	)

	BeforeEach(func() {
		aShortURL = &url.ShortURL{
			Hash:        "12345678",
			OriginalURL: url.OriginalURL{URL: "https://google.com", IsValid: true},
		}

		var err error
		repository, err = postgres.NewDB(postgres.ConnectionDetails{
			User:     "postgres",
			Pass:     "postgres",
			Host:     "localhost",
			Port:     5432,
			Database: "postgres",
			SSLMode:  "disable",
		})

		Expect(err).To(Succeed())
	})

	It("saves the short URL in the database and retrieves it again", func() {
		err := repository.Save(aShortURL)
		Expect(err).To(Succeed())

		retrievedShortURL, err := repository.FindByHash("12345678")

		Expect(err).To(Succeed())
		Expect(retrievedShortURL.Hash).To(Equal(aShortURL.Hash))
		Expect(retrievedShortURL.OriginalURL).To(Equal(aShortURL.OriginalURL))
	})

	Context("when the short URL already exists in the database", func() {
		It("doesn't return an error", func() {
			err := repository.Save(aShortURL)
			Expect(err).To(Succeed())

			err = repository.Save(aShortURL)
			Expect(err).To(Succeed())
		})
	})

	Context("when the short URL doesn't exist in the database", func() {
		It("returns an error saying not found", func() {
			retrievedShortURL, err := repository.FindByHash("non_existing_hash")

			Expect(err).To(MatchError(url.ErrShortURLNotFound))
			Expect(retrievedShortURL).To(BeNil())
		})
	})

	It("Stores click information and retrieves again", func() {
		click := &click.Details{
			Hash: aShortURL.Hash,
			IP:   "192.168.1.1",
		}
		err := repository.SaveClick(click)
		Expect(err).To(Succeed())

		clicks, err := repository.FindClicksByHash(click.Hash)
		Expect(err).To(Succeed())
		Expect(clicks).To(ContainElement(click))
	})

	Context("when click information doesn't exist in the database", func() {
		It("doesn't return an error", func() {
			clicks, err := repository.FindClicksByHash("non_existing_hash")
			Expect(err).To(Succeed())
			Expect(clicks).To(BeEmpty())
		})
	})

})
