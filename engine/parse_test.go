package engine_test

import (
	"testing"

	"github.com/sheirys/zombebattle/engine"
	"github.com/sheirys/zombebattle/engine/types"
)

func TestParse(t *testing.T) {
	testTable := []struct {
		Input         []byte
		ExpectedEvent types.Event
		ExpectedErr   error
	}{
		{
			Input: []byte("start jonas"),
			ExpectedEvent: types.Event{
				Type:  types.EventStart,
				Actor: "JONAS",
			},
			ExpectedErr: nil,
		},
		{
			Input: []byte("shoot 1 2"),
			ExpectedEvent: types.Event{
				Type: types.EventShoot,
				X:    1,
				Y:    2,
			},
			ExpectedErr: nil,
		},
		{
			Input:         []byte("shoot a 2b"),
			ExpectedEvent: types.Event{},
			ExpectedErr:   engine.ErrBadInput,
		},
		{
			Input:         []byte("shoot 1 2b"),
			ExpectedEvent: types.Event{},
			ExpectedErr:   engine.ErrBadInput,
		},
		{
			Input: []byte("join castle1"),
			ExpectedEvent: types.Event{
				Type:  types.EventJoin,
				Actor: "CASTLE1",
			},
			ExpectedErr: nil,
		},
		{
			Input: []byte("new castle1"),
			ExpectedEvent: types.Event{
				Type:  types.EventNew,
				Actor: "CASTLE1",
			},
			ExpectedErr: nil,
		},
		{
			Input:         []byte("fat mama"),
			ExpectedEvent: types.Event{},
			ExpectedErr:   engine.ErrBadInput,
		},
		{
			Input:         []byte("jashgkjhdlkfjhaluefshalusf"),
			ExpectedEvent: types.Event{},
			ExpectedErr:   engine.ErrBadInput,
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
