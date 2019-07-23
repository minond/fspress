package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/minond/fspress"
)

const (
	reloader  = `<script>setTimeout(() => location.reload(), 1000)</script>`
	indexhtml = `<ul>{{range .Posts}}<li><a href="/{{.URL}}">{{.URL}}</a></li>{{end}}</ul>`
)

var (
	blog  *fspress.Blog
	index = template.Must(template.New("index").Parse(indexhtml))

	autoreload = flag.Bool("autoreload", false, "auto reload posts on an interval. requires dev mode")
	catalog    = flag.String("catalog", "catalog.csv", "path to catalog csv file")
	dev        = flag.Bool("dev", false, "run blog in development mode")
	glob       = flag.String("glob", "[0-9]*.md", "directories to check for post files")
	listen     = flag.String("listen", ":8081", "host and port to listen on")
	postTmpl   = flag.String("post-template", "post.tmpl", "path to post template file")
)

func init() {
	flag.Parse()
	blog = fspress.New(*catalog, *postTmpl, *glob)
	if err := blog.Load(); err != nil {
		panic(err)
	}
}

func fileExists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func main() {
	extra := ""
	static := http.FileServer(http.Dir("."))

	if *dev && *autoreload {
		extra += reloader
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("got request to %s", r.URL.Path)

		if r.Method != http.MethodGet {
			log.Printf("ignoring %s %s request", r.Method, r.URL.Path)
			return
		}

		path := r.URL.Path
		if r.URL.Path == "/" {
			path = "index.html"
		}

		if fileExists(path) {
			static.ServeHTTP(w, r)
			return
		}

		if *dev && r.URL.Path == "/" {
			w.Header().Add("Content-Type", "text/html")
			index.Execute(w, blog)
			return
		}

		if *dev {
			log.Println("reloading blog")
			if err := blog.Load(); err != nil {
				panic(err)
			}
		}

		if post := blog.Get(r.URL.Path); post != nil {
			fmt.Fprint(w, post.String()+extra)
		} else {
			static.ServeHTTP(w, r)
		}
	})

	log.Printf("starting server on %s\n", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}
