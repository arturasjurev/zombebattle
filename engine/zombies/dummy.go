package zombies

import (
	"context"
	"log"
	"sync/atomic"
	"time"

	"github.com/sheirys/zombebattle/engine/types"
)

// Dummy is dumb zombie. He will always stand in same position. This has
// unlimited HP so you cannot kill him with arrows. However this zoombie still
// will react if you hit it.
type Dummy struct {
	name      string
	x, y      int64
	events    chan types.Event
	heartbeat *time.Ticker
	ctx       context.Context
}

// Summon is used to initialize zombie and attach world events to it. Also
// context must be passed here to control how long this zombie should exist.
func (z *Dummy) Summon(ctx context.Context, e chan types.Event) error {
	// this is dummy zombie, so whatever name we will generate lets
	// attach "dummy-" in front of it so we will know that this zombie
	// cannot be killed.
	z.name = "dummy-" + PickName()
	z.events = e
	z.ctx = ctx

	log.Printf("zombie '%s' has been summoned!", z.name)
	return nil
}

// Run will start this zombie.
func (z *Dummy) Run() {
	// start living cycle.
	go z.startLiving()
}

// GetName will return zombie name.
func (z *Dummy) GetName() string {
	return z.name
}

// Hit will be called when player hits this zombie.
func (z *Dummy) Hit() {
	log.Printf("zombie '%s' got hit", z.name)
}

// Kill kills zombie. Killed zombie does not move.
// FIXME: implement state e.g.: zombie.IsAlive()
func (z *Dummy) Kill() error {
	z.heartbeat.Stop()
	return nil
}

// GetPos will return current zombie position.
func (z *Dummy) GetPos() (int64, int64) {
	return atomic.LoadInt64(&z.x), atomic.LoadInt64(&z.y)
}

// Reset will move zombie to given position. But it does not respawn zombie if
// zombie is already died.
func (z *Dummy) Reset(x, y int64) error {
	atomic.StoreInt64(&z.x, x)
	atomic.StoreInt64(&z.y, y)
	return nil
}

// Next will force this zombie to move into next location.
func (z *Dummy) Next() {
	z.move()
}

func (z *Dummy) move() {
	move := z.nextMove()
	log.Printf("zombie '%s' has moved '%s'", z.name, move.String())
	z.events <- move
}

func (z *Dummy) startLiving() {
	z.heartbeat = time.NewTicker(3 * time.Second)
	for {
		select {
		case <-z.heartbeat.C:
			z.Next()
		case <-z.ctx.Done():
			z.heartbeat.Stop()
			return
		}
	}
}

func (z *Dummy) nextMove() types.Event {
	// dummy will always stand in same position
	atomic.StoreInt64(&z.x, 5)
	atomic.StoreInt64(&z.y, 5)

	return types.Event{
		Type:  "WALK",
		Actor: z.name,
		X:     atomic.LoadInt64(&z.x),
		Y:     atomic.LoadInt64(&z.y),
	}
}
