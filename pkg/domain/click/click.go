package click

type Clicker struct {
	repository ClickerRepository
}

type ClickDetails struct {
	Hash    string
	Ip string
}

func NewClicker(repository ClickerRepository) *Clicker {
	return &Clicker{repository: repository}
}

func (l *Clicker) LogClick(click *ClickDetails) {
	l.repository.SaveClick(click)
}

