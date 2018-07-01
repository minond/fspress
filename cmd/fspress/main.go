package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/minond/fspress"
)

const (
	reloader = "<script>setTimeout(() => location.reload(), 1000)</script>"
)

var (
	blog *fspress.Blog

	dev      = flag.Bool("dev", false, "Run blog in development mode")
	out      = flag.String("out", "", "Output directory")
	postTmpl = flag.String("post-template", "post.tmpl", "Path to post template file")
	listen   = flag.String("listen", ":8081", "Host and port to listen on")
	glob     = flag.String("glob", "[0-9]*.md", "Directories to check for post files")
)

func init() {
	flag.Parse()
	blog = fspress.Must(fspress.ParseGlob(*postTmpl, *glob))
}

func main() {
	if *out != "" {
		log.Print("compiling posts")
		for _, post := range blog.Posts() {
			log.Printf("saved %s to %s/%s.html", post.Path, *out, post.URL)
			ioutil.WriteFile(*out+"/"+post.URL+".html", []byte(post.Content), 0644)
		}
		log.Print("done")
	} else {
		log.Print("starting server")
		server()
	}
}

func server() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("got request to %s", r.URL.Path)

		if r.Method != http.MethodGet {
			log.Printf("ignoring %s %s request", r.Method, r.URL.Path)
			return
		}

		if *dev {
			log.Println("reloading blog")
			blog = fspress.Must(fspress.ParseGlob(*postTmpl, *glob))
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
