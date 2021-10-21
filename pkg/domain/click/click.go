package click

type Clicker struct {
	repository ClickerRepository
}

type Details struct {
	Hash string
	IP   string
}

func NewClicker(repository ClickerRepository) *Clicker {
	return &Clicker{repository: repository}
}

func (l *Clicker) LogClick(clickDetails *Details) {
	l.repository.SaveClick(clickDetails)
}
