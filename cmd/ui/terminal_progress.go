package ui

import "github.com/schollz/progressbar/v3"

type TerminalWidgetProgressBar struct {
	bar *progressbar.ProgressBar
}

func (p *TerminalWidgetProgressBar) Init(count int) {
	p.bar = progressbar.Default(int64(count))
}

func (p *TerminalWidgetProgressBar) Add(count int) {
	if p.bar != nil {
		err := p.bar.Add(count)
		if err != nil {
		}
	}
}
