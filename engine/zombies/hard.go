package zombies

import (
	"context"
	"time"

	"github.com/sheirys/zombebattle/engine"
)

// Hard holds rabbit zombie. This is basicly copy/paste zoombie struct from
// zombies/easy.go but with different movement logic (random instead of linear).
type Hard struct {
	name      string
	events    chan engine.Event
	heartbeat *time.Ticker
	pos       engine.Pos
	ctx       context.Context
}

// Summon is used to initialize zombie.
func (z *Hard) Summon(e chan engine.Event, ctx context.Context) error {
	z.name = pickName()
	z.events = e
	z.ctx = ctx
	go z.startLiving()
	return nil
}

// Kill kills zombie. Killed zombie does not move.
// FIXME: implement state e.g.: zombie.IsAlive()
func (z *Hard) Kill() error {
	z.heartbeat.Stop()
	return nil
}

// Reset will move zombie to given position. But it does not respawn zombie if
// zombie is already died.
func (z *Hard) Reset(p engine.Pos) error {
	z.pos = p
	return nil
}

// Hit zombie with arrow!
func (z *Hard) Hit() {}

// Pos will return current zombie position
func (z *Hard) Pos() engine.Pos {
	return z.pos
}

func (z *Hard) startLiving() {
	// FIXME: Don't like hardcoded 3s here. Probably should be const
	// EasySpeed, HardSpeed etc.
	z.heartbeat = time.NewTicker(3 * time.Second)
	for {
		select {
		case <-z.heartbeat.C:
			z.move()
		case <-z.ctx.Done():
			z.heartbeat.Stop()
			return
		}
	}
}

func (z *Hard) move() {
	// this is not a zombie ?! this is zombie rabbit !! :panic:
	// jumps randomly
	z.pos.X = randomPos(1, 6)
	z.pos.Y = randomPos(1, 6)
	z.events <- engine.Event{
		Type:  "WALK",
		Actor: z.name,
		X:     z.pos.X,
		Y:     z.pos.Y,
	}
}
