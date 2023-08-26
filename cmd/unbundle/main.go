package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s: unbundle [options] <bundle file path>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	args := os.Args[1:]

	if len(args) != 1 {
		log.Fatalf("Invalid number of arguments. Expected a FHIR bundle file path.")
	}

	err := Unbundle(args[0])
	if err != nil {
		log.Fatalf("Failed to unbundle FHIR bundle: %s", err)
	}
}
