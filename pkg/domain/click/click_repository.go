package click

type ClickerRepository interface {
	FindClicksByHash(hash string) []*Details
	SaveClick(click *Details)
}
