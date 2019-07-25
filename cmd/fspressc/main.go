package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/minond/fspress"
)

var (
	blog *fspress.Blog

	out         = flag.String("out", ".", "output directory")
	postCatalog = flag.String("post-catalog", "catalog.csv", "path to catalog csv file")
	postGlob    = flag.String("post-glob", "[0-9]*.md", "directories to check for post files")
	postTmpl    = flag.String("post-template", "post.tmpl", "path to post template file")
)

func init() {
	flag.Parse()
	blog = fspress.New(*postCatalog, *postTmpl, *postGlob)
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
