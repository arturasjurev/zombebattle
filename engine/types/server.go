package types

// ServerRoom defines rooms holded by server. Room should be pointer to room and
// default defines if this room is default room in server.
type ServerRoom struct {
	Room    Room
	Default bool
}

// Lobby defines what rooms are registered in server. This struct is returned to
// client, when we want to inform him, what rooms are available at this moment.
type Lobby struct {
	Name    string
	Default bool
}
