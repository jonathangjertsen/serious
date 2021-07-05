package ui

import (
	messages "github.com/jonathangjertsen/serious/messages"
)

type Ui interface {
	Run(channel *chan messages.Message)
}
