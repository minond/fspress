package fspress

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	blackfriday "gopkg.in/russross/blackfriday.v2"
)

type Blog struct {
	posts map[string]*post
}

type post struct {
	Path    string
	Content string
	tmpl    *template.Template
}

func Must(blog *Blog, err error) *Blog {
	if err != nil {
		panic(err)
	}

	return blog
}

func ParseFiles(postTmpl string, files []string) (*Blog, error) {
	tmpl := template.Must(template.ParseFiles(postTmpl))
	blog := &Blog{posts: make(map[string]*post)}

	for _, file := range files {
		name := cleanURL(file)
		post := &post{Path: file, tmpl: tmpl}

		if err := post.Load(); err != nil {
			return nil, err
		}

		blog.posts[name] = post
	}

	return blog, nil
}

func ParseGlob(postTmpl string, glob string) (*Blog, error) {
	files, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}

	return ParseFiles(postTmpl, files)
}

func (b *Blog) Get(file string) *post {
	return b.posts[cleanURL(file)]
}

func (p *post) Load() error {
	bytes, err := ioutil.ReadFile(p.Path)
	if err != nil {
		return err
	}

	p.Content = string(blackfriday.Run(bytes))
	return nil
}

func (p *post) String() string {
	buf := &bytes.Buffer{}
	p.tmpl.Execute(buf, p)
	return buf.String()
}

func cleanURL(name string) string {
	re := regexp.MustCompile("^[0-9]+-")
	base := strings.TrimLeft(
		strings.TrimRight(strings.TrimRight(name, ".html"), ".md"), "/")
	return re.ReplaceAllString(base, "")
}
