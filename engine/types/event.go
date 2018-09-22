package types

import "fmt"

// Define expected events. Events are used in rooms to describe what is
// happening.
const (

	// default communication channel commands. These commands are defined
	// in requirements and should always work.
	EventWalk  = "WALK"  // when zombie walks
	EventShoot = "SHOOT" // when player shoots
	EventBoom  = "BOOM"  // when zombie dies
	EventStart = "START" // when client wants start game.

	// extended communication channel commands. These commands extends
	// default communication channel protocol. As there is no documentation
	// about how allow mutliroom server to join multiple clients, we need
	// additional commands to join room or create a new one.
	EventJoin = "JOIN" // join to given room `JOIN woods`
	EventNew  = "NEW"  // create new room
)

// Event will be used for various events in this engine. For example if player
// tries to shoot zombie, when zombie dies or walks. All events should be
// processed into room. Each room implemeation can interpretate events
// differently.
type Event struct {
	Type   string
	Actor  string
	X, Y   int
	Points int
	Hits   []string
}

// String will convert event into human readable string. E.g.:
//
//	WALK zombie 1 7
//
func (e *Event) String() (s string) {
	switch e.Type {
	case EventWalk:
		s = fmt.Sprintf("%s %s %d %d", e.Type, e.Actor, e.X, e.Y)
	case EventShoot:
		s = fmt.Sprintf("%s %d %d", e.Type, e.X, e.Y)
	case EventBoom:
		s = fmt.Sprintf("%s %s %d %v", e.Type, e.Actor, e.Points, e.Hits)
	}
	return
}
