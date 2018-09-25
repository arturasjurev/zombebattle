package types

type Player interface {
	Notify(msg string)
	GetEvent() (Event, bool)
	ProcessEvent(e Event)
	ProduceEvent(e Event)
}
