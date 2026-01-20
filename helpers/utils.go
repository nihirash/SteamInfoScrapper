package helpers

import (
	"log"
)

// ProcessError logs an error without terminating the process.
func ProcessError(err error) {
	log.Printf("Error happens: %s", err.Error())
}
