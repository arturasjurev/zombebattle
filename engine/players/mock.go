package players

import (
	"log"

	"github.com/sheirys/zombebattle/engine/types"
)

// MockPlayer satisfies engine.player interface and can be used in tests.
type MockPlayer struct {
	Events chan types.Event
}

// Notify will print notfy message for client - in this case in log console.
func (m *MockPlayer) Notify(msg string) {
	log.Printf(msg + "\n")
}

// GetEvent will return event produced by this mock client.
func (m *MockPlayer) GetEvent() (types.Event, bool) {
	event, ok := <-m.Events
	return event, ok
}

// ProcessEvent will handle event.
// FIXME: implement this.
func (m *MockPlayer) ProcessEvent(e types.Event) {}

// ProduceEvent will add event into mock client event stream.
func (m *MockPlayer) ProduceEvent(e types.Event) {
	m.Events <- e
}
