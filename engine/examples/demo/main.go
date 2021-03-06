package main

import (
	"github.com/sheirys/zombebattle/engine"
	"github.com/sheirys/zombebattle/engine/rooms"
	"github.com/sheirys/zombebattle/engine/types"
	"github.com/sheirys/zombebattle/engine/zombies"
)

func main() {

	server := &engine.Server{
		Addr: ":3333",
		DefaultRoom: &rooms.TrainingGrounds{
			Zombies: []types.Zombie{
				&zombies.Dummy{},
			},
		},
		Rooms: []types.ServerRoom{
			{Room: &rooms.TheWall{}, Default: false},
		},
	}

	server.Run()

}
