package click_test

import (
	"github.com/WebEngineeringGroupI/backend/pkg/domain/click"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Click logger", func() {
	var (
		logger     *click.Logger
		repository click.ClickLoggerRepository
		shortURL   url.ShortURL
	)

	BeforeEach(func() {
		repository = &FakeClickLoggerRepository{logs: map[string][]string{}}
		logger = click.NewLogger(repository)
		shortURL = url.ShortURL{Hash: "12345678", LongURL: "https://google.com"}
	})

	Context("when providing a short URL and an @IP", func() {
		It("stores the the @IP in a repository", func() {
			ip := "192.168.1.1"

			logger.LogIP(&shortURL, ip)
			log := repository.FindByHash(shortURL.Hash)
			Expect(log).To(ContainElement(ip))
		})
	})

})

type FakeClickLoggerRepository struct {
	logs map[string][]string
}

func (f *FakeClickLoggerRepository) Save(url *url.ShortURL, ip string) {
	f.logs[url.Hash] = append(f.logs[url.Hash], ip)
}

func (f *FakeClickLoggerRepository) FindByHash(hash string) []string {
	log, ok := f.logs[hash]
	if !ok {
		return nil
	}

	return log
}
