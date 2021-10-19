package click_test

import (
	"github.com/WebEngineeringGroupI/backend/pkg/domain/click"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Click logger", func() {
	var (
		clicker     *click.Clicker
		repository click.ClickerRepository
		aShortURL *url.ShortURL
	)

	BeforeEach(func() {
		repository = &FakeClickerRepository{clicks: map[string][]*click.ClickDetails{}}
		clicker = click.NewClicker(repository)
		aShortURL = &url.ShortURL{ Hash: "12345678", LongURL: "https://google.com"}
	})

	Context("when providing click details", func() {
		It("logs click details in a repository", func() {

			click := &click.ClickDetails {
				Hash: aShortURL.Hash,
				Ip : "192.168.1.1",
			}

			clicker.LogClick(click)
			clicks := repository.FindClicksByHash(aShortURL.Hash)
			Expect(clicks).To(ContainElement(click))
		})
	})

})

type FakeClickerRepository struct {
	clicks map[string][]*click.ClickDetails
}

func (f *FakeClickerRepository) SaveClick(click *click.ClickDetails) {
	f.clicks[click.Hash] = append(f.clicks[click.Hash], click)
}

func (f *FakeClickerRepository) FindClicksByHash(hash string) []*click.ClickDetails {
	clicks, ok := f.clicks[hash]
	if !ok {
		return nil
	}

	return clicks
}
