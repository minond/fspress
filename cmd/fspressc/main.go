package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/minond/fspress"
)

var (
	blog *fspress.Blog

	catalog  = flag.String("catalog", "catalog.csv", "path to catalog csv file")
	glob     = flag.String("glob", "[0-9]*.md", "directories to check for post files")
	out      = flag.String("out", ".", "output directory")
	postTmpl = flag.String("post-template", "post.tmpl", "path to post template file")
)

func init() {
	flag.Parse()
	blog = fspress.New(*catalog, *postTmpl, *glob)
	if err := blog.Load(); err != nil {
		panic(err)
	}
}

func main() {
	log.Print("compiling posts")
	for _, post := range blog.Posts {
		log.Printf("saved %s to %s/%s.html", post.Path, *out, post.URL)
		ioutil.WriteFile(*out+"/"+post.URL+".html", []byte(post.String()), 0644)
	}
	log.Print("done")
}
