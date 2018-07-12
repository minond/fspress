package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/minond/fspress"
)

const (
	reloader  = `<script>setTimeout(() => location.reload(), 1000)</script>`
	indexhtml = `<ul>{{range .Posts}}<li><a href="/{{.URL}}">{{.URL}}</a></li>{{end}}</ul>`
)

var (
	blog  *fspress.Blog
	index = template.Must(template.New("index").Parse(indexhtml))

	dev        = flag.Bool("dev", false, "run blog in development mode")
	autoreload = flag.Bool("autoreload", false, "auto reload posts on an interval. requires dev mode")
	postTmpl   = flag.String("post-template", "post.tmpl", "path to post template file")
	listen     = flag.String("listen", ":8081", "host and port to listen on")
	glob       = flag.String("glob", "[0-9]*.md", "directories to check for post files")
)

func init() {
	flag.Parse()
	blog = fspress.Must(fspress.ParseGlob(*postTmpl, *glob))
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("got request to %s", r.URL.Path)

		if r.Method != http.MethodGet {
			log.Printf("ignoring %s %s request", r.Method, r.URL.Path)
			return
		}

		if *dev && r.URL.Path == "/" {
			w.Header().Add("Content-Type", "text/html")
			index.Execute(w, blog)
			return
		}

		if *dev {
			log.Println("reloading blog")
			blog = fspress.Must(fspress.ParseGlob(*postTmpl, *glob))
		}

		if post := blog.Get(r.URL.Path); post != nil {
			fmt.Fprint(w, post.String())
		}

		if *dev && *autoreload {
			fmt.Fprint(w, reloader)
		}
	})

	log.Printf("starting server on %s\n", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}
