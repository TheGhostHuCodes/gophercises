package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/TheGhostHuCodes/gophercises/exercise02/urlshort"
	bolt "github.com/coreos/bbolt"
)

func main() {
	yamlFilename := flag.String("yaml", "", "A YAML file containing path to url mappings.")
	flag.Parse()
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	var yaml []byte
	if *yamlFilename != "" {
		file, err := os.Open(*yamlFilename)
		if err != nil {
			log.Fatalf("Could not open YAML file: '%v'", err)
		}
		yaml, err = ioutil.ReadAll(file)
		if err != nil {
			log.Fatalf("Could not read YAML file: '%v'", err)
		}
	} else {
		yaml = []byte(`
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`)
	}
	yamlHandler, err := urlshort.YAMLHandler(yaml, mapHandler)
	if err != nil {
		panic(err)
	}

	json := []byte(`
[
	{"path": "/hackernews", "url": "https://news.ycombinator.com/"},
	{"path": "/arstech", "url": "https://arstechnica.com/"},
	{"path": "/xkcd", "url": "https://xkcd.com/"}
]
`)
	jsonHandler, err := urlshort.JSONHandler(json, yamlHandler)
	if err != nil {
		panic(err)
	}

	db, err := bolt.Open("shortener.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	defer db.Close()
	if err != nil {
		panic(fmt.Errorf("Error opening database: %s", err))
	}

	err = populateDb(db)
	if err != nil {
		panic(err)
	}

	boltHandler, err := urlshort.BoltHandler(db, jsonHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", boltHandler)
}

func populateDb(db *bolt.DB) error {
	// Create "paths" bucket.
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(urlshort.BoltDBPathsBucketName))
		if err != nil {
			return fmt.Errorf("Error creating bucket : %s", err)
		}
		return nil
	})

	// Populate "paths" bucket.
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(urlshort.BoltDBPathsBucketName))
		err := b.Put([]byte("/boltdb"), []byte("https://github.com/coreos/bbolt"))
		return err
	})
	if err != nil {
		return fmt.Errorf("Error putting into bucket: %s", err)
	}

	return nil
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
