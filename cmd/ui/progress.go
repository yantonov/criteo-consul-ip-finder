package ui

type ProgressBar interface {
	Init(count int)
	Add(count int)
}
