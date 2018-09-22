package engine_test

import (
	"testing"

	"github.com/sheirys/zombebattle/engine"
)

func TestParse(t *testing.T) {
	testTable := []struct {
		Input         []byte
		ExpectedEvent engine.Event
		ExpectedErr   error
	}{
		{
			Input: []byte("start jonas"),
			ExpectedEvent: engine.Event{
				Type:  engine.EventStart,
				Actor: "JONAS",
			},
			ExpectedErr: nil,
		},
		{
			Input: []byte("shoot 1 2"),
			ExpectedEvent: engine.Event{
				Type: engine.EventShoot,
				X:    1,
				Y:    2,
			},
			ExpectedErr: nil,
		},
		{
			Input: []byte("join castle1"),
			ExpectedEvent: engine.Event{
				Type:  engine.EventJoin,
				Actor: "CASTLE1",
			},
			ExpectedErr: nil,
		},
	}

	for i, v := range testTable {
		event, err := engine.Parse(v.Input)
		if event.Type != v.ExpectedEvent.Type {
			t.Errorf("incorrect event type. case: %d got: %s want: %s", i, event.Type, v.ExpectedEvent.Type)
		}
		if event.X != v.ExpectedEvent.X {
			t.Errorf("incorrect event x. case: %d got: %d want: %d", i, event.X, v.ExpectedEvent.X)
		}
		if event.Y != v.ExpectedEvent.Y {
			t.Errorf("incorrect event y. case: %d got: %d want: %d", i, event.Y, v.ExpectedEvent.Y)
		}
		if event.Actor != v.ExpectedEvent.Actor {
			t.Errorf("incorrect event actor. case: %d got: %s want: %s", i, event.Actor, v.ExpectedEvent.Actor)
		}
		if err != v.ExpectedErr {
			t.Errorf("incorrect error. case: %d got: %v want: %v", i, err, v.ExpectedErr)
		}
	}
}
