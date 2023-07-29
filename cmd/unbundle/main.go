package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		log.Fatalf("Invalid number of arguments. Expected a FHIR bundle file path.")
	}

	err := Unbundle(args[0])
	if err != nil {
		log.Fatalf("Failed to unbundle FHIR bundle: %s", err)
	}
}
