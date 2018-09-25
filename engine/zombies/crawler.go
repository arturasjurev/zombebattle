package zombies

import (
	"context"
	"log"
	"sync/atomic"
	"time"

	"github.com/sheirys/zombebattle/engine/types"
)

// Crawler is a zombie that shiuld be used in TheWall room. This zombie will
// try to climb on the wall and kill archer. When zombie reaches wall it is room
// responsibility to kill that zombie and respawn it.
type Crawler struct {
	name       string
	x, y       int64
	events     chan types.Event
	timeToMove *time.Ticker
	ctx        context.Context
}

// Summon is used to initialize zombie and attach world events to it. Also
// context must be passed here to control how long this zombie should exist.
func (z *Crawler) Summon(ctx context.Context, e chan types.Event) error {
	// lets randomly generate name for this zombie. As we want to know what
	// type of zombie this is, we will hardcode `crawler` in front of the
	// name.
	z.name = "crawler-" + PickName()
	z.events = e
	z.ctx = ctx

	log.Printf("zombie '%s' has been summoned!", z.name)
	return nil
}

// Run will start this zombie.
func (z *Crawler) Run() {
	// start living cycle.
	go z.startLiving()
}

// GetName will return zombie name.
func (z *Crawler) GetName() string {
	return z.name
}

// Hit will be called when player hits this zombie. As this zombie should be
// used in TheWall room, it is room responsibility to kill and respawn zombie.
func (z *Crawler) Hit() {
	log.Printf("zombie '%s' got hit", z.name)
}

// Kill zombie. Room should call this when we want to force-kill this zombie.
func (z *Crawler) Kill() error {
	z.timeToMove.Stop()
	return nil
}

// GetPos will return current zombie position.
func (z *Crawler) GetPos() (int64, int64) {
	return atomic.LoadInt64(&z.x), atomic.LoadInt64(&z.y)
}

// Reset will move zombie to given position. But it does not respawn zombie if
// zombie is already died.
func (z *Crawler) Reset(x, y int64) error {
	atomic.StoreInt64(&z.x, x)
	atomic.StoreInt64(&z.y, y)
	return nil
}

// Next will force to move this zombie into next location.
func (z *Crawler) Next() {
	z.move()
}

func (z *Crawler) move() {
	move := z.nextMove()
	log.Printf("zombie '%s' has moved '%s'", z.name, move.String())
	z.events <- move
}

func (z *Crawler) startLiving() {
	z.timeToMove = time.NewTicker(3 * time.Second)
	for {
		select {
		case <-z.timeToMove.C:
			z.Next()
		case <-z.ctx.Done():
			z.timeToMove.Stop()
			return
		}
	}
}

func (z *Crawler) nextMove() types.Event {
	// we expect that wall will be reached when zoombie will reach X axis
	// 0 point. So we will decrease x on every move.
	atomic.AddInt64(&z.x, -1)

	return types.Event{
		Type:  "WALK",
		Actor: z.name,
		X:     atomic.LoadInt64(&z.x),
		Y:     atomic.LoadInt64(&z.y),
	}
}
