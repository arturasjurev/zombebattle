package engine

import (
	"errors"
	"strconv"
	"strings"

	"github.com/sheirys/zombebattle/engine/types"
)

var (
	// ErrBadInput will be returned when command cannot be identified or
	// parsed into types.Event
	ErrBadInput = errors.New("bad command, won't parse")
)

// Parse will parse player input and produces Event from that input. See
// event.go for more information about events.
func Parse(b []byte) (types.Event, error) {
	s := string(b)
	s = strings.Trim(s, "\n")
	s = strings.Trim(s, "\r")
	s = strings.Trim(s, " ")

	// expect that command has at last 2 strings.
	args := strings.Split(s, " ")
	if len(args) < 2 {
		return types.Event{}, ErrBadInput
	}

	// commands should be case-insensitive
	for i, v := range args {
		args[i] = strings.ToUpper(v)
	}

	switch {
	// parse SHOOT command e.g.: SHOOT 1 2
	case args[0] == types.EventShoot && len(args) == 3:
		return parseShoot(args)
	// parse START command e.g.: START ivan
	case args[0] == types.EventStart && len(args) == 2:
		return parseStart(args)
	// parse JOIN command e.g.: JOIN woodstock
	case args[0] == types.EventJoin && len(args) == 2:
		return parseJoin(args)
	// parse NEW command e.g.: NEW world2
	case args[0] == types.EventNew && len(args) == 2:
		return parseNew(args)
	default:
		return types.Event{}, ErrBadInput
	}
}

// parseStart will parse START command and produce EventStart event. Here client
// name will be stored as Actor.
func parseStart(cmd []string) (types.Event, error) {
	return types.Event{
		Type:  types.EventStart,
		Actor: cmd[1],
	}, nil
}

// parseJoin will parse JOIN comand and produce EventJoin event. Here requested
// room will be stored as Actor.
func parseJoin(cmd []string) (types.Event, error) {
	return types.Event{
		Type:  types.EventJoin,
		Actor: cmd[1],
	}, nil
}

// parseNew will parse NEW command and produce EventNew event. Here requested
// room name will be stored as Actor.
func parseNew(cmd []string) (types.Event, error) {
	return types.Event{
		Type:  types.EventNew,
		Actor: cmd[1],
	}, nil
}

// parseShoot will parse SHOOT command and produce EventShoot event. As we
// cannot know from input players name, we cannot store it as Actor. So
// this event should be appended with player name as Actor latter in room or
// client struct.
func parseShoot(cmd []string) (types.Event, error) {

	// parse X from command
	x, err := strconv.Atoi(cmd[1])
	if err != nil {
		return types.Event{}, ErrBadInput
	}

	// parse Y from command
	y, err := strconv.Atoi(cmd[2])
	if err != nil {
		return types.Event{}, ErrBadInput
	}

	event := types.Event{
		Type: types.EventShoot,
		X:    int64(x),
		Y:    int64(y),
		// FIXME: this is bad. We reach player name only in room, so
		// for demo, we will append this event with player name in room.
		//Player: "player",
	}
	return event, nil
}
