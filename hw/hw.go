package hw

type Hw interface {
	GetPorts() []string
	Selected() (int, *string)
}
