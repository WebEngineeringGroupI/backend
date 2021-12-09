package postgres_test

import (
	"context"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/click"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/postgres"
)

var _ = Describe("Postgres", func() {
	var (
		repository *postgres.DB
		ctx        context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		var err error
		repository, err = postgres.NewDB(postgres.ConnectionDetails{
			User:     "postgres",
			Pass:     "postgres",
			Host:     "localhost",
			Port:     5432,
			Database: "postgres",
			SSLMode:  "disable",
		})

		Expect(err).ToNot(HaveOccurred())
	})

	It("saves the short URL in the database and retrieves it again", func() {
		err := repository.SaveShortURL(ctx, aShortURL())
		Expect(err).ToNot(HaveOccurred())

		retrievedShortURL, err := repository.FindShortURLByHash(ctx, "12345678")

		Expect(err).ToNot(HaveOccurred())
		Expect(retrievedShortURL.Hash).To(Equal("12345678"))
		Expect(retrievedShortURL.OriginalURL).To(Equal(url.OriginalURL{URL: "https://google.com", IsValid: true}))
	})

	Context("when the short URL already exists in the database", func() {
		It("doesn't return an error", func() {
			err := repository.SaveShortURL(ctx, aShortURL())
			Expect(err).ToNot(HaveOccurred())

			err = repository.SaveShortURL(ctx, aShortURL())
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when the short URL doesn't exist in the database", func() {
		It("returns an error saying not found", func() {
			retrievedShortURL, err := repository.FindShortURLByHash(ctx, "non_existing_hash")

			Expect(err).To(MatchError(url.ErrShortURLNotFound))
			Expect(retrievedShortURL).To(BeNil())
		})
	})

	It("saves a load balanced URL and retrieves it again", func() {
		err := repository.SaveLoadBalancedURL(ctx, aLoadBalancedURL())
		Expect(err).ToNot(HaveOccurred())

		loadBalancedURL, err := repository.FindLoadBalancedURLByHash(ctx, "12345678")
		Expect(err).ToNot(HaveOccurred())
		Expect(loadBalancedURL.Hash).To(Equal("12345678"))
		Expect(loadBalancedURL.LongURLs).To(ConsistOf(
			url.OriginalURL{URL: "https://google.com", IsValid: false},
			url.OriginalURL{URL: "https://youtube.com", IsValid: false},
			url.OriginalURL{URL: "https://reddit.com", IsValid: true},
		))
	})

	Context("when the load balanced URL already exists in the database", func() {
		It("doesn't return an error", func() {
			err := repository.SaveLoadBalancedURL(ctx, aLoadBalancedURL())
			Expect(err).ToNot(HaveOccurred())

			err = repository.SaveLoadBalancedURL(ctx, aLoadBalancedURL())
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when the load balanced URL doesn't exist in the database", func() {
		It("returns an error saying not found", func() {
			loadBalancedURL, err := repository.FindLoadBalancedURLByHash(ctx, "non_existing_hash")

			Expect(err).To(MatchError(url.ErrValidURLNotFound))
			Expect(loadBalancedURL).To(BeNil())
		})
	})

	It("Stores click information and retrieves again", func() {
		click := &click.Details{
			Hash: "12345678",
			IP:   "192.168.1.1",
		}
		err := repository.SaveClick(ctx, click)
		Expect(err).ToNot(HaveOccurred())

		clicks, err := repository.FindClicksByHash(ctx, click.Hash)
		Expect(err).ToNot(HaveOccurred())
		Expect(clicks).To(ContainElement(click))
	})

	Context("when click information doesn't exist in the database", func() {
		It("doesn't return an error", func() {
			clicks, err := repository.FindClicksByHash(ctx, "non_existing_hash")
			Expect(err).ToNot(HaveOccurred())
			Expect(clicks).To(BeEmpty())
		})
	})

})

func aShortURL() *url.ShortURL {
	return &url.ShortURL{
		Hash:        "12345678",
		OriginalURL: url.OriginalURL{URL: "https://google.com", IsValid: true},
	}
}

func aLoadBalancedURL() *url.LoadBalancedURL {
	return &url.LoadBalancedURL{
		Hash: "12345678",
		LongURLs: []url.OriginalURL{
			{URL: "https://google.com", IsValid: false},
			{URL: "https://youtube.com", IsValid: false},
			{URL: "https://reddit.com", IsValid: true},
		},
	}
}
