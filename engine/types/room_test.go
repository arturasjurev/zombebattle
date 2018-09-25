package types_test

import (
	"reflect"
	"testing"

	"github.com/sheirys/zombebattle/engine/rooms"
	"github.com/sheirys/zombebattle/engine/types"
)

// TestRoomInterface will test all rooms that satisfies this interface. As all
// rooms should implement Room interface we expect that they will behave same
// with same arguments.
func TestRoomInterface(t *testing.T) {
	testWithInterface(t, &rooms.TrainingGrounds{})
}

func testWithInterface(t *testing.T, i types.Room) {
	// functions that should be tested with this interface implementation.
	testFunc := []func(*testing.T, types.Room){
		testName,
	}

	// call all testable functions.
	for _, fn := range testFunc {
		t.Logf("using (%s) as interface\n", reflect.TypeOf(i).String())
		fn(t, i)
	}
}

func testName(t *testing.T, i types.Room) {
	testTable := []struct {
		SetName, ExpectedName string
	}{
		{"test1", "test1"},
		{"name", "name"},
	}

	for idx, c := range testTable {
		i.SetName(c.SetName)
		if c.ExpectedName != i.Name() {
			t.Errorf("wrong name: case %d, got:'%s', want:'%s'", idx, i.Name(), c.ExpectedName)
		}
	}
}
