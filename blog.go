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
	posts    map[string]*post
	postTmpl *template.Template
}

type post struct {
	path    string
	Content string
	tmpl    *template.Template
}

func New(postTmpl string, files []string) (*Blog, error) {
	tmpl := template.Must(template.ParseFiles(postTmpl))
	blog := &Blog{posts: make(map[string]*post)}

	for _, file := range files {
		name := CleanURL(file)
		post := &post{path: file, tmpl: tmpl}

		if err := post.Load(); err != nil {
			return nil, err
		}

		blog.posts[name] = post
	}

	return blog, nil
}

func (b *Blog) Paths() (paths []string) {
	for _, post := range b.posts {
		paths = append(paths, post.path)
	}
	return
}

func (b *Blog) Get(file string) *post {
	return b.posts[CleanURL(file)]
}

func (p *post) Load() error {
	bytes, err := ioutil.ReadFile(p.path)
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

func CleanURL(name string) string {
	re := regexp.MustCompile("^[0-9]{10}-")
	base := strings.TrimLeft(
		strings.TrimRight(strings.TrimRight(name, ".html"), ".md"), "/")
	return re.ReplaceAllString(base, "")
}

func FindPostFiles(glob string) ([]string, error) {
	return filepath.Glob(glob)
}
