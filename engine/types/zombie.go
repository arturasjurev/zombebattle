package types

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
	Summon(ctx context.Context, moves chan Event) error

	Run()

	// GetName will return zombies name.
	GetName() string

	// Kill zombie now. This will be called when room decides that this
	// zombie should die.
	Kill() error

	// Reset will move zombie to given position. But it does not respawn
	// zombie if zombie is already died.
	Reset(x, y int64) error

	// GetPos should return position of zombie.
	GetPos() (x, y int64)

	// Hit will be called when zombie gets an arrow in his ass.
	Hit()

	// Next should force zombie to move.
	Next()
}
