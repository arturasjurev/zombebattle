package rooms

import (
	"context"

	"github.com/sheirys/zombebattle/engine"
)

// TrainingGrounds satisfies engine.Room interface and can be used as playable
// room. This room implementation can handle multiple players and zombies at
// same time.
type TrainingGrounds struct {
	Zombies      []engine.Zombie
	players      []*engine.Client
	playerEvents chan engine.Event
	zombieEvents chan engine.Event
	ctx          context.Context
	stopFunc     context.CancelFunc
	name         string
}

// Name will return rooms name
func (p *TrainingGrounds) Name() string {
	return p.name
}

// AddPlayer will attach client to this room.
func (p *TrainingGrounds) AddPlayer(c *engine.Client) error {
	go func() {
		p.players = append(p.players, c)
		for {
			event, open := <-c.GetEvent()
			if !open {
				return
			}
			p.playerEvents <- event
		}
	}()
	return nil
}

// AddZombie will attach zombie to this room.
func (p *TrainingGrounds) AddZombie(z engine.Zombie) error {
	p.Zombies = append(p.Zombies, z)
	go z.Summon(p.zombieEvents, p.ctx)
	return nil
}

// Stop stops this room and kills all zombies.
func (p *TrainingGrounds) Stop() error {
	defer func() {
		close(p.zombieEvents)
		close(p.playerEvents)
	}()
	p.stopFunc()
	return nil
}

// Run will initialize this room.
func (p *TrainingGrounds) Run() error {
	p.init()
	go func() {
		for {
			select {
			case playerEvent := <-p.playerEvents:
				for _, zombie := range p.Zombies {
					go zombie.ProcessEvent(playerEvent)
				}
			case zombieEvent := <-p.zombieEvents:
				for _, player := range p.players {
					go player.ProcessEvent(zombieEvent)
				}
			}
		}
	}()
	return nil
}

// init will do some room preparations.
func (p *TrainingGrounds) init() error {
	p.name = "training grounds"
	p.zombieEvents = make(chan engine.Event)
	p.playerEvents = make(chan engine.Event)
	p.ctx, p.stopFunc = context.WithCancel(context.Background())

	// summon all pre-defined zombies.
	for _, zombie := range p.Zombies {
		zombie.Summon(p.zombieEvents, p.ctx)
	}
	return nil
}
