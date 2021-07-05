package ui

import (
	hw "github.com/jonathangjertsen/serious/hw"
)

type Ui interface {
	Run(hw.Hw)
}
