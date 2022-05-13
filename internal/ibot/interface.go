package ibot

type Interface interface {
	SendMessage(text string) error
}
