package rooms

import (
	"context"
	"errors"
	"log"
	"sync/atomic"

	"github.com/sheirys/zombebattle/engine/types"
	"github.com/sheirys/zombebattle/engine/zombies"
)

// General settings for this room.
const (
	TheWallMapWidth       = 29 // map size 10x30 but we count from 0 here.
	TheWallMapHeight      = 9  // map size
	TheWallMaxPlayerScore = 5  // zombies needs to kill before victory
	TheWallMaxZombieScore = 5  // zombies reach the wall before game over
)

// TheWall satisfies engine.Room interface and can be used as playable room.
// This room has a wall on X0 axis where zombies tries to reach it from right
// side. This room should use only zoombies.Crawler type zombies. When zombie
// reaches the wall, room should kill him, and respawn into random position.
// By default, this room does not have predefined zombies. If server requires it
// then these zombies will be spawned. We will spawn a new zombie each time when
// player joins this room.
type TheWall struct {
	Zombies      []types.Zombie
	players      []types.Player
	playerEvents chan types.Event
	zombieEvents chan types.Event
	name         string

	// room settings
	width, height int64 // map size
	playerScore   int64 // how many zombies must be killed before win?
	zombieScore   int64 // how many times wall can be reached by zombies?

	// room systems
	ctx      context.Context
	stopFunc context.CancelFunc
	running  bool
}

// Name will return room name.
func (p *TheWall) Name() string {
	return p.name
}

// SetName will set name for this room.
func (p *TheWall) SetName(n string) {
	p.name = n
}

// AddPlayer will attach client to this room. Everytime when we attach new
// player, new crawler will be spawned.
func (p *TheWall) AddPlayer(player types.Player) error {
	p.players = append(p.players, player)
	player.Notify(p.hello())

	// check scores. Maybe this room is already in end state.
	p.checkScores()

	// add zombies only then, when we do not have a winner of this room.
	if p.running {
		crawler := &zombies.Crawler{}
		p.AddZombie(crawler)
	}

	go func() {
		for {
			// handle player events
			if event, open := player.GetEvent(); open && p.running {
				p.playerEvents <- event
				continue
			}
			return
		}
	}()
	return nil
}

// AddZombie will attach zombie to this room.
func (p *TheWall) AddZombie(z types.Zombie) error {
	z.Reset(p.width, zombies.RandomPos(0, p.height))
	p.Zombies = append(p.Zombies, z)
	z.Summon(p.ctx, p.zombieEvents)
	z.Run()
	return nil
}

// Init will do some room preparations.
func (p *TheWall) Init() error {
	if p.name == "" {
		p.name = "THE-WALL"
	}
	p.zombieEvents = make(chan types.Event, 1)
	p.playerEvents = make(chan types.Event, 1)
	p.ctx, p.stopFunc = context.WithCancel(context.Background())

	p.width = TheWallMapWidth
	p.height = TheWallMapHeight
	p.running = true

	// summon all pre-defined zombies.
	for _, zombie := range p.Zombies {
		zombie.Reset(p.width, zombies.RandomPos(0, p.height))
		zombie.Summon(p.ctx, p.zombieEvents)
		zombie.Run()
	}
	return nil
}

// Stop stops this room and kills all zombies.
func (p *TheWall) Stop() error {
	defer func() {
		close(p.zombieEvents)
		close(p.playerEvents)
	}()
	p.stopFunc()
	return nil
}

// Run will initialize this room.
func (p *TheWall) Run() error {
	go func() {
		for {
			if err := p.Process(); err != nil {
				return
			}
		}
	}()
	return nil
}

func (p *TheWall) Process() error {

	select {
	// handle player event
	case playerEvent, ok := <-p.playerEvents:
		if !ok {
			// react when channel is closed
			return errors.New("player channel closed")
		}
		if playerEvent.Type == types.EventShoot {
			// return shot result to players
			booms := p.processShootEvent(playerEvent)
			p.sendEventToPlayers(booms)
		}
	// handle zombie event
	case zombieEvent, ok := <-p.zombieEvents:
		if !ok {
			// reach when channel is closed
			return errors.New("zombie channel closed")
		}
		// check maybe zombie reached the wall?
		p.processMoveEvent(zombieEvent)
		p.sendEventToPlayers(zombieEvent)
	}
	return nil
}

// ZombiesWon will return true if zombies won this room.
func (p *TheWall) ZombiesWon() bool {
	return p.getZombieScores() >= TheWallMaxPlayerScore
}

// PlayersWon will return true if players won this room.
func (p *TheWall) PlayersWon() bool {
	return p.getPlayerScores() >= TheWallMaxPlayerScore
}

func (p *TheWall) sendEventToPlayers(e types.Event) {
	for _, player := range p.players {
		go player.ProcessEvent(e)
	}
}

// processMoveEvent will check how zombies are moving and where they are. Here
// we will check if zombie reached the wall. If reached then add points to
// zombie team and respawn it on the left.
func (p *TheWall) processMoveEvent(e types.Event) {
	if e.X == 0 {
		p.incZombieScores()
		p.checkScores()
		log.Printf("zombie %s reached the wall", e.Actor)
		for _, zombie := range p.Zombies {
			if zombie.GetName() == e.Actor {
				zombie.Reset(p.width, zombies.RandomPos(0, p.height))
				break
			}
		}
	}
}

// processShootEvent will handle shoot event from player. Here we will check
// each zombie position and check if any zombies are hit. In the end we will
// produce BOOM event here wit points count and hit zombies.
func (p *TheWall) processShootEvent(e types.Event) types.Event {
	hits := []string{}
	for _, zombie := range p.Zombies {
		x, y := zombie.GetPos()
		if x == e.X && y == e.Y {
			hits = append(hits, zombie.GetName())
			if zombie.Hit() {
				zombie.Reset(p.width, zombies.RandomPos(0, p.height))
				p.incPlayerScores()
				p.checkScores()
			}
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

// checkScores should be called everytime when player hits a zombie or zombie
// reaches the wall. Here we will check how many scores has zombies vs players
// and decide if we need to continue this room, or someone wins.
// FIXME: implement this.
func (p *TheWall) checkScores() {
	log.Printf("scores for map %s", p.name)
	log.Printf("zombies has %d/%d points", p.zombieScore, TheWallMaxZombieScore)
	log.Printf("players has %d/%d points", p.playerScore, TheWallMaxPlayerScore)

	if p.ZombiesWon() {
		p.endGame("zombies win")
		return
	}

	if p.PlayersWon() {
		p.endGame("players win")
		return
	}
}

// endGame will end this room. We will notify each player about winners of this
// room, drop connections and stop all zombies in this room.
func (p *TheWall) endGame(reason string) {
	for _, player := range p.players {
		player.Notify("# " + reason + "\n")
		player.Drop()
	}
	if p.running {
		p.Stop()
		p.running = false
	}
	return
}

// hello will produce hello message of this room, that will be sent to player
// when new player appears.
func (p *TheWall) hello() string {
	msg := "# " + p.name + "\n"
	msg += "# Zombies are coming !!! Prepare your bows warriors !!!\n"
	return msg
}

func (p *TheWall) getPlayerScores() int64 {
	return atomic.LoadInt64(&p.playerScore)
}

func (p *TheWall) getZombieScores() int64 {
	return atomic.LoadInt64(&p.zombieScore)
}

func (p *TheWall) incZombieScores() {
	atomic.AddInt64(&p.zombieScore, 1)
}

func (p *TheWall) incPlayerScores() {
	atomic.AddInt64(&p.playerScore, 1)
}
