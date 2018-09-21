package zombies

import (
	"context"
	"log"
	"time"

	"github.com/sheirys/zombebattle/engine"
)

type Easy struct {
	name      string
	events    chan engine.Event
	heartBeat *time.Ticker
	pos       engine.Pos
	ctx       context.Context
}

// Summon is used to initialize zombie.
func (z *Easy) Summon(e chan engine.Event, ctx context.Context) error {
	z.name = pickName()
	z.events = e
	z.ctx = ctx
	go z.startLiving()
	log.Printf("zombie '%s' has been summoned!", z.name)
	return nil
}

// Kill zombie. Killed zombie does not move.
// FIXME: implement state e.g.: zombie.IsAlive()
func (z *Easy) Kill() error {
	z.heartBeat.Stop()
	return nil
}

// Reset will move zombie to given position. But it does not respawn zombie if
// zombie is already died.
func (z *Easy) Reset(p engine.Pos) error {
	z.pos = p
	return nil
}

// Hit zombie with arrow!
func (z *Easy) Hit() {}

func (z *Easy) ProcessEvent(e engine.Event) {
	log.Printf("zombie %s got event %s", z.name, e.String())
}

// Pos will return current zombie position
func (z *Easy) Pos() engine.Pos {
	return z.pos
}

func (z *Easy) startLiving() {
	// FIXME: Don't like hardcoded 3s here. Probably should be const
	// EasySpeed, HardSpeed etc.
	z.heartBeat = time.NewTicker(3 * time.Second)
	for {
		select {
		case <-z.heartBeat.C:
			move := z.nextMove()
			log.Printf("zombie '%s' has moved '%s'", z.name, move.String())
			z.events <- move
		case <-z.ctx.Done():
			z.heartBeat.Stop()
			return
		}
	}
}

func (z *Easy) nextMove() engine.Event {
	// omg, omg this dumb Easy zombie can move only in X axis !!
	z.pos.X++
	return engine.Event{
		Type:  "WALK",
		Actor: z.name,
		X:     z.pos.X,
		Y:     z.pos.Y,
	}
}
