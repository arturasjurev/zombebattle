package main

import (
	"bufio"
	"log"
	"net"

	"github.com/sheirys/zombebattle/engine"
	"github.com/sheirys/zombebattle/engine/communication"
	"github.com/sheirys/zombebattle/engine/rooms"
	"github.com/sheirys/zombebattle/engine/zombies"
)

func main() {

	l, err := net.Listen("tcp", ":3333")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		log.Printf("server started on :3333")
		conn, err := l.Accept()
		if err != nil {
			log.Printf("error while accepting: %s", err)
			continue
		}
		go startSession(conn)
	}
}

func startSession(conn net.Conn) {
	log.Printf("connected: %s", conn.RemoteAddr().String())
	log.Printf("creating new room for new client")
	conn.Write([]byte("# hello unknown warrior!\n"))
	conn.Write([]byte("# use `START <name>` to start adventure >:o\n"))

	room := rooms.PlayRoom{}
	room.Init()

	go userEventsToRoom(conn, &room)

	room.AddZombie(&zombies.Easy{})
	room.AddZombie(&zombies.Hard{})

	for {
		// FIXME: this fanout method doeas not work with multiple users.
		event := <-room.EventStream()
		conn.Write([]byte(event.String() + "\n"))
		log.Println(event.String())
	}

}

func userEventsToRoom(c net.Conn, r engine.Room) {
	// FIXME: this is hack for now. this logic should be moved to sessionHandler
	// or something similar.
	username := ""
	for {
		input, err := bufio.NewReader(c).ReadBytes('\n')
		if err != nil {
			log.Printf("parse error: %s", err)
			continue

		}
		event, err := communication.Parse(input)
		if err != nil {
			log.Printf("parse error: %s", err)
			continue
		}

		switch event.Type {
		case engine.EventStart:
			username = event.Actor
			r.AddEvent(event)
		case engine.EventShoot:
			event.Actor = username
			r.AddEvent(event)
		}
	}
}
