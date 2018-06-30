package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var (
	listen = flag.String("listen", ":8081", "Host and port to listen on")
	source = flag.String("source", ".", "Directories to check for post files")
)

func main() {
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			log.Printf("ignoring %s %s request", r.Method, r.URL.Path)
			return
		}

		name := strings.TrimLeft(strings.TrimRight(r.URL.Path, ".html"), "/")

		log.Printf("got request to %s", r.URL.Path)
		log.Printf("serving %s", name)

		fmt.Fprint(w, "hi")
	})

	log.Fatal(http.ListenAndServe(*listen, nil))
}
