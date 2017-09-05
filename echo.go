package echo

import (
	"log"
	"math/rand"
	"os"
	"time"
)

// Initialize the package and random numbers, etc.
func init() {
	// Set the random seed to something different each time.
	rand.Seed(time.Now().Unix())

	// Initialize our debug logging with our prefix
	logger = log.New(os.Stdout, "[echo] ", log.Lmicroseconds)
}
