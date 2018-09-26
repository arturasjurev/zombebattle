package zombies_test

import (
	"context"
	"testing"

	"github.com/sheirys/zombebattle/engine/types"
	"github.com/sheirys/zombebattle/engine/zombies"
)

func TestCrawler(t *testing.T) {
	events := make(chan types.Event, 1)
	ctx := context.Background()
	event := types.Event{}

	crawler := &zombies.Crawler{}
	crawler.Summon(ctx, events)
	crawler.Run()

	crawler.Reset(5, 5)
	crawler.Next()
	event = <-events

	if event.Type != types.EventWalk {
		t.Errorf("unexpected event. got: '%s', want: '%s'", event.Type, types.EventWalk)
	}

	if event.X != 4 {
		t.Errorf("unexpected x position. got: %d, want: %d", event.X, 5)
	}

	if event.Y != 5 {
		t.Errorf("unexpected x position. got: %d, want: %d", event.Y, 5)
	}

	x, y := crawler.GetPos()

	if x != 4 {
		t.Errorf("unexpected x position. got: %d, want: %d", event.X, 5)
	}

	if y != 5 {
		t.Errorf("unexpected x position. got: %d, want: %d", event.Y, 5)
	}

	if name := crawler.GetName(); name == "" {
		t.Errorf("dummy does not have name. expected name.")
	}

	if !crawler.Hit() {
		t.Errorf("crawler should be dead now")
	}
}
