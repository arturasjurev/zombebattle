package engine

import (
	"bufio"
	"log"
	"net"
)

type Client struct {
	Name        string
	Conn        net.Conn
	eventStream chan Event
}

func (c *Client) Run() {
	c.eventStream = make(chan Event)
	for {
		input, err := bufio.NewReader(c.Conn).ReadBytes('\n')
		if err != nil {
			log.Printf("parse error: %s", err)
			return
		}
		event, err := Parse(input)
		if err != nil {
			log.Printf("parse error command: %s", err)
			continue
		}

		switch event.Type {
		case EventStart:
			c.Name = event.Actor
		case EventShoot:
			event.Actor = c.Name
		}
		c.eventStream <- event
	}
}

func (c *Client) Identify() {
	for {
		input, err := bufio.NewReader(c.Conn).ReadBytes('\n')
		if err != nil {
			log.Printf("parse error: %s", err)
			return
		}
		event, err := Parse(input)
		if err != nil {
			log.Printf("parse error command: %s", err)
			continue
		}
		if event.Type == EventStart {
			c.Name = event.Actor
			return
		}
	}

}

func (c *Client) AskRoom() string {
	return ""
}

func (c *Client) Kill() {}

func (c *Client) GetEvent() <-chan Event {
	return c.eventStream
}

func (c *Client) ProcessEvent(e Event) {
	c.Conn.Write([]byte(e.String() + "\n"))
}
