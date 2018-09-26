package rooms

import (
	"context"

	"github.com/sheirys/zombebattle/engine/types"
)

// TrainingGrounds satisfies engine.Room interface and can be used as playable
// room. This room implementation can handle multiple players and zombies at
// same time.
type TrainingGrounds struct {
	Zombies      []types.Zombie
	players      []types.Player
	playerEvents chan types.Event
	zombieEvents chan types.Event
	ctx          context.Context
	stopFunc     context.CancelFunc
	name         string
}

// Name will return rooms name.
func (p *TrainingGrounds) Name() string {
	return p.name
}

// SetName will set name for this room.
func (p *TrainingGrounds) SetName(n string) {
	p.name = n
}

// AddPlayer will attach client to this room.
func (p *TrainingGrounds) AddPlayer(player types.Player) error {
	player.Notify(p.hello())
	go func() {
		p.players = append(p.players, player)
		for {
			if event, open := player.GetEvent(); open {
				p.playerEvents <- event
			} else {
				return
			}
		}
	}()
	return nil
}

// AddZombie will attach zombie to this room.
func (p *TrainingGrounds) AddZombie(z types.Zombie) error {
	p.Zombies = append(p.Zombies, z)
	z.Summon(p.ctx, p.zombieEvents)
	z.Run()
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

// Init will do some room preparations.
func (p *TrainingGrounds) Init() error {
	if p.name == "" {
		p.name = "TRAINING-GROUNDS"
	}
	p.zombieEvents = make(chan types.Event)
	p.playerEvents = make(chan types.Event)
	p.ctx, p.stopFunc = context.WithCancel(context.Background())

	// summon all pre-defined zombies.
	for _, zombie := range p.Zombies {
		zombie.Summon(p.ctx, p.zombieEvents)
		zombie.Run()
	}
	return nil
}

// Run will initialize this room.
func (p *TrainingGrounds) Run() error {
	go func() {
		for {
			if err := p.Process(); err != nil {
				return
			}
		}
	}()
	return nil
}

func (p *TrainingGrounds) Process() error {
	select {
	case playerEvent := <-p.playerEvents:
		switch playerEvent.Type {
		case types.EventShoot:
			booms := p.processShootEvent(playerEvent)
			p.sendEventToPlayers(booms)
		default:
		}
	case zombieEvent := <-p.zombieEvents:
		p.sendEventToPlayers(zombieEvent)
	}
	return nil
}

// ZombiesWon will always return false here, because we cannot win in training.
func (p *TrainingGrounds) ZombiesWon() bool {
	return false
}

// PlayersWon will always return false here, because we cannot win in training.
func (p *TrainingGrounds) PlayersWon() bool {
	return false
}

func (p *TrainingGrounds) sendEventToPlayers(e types.Event) {
	for _, player := range p.players {
		go player.ProcessEvent(e)
	}
}

func (p *TrainingGrounds) processShootEvent(e types.Event) types.Event {
	hits := []string{}
	for _, zombie := range p.Zombies {
		x, y := zombie.GetPos()
		if x == e.X && y == e.Y {
			hits = append(hits, zombie.GetName())
			zombie.Hit()
		}
	}
	shootResult := types.Event{
		Type:   types.EventBoom,
		Actor:  e.Actor,
		Points: len(hits),
		Hits:   hits,
	}
	return shootResult
}

func (p *TrainingGrounds) hello() string {
	msg := "# You appeared in sandy yard. Sharp stones are \n"
	msg += "# tickling your legs. You feel uncomfortable. In\n"
	msg += "# front of you there is a sign: \n"
	msg += "# \n"
	msg += "# ================================================ \n"
	msg += "# " + p.name + "\n"
	msg += "# Welcome to the training grounds. Here you can try\n"
	msg += "# to kill zombies with your friends. Use your bow\n"
	msg += "# with `SHOOT <x> <y>` command if you want to kill\n"
	msg += "# a zombie. Be awere that some zoombies needs to \n"
	msg += "# be hit twice or more before they die.\n"
	msg += "# ================================================ \n"
	return msg
}
