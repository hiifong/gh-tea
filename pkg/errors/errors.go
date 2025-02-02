package errors

import (
	"github.com/charmbracelet/log"
)

func Check(msg interface{}) {
	if msg != nil {
		log.Fatalf("Error: %v", msg)
	}
}
