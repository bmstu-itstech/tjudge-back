package shared

type Event interface {
	Name() string
	IsEvent()
}
