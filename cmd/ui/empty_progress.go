package ui

type EmptyProgressBar struct{}

func (p *EmptyProgressBar) Init(count int) {

}

func (p *EmptyProgressBar) Add(count int) {}
