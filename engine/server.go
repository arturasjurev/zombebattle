package engine

import (
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sheirys/zombebattle/engine/rooms"
	"github.com/sheirys/zombebattle/engine/types"
)

var killSignals = []os.Signal{syscall.SIGTERM, os.Interrupt, os.Kill}

// Server holds information about game server.
type Server struct {
	Addr        string
	DefaultRoom types.Room
	Rooms       []types.ServerRoom
	newClient   chan net.Conn
	stop        chan os.Signal
	command     chan types.Event
	roomsMtx    sync.Mutex
}

// Run starts to listen for events and handle them.
func (s *Server) Run() {
	s.init()
	s.startRooms()
	for {
		select {
		case connection := <-s.newClient:
			go s.acceptClient(connection)
		case command := <-s.command:
			// FIXME: For now only one EventNew is supported.
			// FIXME: New maps will always be rooms.TheWall rooms.
			s.createRoom(command.Actor, command.X)
		case <-s.stop:
			s.Shutdown()
			return
		}
	}
}

// Shutdown will stop server
func (s *Server) Shutdown() {}

// AddRoom will registers new room into server.
func (s *Server) AddRoom(r types.ServerRoom) {
	s.roomsMtx.Lock()
	s.Rooms = append(s.Rooms, r)
	s.roomsMtx.Unlock()
}

func (s *Server) init() error {
	s.newClient = make(chan net.Conn)
	s.stop = make(chan os.Signal)
	s.command = make(chan types.Event)

	signal.Notify(s.stop, killSignals...)

	// move default room to the lobby so it will be visible for players.
	if s.DefaultRoom != nil {
		s.AddRoom(types.ServerRoom{
			Room:    s.DefaultRoom,
			Default: true,
		})
	}

	return s.listen()
}

func (s *Server) startRooms() {
	s.roomsMtx.Lock()
	for _, r := range s.Rooms {
		r.Room.Run()
	}
	s.roomsMtx.Unlock()
}

func (s *Server) createRoom(name string, world int64) {
	log.Printf("creating new room")
	room := &rooms.TrainingGrounds{}
	room.SetName(name)
	s.AddRoom(types.ServerRoom{
		Room:    room,
		Default: false,
	})
}

// acceptClient will be called when new connection appears in server.
func (s *Server) acceptClient(c net.Conn) {
	client := &Client{
		Name: "unknown warrior",
		Conn: c,
	}

	// show possible rooms to client. Client can select where he wants to
	// join with `JOIN` command.
	client.ShowLobby(s.lobby())

	// wait until client produces EventStart. Also client can select room
	// where he wants to join or even create new room with `NEW` command.
	// So `JOIN`, `NEW` and `START` commands will be processed here.
	if err := client.WaitForStart(s.command); err != nil {
		log.Printf("WaitForStart returned error: %s", err)
		return
	}

	// we expect that client selected room with JOIN command. If no, then
	// client will be forced to join to default room.
	if s.DefaultRoom != nil && client.SelectedRoom() == "" {
		log.Printf("force joined")
		s.DefaultRoom.AddPlayer(client)
		go client.Run()
		return
	}

	// join client to required room.
	s.roomsMtx.Lock()
	for _, r := range s.Rooms {
		if r.Room.Name() == client.SelectedRoom() {
			r.Room.AddPlayer(client)
			break
		}
	}
	s.roomsMtx.Unlock()
}

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

func (s *Server) lobby() (lobby []types.Lobby) {
	s.roomsMtx.Lock()
	for _, r := range s.Rooms {
		lobby = append(lobby, types.Lobby{
			Name:    r.Room.Name(),
			Default: r.Default,
		})
	}
	s.roomsMtx.Unlock()
	return
}
