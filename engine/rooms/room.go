package rooms

import (
	"context"
	"sync"

	"github.com/sheirys/zombebattle/engine"
)

// PlayRoom satisfies engine.Room interface and can be used as playable room.
// This room implementation can handle multiple players and zombies at same time.
type PlayRoom struct {
	zombies   []engine.Zombie
	events    chan engine.Event
	scores    map[string]int
	scoresMtx sync.Mutex
	ctx       context.Context
}

// Init() will initialize various variables required for functional room. This
// should be called before use of room.
// FIXME: implement something like NewRoom() where Init() will be called before
// return.
func (p *PlayRoom) Init() error {
	p.scores = make(map[string]int)
	p.events = make(chan engine.Event, 50) //FIXME: hardcoded for demo
	p.scoresMtx = sync.Mutex{}
	return nil
}

// Join will register player in this room. This can be called only before Start()
func (p *PlayRoom) Join(name string) error {
	// FIXME: check if we are calling this before Start() function.
	p.pointsSet(name, 0)
	return nil
}

// Stop stops this room and kills all zombies.
func (p *PlayRoom) Stop() error {
	close(p.events)
	return nil
}

// AddZombie registers new zombie into this room. This can be called only when
// room does not entered start state.
func (p *PlayRoom) AddZombie(z engine.Zombie) error {
	// FIXME: check if we are calling this before Start() function.
	p.zombies = append(p.zombies, z)
	return nil
}

func (p *PlayRoom) AddEvent(e engine.Event) {
	p.events <- e
}

// Shoot should be called when player tries to shoot zombie.
func (p *PlayRoom) Shoot(name string, target engine.Pos) error {
	for _, zombie := range p.zombies {
		// FIXME: mutex in zombies should be used because here possible
		// race condition can appear. When we are looping zoombies
		// positions may be changed, but for demo Im fine with this.
		if zombie.Pos() == target {
			zombie.Reset(engine.Pos{0, 0})
			points := p.pointsInc(name)
			p.events <- engine.Event{
				Type:   "BOOM",
				Actor:  name,
				Points: points,
			}
		}
	}
	return nil
}

// EventStream will stream all room events.
func (p *PlayRoom) EventStream() <-chan engine.Event {
	return p.events
}

func (p *PlayRoom) Scores() engine.Scoreboard {
	return nil
}

// start will initialize all zombies
func (p *PlayRoom) start(ctx context.Context) error {
	p.ctx = ctx
	for _, zombie := range p.zombies {
		// summon all registered zombies and attach their events into
		// our events events stream.
		go zombie.Summon(p.events, ctx)
	}
	return nil
}

// pointsSet will set points for player in scoreboard. This should be called when
// new player joins our room to initialize points for him.
func (p *PlayRoom) pointsSet(name string, points int) {
	p.scoresMtx.Lock()
	p.scores[name] = points
	p.scoresMtx.Unlock()
}

// pointsInc will increment points by 1 for provided player.
func (p *PlayRoom) pointsInc(name string) (points int) {
	p.scoresMtx.Lock()
	p.scores[name]++
	points = p.scores[name]
	p.scoresMtx.Unlock()
	return
}
