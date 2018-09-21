package engine

import "context"

// Zombie defines what we expect from zombie and how we are going to control him.
// We will have different type of zombies;- a dumb zombie which will move only
// in x axis (zombie.Easy), a rabbit zombie which will jump in random coordinates
// (we call it zombie.Hard). This interface allows us to implement different
// kind of zombies like twitter zombie which will move by some random tweets or
// telnet zombie which can be controlled by network, because you know, why not?
type Zombie interface {

	// Summon will spanw a zombie and all his movements/events will be sent
	// to e chan. Context can be used to stop zombies.
	Summon(e chan Event, ctx context.Context) error

	ProcessEvent(e Event)
	// Kill zombie now. This will be called when room decides that this
	// zombie should die.
	Kill() error

	// Hit means that this zombie now has an arrow in his ass.
	Hit()

	// Reset will move zombie to given position. But it does not respawn
	// zombie if zombie is already died.
	Reset(p Pos) error

	// Post will return current zombie possition.
	Pos() Pos
}
