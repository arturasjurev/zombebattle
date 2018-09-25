package zombies_test

import (
	"context"
	"testing"

	"github.com/sheirys/zombebattle/engine/types"
	"github.com/sheirys/zombebattle/engine/zombies"
)

func TestDummy(t *testing.T) {
	events := make(chan types.Event, 1)
	ctx := context.Background()
	event := types.Event{}

	dummy := &zombies.Dummy{}
	dummy.Summon(ctx, events)
	dummy.Run()

	dummy.Reset(5, 5)
	dummy.Next()
	event = <-events

	if event.Type != types.EventWalk {
		t.Errorf("unexpected event. got: '%s', want: '%s'", event.Type, types.EventWalk)
	}

	if event.X != 5 {
		t.Errorf("unexpected x position. got: %d, want: %d", event.X, 5)
	}

	if event.Y != 5 {
		t.Errorf("unexpected x position. got: %d, want: %d", event.Y, 5)
	}
}
