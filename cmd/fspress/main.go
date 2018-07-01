package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/minond/fspress"
)

const reloader = "<script>setTimeout(() => location.reload(), 1000)</script>"

var (
	dev      = flag.Bool("dev", false, "Run blog in development mode")
	postTmpl = flag.String("post-template", "post.tmpl", "Path to post template file")
	listen   = flag.String("listen", ":8081", "Host and port to listen on")
	glob     = flag.String("glob", "[0-9]*.md", "Directories to check for post files")
)

func main() {
	flag.Parse()

	blog := newBlog()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("got request to %s", r.URL.Path)

		if r.Method != http.MethodGet {
			log.Printf("ignoring %s %s request", r.Method, r.URL.Path)
			return
		}

		if *dev {
			log.Println("reloading blog")
			blog = newBlog()
		}

		if post := blog.Get(r.URL.Path); post != nil {
			fmt.Fprint(w, post.String())
		}

		if *dev {
			fmt.Fprint(w, reloader)
		}
	})

	log.Printf("starting server on %s\n", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}

func newBlog() *fspress.Blog {
	log.Printf("building blog from %s\n", *glob)
	files, err := fspress.FindPostFiles(*glob)
	if err != nil {
		log.Fatalf("error finding post files: %v\n", err)
	}

	blog, err := fspress.New(*postTmpl, files)
	if err != nil {
		log.Fatalf("error generating blog: %v\n", err)
	}

	for _, path := range blog.Paths() {
		log.Printf("serving %s on /%s\n", path, fspress.CleanURL(path))
	}

	return blog
}
