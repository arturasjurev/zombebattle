package engine

import "errors"

var (
	ErrBadInput = errors.New("bad command, won't parse")
)

// CommunicationChannel defines how communication should behave. We need this
// interface, as default Communication Channel protocol does not specify how
// multiple players should start game, so we will need to implement other
// CommunicationChannelExtended protocol. This interface will allow use to
// choose different protocol for server.
type CommunicationChannel interface {
	// Parse should read bytes and return decoded Event
	Parse(b []byte) (Event, error)
}
