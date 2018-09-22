package types

type Player interface {
	Run()
	WaitForStart(server chan Event) error
	ShowLobby([]Lobby)
	Notify(msg string)
	SelectedRoom() string
	GetEvent() <-chan Event
	ProcessEvent(e Event)
}
