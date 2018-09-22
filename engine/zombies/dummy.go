package zombies

import (
	"context"
	"log"
	"time"

	"github.com/sheirys/zombebattle/engine/types"
)

// Dummy is dumb zombie. He will always stand in same position. This has
// unlimited HP so you cannot kill him with arrows. However this zoombie still
// will react if you hit it.
type Dummy struct {
	name      string
	x, y      int
	events    chan types.Event
	heartbeat *time.Ticker
	ctx       context.Context
}

// Summon is used to initialize zombie and attach world events to it. Also
// context must be passed here to control how long this zombie should exist.
func (z *Dummy) Summon(e chan types.Event, ctx context.Context) error {
	// this is dummy zombie, so whatever name we will generate lets
	// attach "dummy-" in front of it so we will know that this zombie
	// cannot be killed.
	z.name = "dummy-" + pickName()
	z.events = e
	z.ctx = ctx

	// start living cycle.
	go z.startLiving()

	log.Printf("zombie '%s' has been summoned!", z.name)
	return nil
}

func (z *Dummy) GetName() string {
	return z.name
}

func (z *Dummy) Hit() {
	log.Printf("zombie '%s' got hit", z.name)
}

// Kill kills zombie. Killed zombie does not move.
// FIXME: implement state e.g.: zombie.IsAlive()
func (z *Dummy) Kill() error {
	z.heartbeat.Stop()
	return nil
}

func (z *Dummy) GetPos() (int, int) {
	return z.x, z.y
}

// Reset will move zombie to given position. But it does not respawn zombie if
// zombie is already died.
func (z *Dummy) Reset(x, y int) error {
	z.x = x
	z.y = y
	return nil
}

/*
func (z *Dummy) ProcessEvent(e types.Event) {
	log.Printf("zombie '%s' got event %s", z.name, e.String())
	// if event is shoot and arrow coordinates match zoombie coordinates
	// then this zoombie is hit.
	if e.Type == types.EventShoot && e.X == z.x && e.Y == z.y {
		log.Printf("zombie '%s' got hit", z.name)

		// zombie does not know how many points player has, so we will
		// attach points in the room.
		z.events <- types.Event{
			Type:  types.EventBoom,
			Actor: e.Actor,
			Hits:  []string{z.name},
		}
	}

}
*/

func (z *Dummy) startLiving() {
	z.heartbeat = time.NewTicker(3 * time.Second)
	for {
		select {
		case <-z.heartbeat.C:
			move := z.nextMove()
			log.Printf("zombie '%s' has moved '%s'", z.name, move.String())
			z.events <- move
		case <-z.ctx.Done():
			z.heartbeat.Stop()
			return
		}
	}
}

func (z *Dummy) nextMove() types.Event {
	// dummy will always stand in same position
	z.x = 5
	z.y = 5
	return types.Event{
		Type:  "WALK",
		Actor: z.name,
		X:     z.x,
		Y:     z.y,
	}
}
