package main

import (
	"github.com/sheirys/zombebattle/engine"
	"github.com/sheirys/zombebattle/engine/rooms"
)

func main() {

	server := &engine.Server{
		Addr:        ":3333",
		DefaultRoom: &rooms.TheWall{},
	}

	server.Run()

}
