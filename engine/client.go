package engine

import (
	"bufio"
	"log"
	"net"

	"github.com/sheirys/zombebattle/engine/types"
)

type Client struct {
	Name         string
	Conn         net.Conn
	eventStream  chan types.Event
	selectedRoom string
}

func (c *Client) Run() {
	c.eventStream = make(chan types.Event)
	for {
		input, err := bufio.NewReader(c.Conn).ReadBytes('\n')
		if err != nil {
			log.Printf("client disconnected")
			return
		}
		event, err := Parse(input)
		if err != nil {
			log.Printf("parse error command: %s", err)
			continue
		}

		switch event.Type {
		case types.EventStart:
			c.Name = event.Actor
		case types.EventShoot:
			event.Actor = c.Name
		}
		c.eventStream <- event
	}
}

// WaitForStart will block until client produces START event. Before that
// client can select room where he wants to join with `JOIN` command or
// create new world with `NEW` command.
func (c *Client) WaitForStart(server chan types.Event) error {
	for {
		input, err := bufio.NewReader(c.Conn).ReadBytes('\n')
		if err != nil {
			log.Printf("client disconnected")
			return err
		}
		event, err := Parse(input)
		if err != nil {
			log.Printf("parse error command: %s", err)
			continue
		}
		if event.Type == types.EventJoin {
			c.selectedRoom = event.Actor
		}
		if event.Type == types.EventNew {
			// if client wants to create a new room send this
			// command to server, so server creates new room.
			server <- event
		}
		if event.Type == types.EventStart {
			c.Name = event.Actor
			return nil
		}
	}
}

// ShowLobby will show possible rooms to client. Client should select room
// with `JOIN <room>` before starting game. If client does not select room
// then player will be forced to join to default room.
func (c *Client) ShowLobby(lobby []types.Lobby) {
	msg := "# Please select room from list above. If you do not\n"
	msg += "# select the room, after `START <name>` you will be\n"
	msg += "# forced into default room. Please select room with\n"
	msg += "# `JOIN <room>` command.\n"
	msg += "# \n"

	for _, room := range lobby {
		msg += "#    " + room.Name
		if room.Default {
			msg += " (default)"
		}
		msg += "\n"
	}

	msg += "# \n"
	msg += "# you can use `NEW <name>` to create a new world.\n"
	c.Conn.Write([]byte(msg))
}

// Notify will send cotification to client. This is used by room to print
// various information to client.
func (c *Client) Notify(msg string) {
	c.Conn.Write([]byte(msg))
}

// SelectedRoom will return room name that client wants to join.
func (c *Client) SelectedRoom() string {
	return c.selectedRoom
}

// GetEvent will pass action what client is trying to do. When we parse input
// from client, we will transform input to event, and all events will be queued
// into this channel.
func (c *Client) GetEvent() <-chan types.Event {
	return c.eventStream
}

// ProcessEvent will handle event passed by room. For example if zombie dies
// or other player is shooting or someone wins the room.
func (c *Client) ProcessEvent(e types.Event) {
	c.Conn.Write([]byte(e.String() + "\n"))
}
