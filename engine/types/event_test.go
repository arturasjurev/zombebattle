package types_test

import (
	"testing"

	"github.com/sheirys/zombebattle/engine/types"
)

func TestEventToString(t *testing.T) {
	testTable := []struct {
		Event          types.Event
		ExpectedString string
	}{
		{
			Event: types.Event{
				Type:  types.EventWalk,
				Actor: "zombie",
				X:     1,
				Y:     2,
			},
			ExpectedString: "WALK zombie 1 2",
		},
		{
			Event: types.Event{
				Type:  types.EventShoot,
				Actor: "player",
				X:     1,
				Y:     2,
			},
			ExpectedString: "SHOOT 1 2",
		},
		{
			Event: types.Event{
				Type:   types.EventBoom,
				Actor:  "player",
				Points: 1,
				Hits:   []string{},
			},
			ExpectedString: "BOOM player 1 []",
		},
	}

	for idx, c := range testTable {
		str := c.Event.String()
		if c.ExpectedString != str {
			t.Errorf("incorrect event format: case %d, got: '%s', want: '%s'", idx, str, c.ExpectedString)
		}
	}
}
