package main

import (
	"fmt"
	"log"
	"os"

	"github.com/TheGhostHuCodes/gophercises/exercise04/link"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Printf("Recieved %d arguments (%v), expected 1 argument\n", len(args), args)
		os.Exit(1)
	}

	filename := args[0]
	r, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Unable to open file %s, received error: %v", args[0], err)
	}

	links, err := link.Parse(r)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", links)
}
