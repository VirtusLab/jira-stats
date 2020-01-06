package analyzer

import (
	"log"
	"time"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("TIMING [%s] took %s", name, elapsed)
}
