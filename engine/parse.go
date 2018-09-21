package engine

import (
	"strconv"
	"strings"
)

// Parse satisfies CommunicationChannel interface and can be used as player
// input parser. This version does not support extended commands set.
func Parse(b []byte) (Event, error) {
	s := string(b)
	s = strings.Trim(s, "\n")
	s = strings.Trim(s, "\r")
	s = strings.Trim(s, " ")

	// expect that command has at last 2 strings.
	args := strings.Split(s, " ")
	if len(args) < 2 {
		return Event{}, ErrBadInput
	}

	// commands should be case-insensitive
	for i, v := range args {
		args[i] = strings.ToUpper(v)
	}

	switch {
	// parse SHOOT command e.g.: SHOOT 1 2
	case args[0] == EventShoot && len(args) == 3:
		return parseShoot(args)
	// parse START command e.g.: START ivan
	case args[0] == EventStart && len(args) == 2:
		return parseStart(args)
	default:
		return Event{}, ErrBadInput
	}

	// FIXME: wtf, if we reach this then this is a bug.
	// return engine.Event{}, engine.ErrBadInput
}

func parseShoot(cmd []string) (Event, error) {

	// parse X from command
	x, err := strconv.Atoi(cmd[1])
	if err != nil {
		return Event{}, err
	}

	// parse Y from command
	y, err := strconv.Atoi(cmd[2])
	if err != nil {
		return Event{}, err
	}

	event := Event{
		Type: EventShoot,
		X:    x,
		Y:    y,
		// FIXME: this is bad. We reach player name only in room, so
		// for demo, we will append this event with player name in room.
		//Player: "player",
	}
	return event, nil
}

func parseStart(cmd []string) (Event, error) {
	return Event{
		Type:  EventStart,
		Actor: cmd[1],
	}, nil
}
