package rooms_test

import (
	"testing"
	"time"

	"github.com/sheirys/zombebattle/engine/players"
	"github.com/sheirys/zombebattle/engine/rooms"
	"github.com/sheirys/zombebattle/engine/types"
	"github.com/sheirys/zombebattle/engine/zombies"
)

func TestTrainingGrounds(t *testing.T) {

	zombie := &zombies.Dummy{}
	player := &players.MockPlayer{
		Events: make(chan types.Event),
	}

	room := &rooms.TrainingGrounds{}
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

	if room.PlayersWon() {
		t.Errorf("should be impossible to win TrainingGrounds for players")
	}

	if room.ZombiesWon() {
		t.Errorf("should be impossible to win TrainingGrounds for zombies")
	}
}
