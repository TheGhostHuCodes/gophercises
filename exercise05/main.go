package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Printf("Recieved %d arguments (%v), expected 1 argument\n", len(args), args)
		os.Exit(1)
	}

	urlForMapping := args[0]

	resp, err := http.Get(urlForMapping)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
