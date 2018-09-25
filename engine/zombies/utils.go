package zombies

import (
	"fmt"
	"math/rand"
	"time"
)

// this file contains various utils needed for zombies.

// first and last name for zombie will be generated from these strings.
var (
	fname = []string{"batman", "brain", "wombat", "leg", "eye", "night"}
	lname = []string{"drinker", "eater", "nomnom", "imbecile", "knight"}
)

func init() {
	rand.Seed(time.Now().Unix())
}

// PickName will generate random name for zombie e.g.: wombat-imbecile.
func PickName() string {
	a := fname[rand.Intn(len(fname))]
	b := lname[rand.Intn(len(lname))]
	return fmt.Sprintf("%s-%s", a, b)
}

// RandomPos will return random position for zombie. This is used for Hard
// zombie because he wants to jump like zombie rabbit.
func RandomPos(min, max int64) int64 {
	return rand.Int63n(max-min) + min
}
