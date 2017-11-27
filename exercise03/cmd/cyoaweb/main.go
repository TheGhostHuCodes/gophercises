package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/TheGhostHuCodes/gophercises/exercise03"
)

func main() {
	port := flag.Int("port", 3000, "The port to start the CYOA web applicaiton on.")
	filename := flag.String("file", "gopher.json", "The JSON file with the CYOA story.")
	flag.Parse()
	fmt.Printf("Using the story in %s.\n", *filename)

	f, err := os.Open(*filename)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	story, err := cyoa.JSONStory(f)
	if err != nil {
		panic(err)
	}

	h := cyoa.NewHandler(story, nil)
	fmt.Printf("Starting the server on port : %d\n", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}
