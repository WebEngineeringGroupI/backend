package click

type ClickerRepository interface {
	FindClicksByHash(hash string) []*ClickDetails
	SaveClick(click *ClickDetails)
}