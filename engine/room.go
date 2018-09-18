package engine

// Scoreboard is used to hold room scores
type Scoreboard map[string]uint

// Room defines what we expect from room. Please notice that we do not provice
// map size here. Map size is not a requirement here and if you want you can
// implement unlimited size map here. Also how scores are calculated depends
// on implementation.
type Room interface {

	// Init should be called before room usage. Here various preparations
	// should be applied.
	Init() error

	// Stop should stop the room and kill everyone inside it. This is like
	// end of the world.
	Stop() error

	// AddZombie should add zompie into room.
	AddZombie(z Zombie) error

	// AddEvent should add new event into room. Here user defined events
	// should be pushed. Zombie events will be pushed internaly and you
	// dont neet to worry about zombies here.
	AddEvent(e Event)

	// EventStream should proxy all room events here. It is used to log
	// events in console and send to users.
	EventStream() <-chan Event

	// Scores should return room scores.
	Scores() Scoreboard
}
