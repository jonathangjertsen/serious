package ui

import (
	messages "github.com/jonathangjertsen/serious/messages"
	"time"
)

type Ui interface {
	Run(channel *chan messages.Message)
	StartReadTask(time.Duration)
}
