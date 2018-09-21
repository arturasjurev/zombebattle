package engine

import (
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var killSignals = []os.Signal{syscall.SIGTERM, os.Interrupt, os.Kill}

type ServerRoom struct {
	Room    Room
	Default bool
}

type Server struct {
	Addr        string
	DefaultRoom Room

	newClient chan net.Conn
	stop      chan os.Signal
	command   chan Event
	roomsMtx  sync.Mutex
	Rooms     []ServerRoom
}

func (s *Server) init() error {
	s.newClient = make(chan net.Conn)
	s.stop = make(chan os.Signal)
	s.command = make(chan Event)

	signal.Notify(s.stop, killSignals...)

	if s.DefaultRoom != nil {
		s.AddRoom(ServerRoom{
			Room:    s.DefaultRoom,
			Default: true,
		})
	}

	return s.listen()
}

func (s *Server) Run() {
	s.init()
	s.startRooms()
	working := true
	for working {
		select {
		case connection := <-s.newClient:
			go s.acceptClient(connection)
		case command := <-s.command:
			s.acceptCommand(command)
		case <-s.stop:
			s.Shutdown()
			working = false
		}
	}
}

func (s *Server) Shutdown() {}

func (s *Server) AddRoom(r ServerRoom) {
	s.roomsMtx.Lock()
	s.Rooms = append(s.Rooms, r)
	s.roomsMtx.Unlock()
}

func (s *Server) startRooms() {
	s.roomsMtx.Lock()
	for _, r := range s.Rooms {
		r.Room.Run()
	}
	s.roomsMtx.Unlock()
}

func (s *Server) acceptClient(c net.Conn) {
	client := &Client{
		Name: "unknown warrior",
		Conn: c,
	}
	client.Identify()

	log.Printf("client identified as %s", client.Name)
	if s.DefaultRoom != nil {
		log.Printf("client %s force-joined to world '%s'",
			client.Name,
			s.DefaultRoom.Name(),
		)
		s.DefaultRoom.AddPlayer(client)
		go client.Run()

	}

	asked := client.AskRoom()
	s.roomsMtx.Lock()
	for _, r := range s.Rooms {
		if r.Room.Name() == asked {
			r.Room.AddPlayer(client)
			break
		}
	}
	s.roomsMtx.Unlock()
}

func (s *Server) acceptCommand(c Event) {}

func (s *Server) listen() error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	log.Printf("listening on %s", s.Addr)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("cannot accept connection: %s", err)
				continue
			}
			log.Printf("accepted connection from %s", conn.RemoteAddr())
			s.newClient <- conn
		}
	}()
	return nil
}

func (s *Server) lobby() (lobby []string) {
	s.roomsMtx.Lock()
	for _, r := range s.Rooms {
		lobby = append(lobby, r.Room.Name())
	}
	s.roomsMtx.Unlock()
	return
}
