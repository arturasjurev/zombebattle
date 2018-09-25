package rooms_test

import (
	"testing"
	"time"

	"github.com/sheirys/zombebattle/engine/players"
	"github.com/sheirys/zombebattle/engine/rooms"
	"github.com/sheirys/zombebattle/engine/types"
	"github.com/sheirys/zombebattle/engine/zombies"
)

func TestTheWallShoot(t *testing.T) {

	zombie := &zombies.Crawler{}
	player := &players.MockPlayer{
		Events: make(chan types.Event),
	}

	room := &rooms.TheWall{}
	room.Run()
	room.AddZombie(zombie)
	room.AddPlayer(player)

	// move zombie to known position, so we know where to shot
	zombie.Reset(2, 5)

	// when we hit a zombie, this zombie should be respawned to pos(0,Y)
	player.ProduceEvent(types.Event{
		Type: types.EventShoot,
		X:    2,
		Y:    5,
	})

	// let the room process chan events
	time.Sleep(time.Microsecond)

	x, _ := zombie.GetPos()

	if x != rooms.TheWallMapWidth {
		t.Errorf("wrong zombie1 x position: got: %d, want: 0", x)
	}
}

func TestTheWallPlayersWin(t *testing.T) {

	zombie := &zombies.Crawler{}
	player := &players.MockPlayer{
		Events: make(chan types.Event),
	}

	room := &rooms.TheWall{}
	room.Run()
	room.AddZombie(zombie)
	room.AddPlayer(player)

	for i := 0; i < rooms.TheWallMaxPlayerScore; i++ {
		x, y := zombie.GetPos()
		player.ProduceEvent(types.Event{
			Type: types.EventShoot,
			X:    x,
			Y:    y,
		})
		time.Sleep(time.Millisecond)
	}

	if !room.PlayersWon() {
		t.Errorf("expected victory for players")
	}

	if room.ZombiesWon() {
		t.Errorf("zombies should lose")
	}
}

func TestTheWallZombiesWin(t *testing.T) {

	zombie := &zombies.Crawler{}

	room := &rooms.TheWall{}
	room.Run()
	room.AddZombie(zombie)

	for i := 0; i < rooms.TheWallMaxZombieScore; i++ {
		zombie.Reset(1, 1)
		zombie.Next()
		time.Sleep(time.Millisecond)
	}

	if !room.ZombiesWon() {
		t.Errorf("expected victory for zombies")
	}

	if room.PlayersWon() {
		t.Errorf("players should lose")
	}
}
