// Package fspress generates blog content from provided templates and file
// system content.
package fspress

import (
	"bytes"
	"errors"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	blackfriday "github.com/russross/blackfriday/v2"
)

var fileNameDateRe = regexp.MustCompile("([0-9]+)-")
var fileNameSlugRe = regexp.MustCompile("[0-9]+-(.+)\\.")

// Blog holds all information about this blog and its posts in memory.
type Blog struct {
	Posts map[string]*Post
}

// Post hold information about a blog post
type Post struct {
	URL     string
	Path    string
	Slug    string
	Content string
	Date    time.Time
	tmpl    *template.Template
}

// Must asserts a Blog creation was successful and panics when it is not.
func Must(blog *Blog, err error) *Blog {
	if err != nil {
		panic(err)
	}

	return blog
}

// ParseFiles generates a Blog using the provided array of file paths
func ParseFiles(postTmpl string, files []string) (*Blog, error) {
	tmpl := template.Must(template.ParseFiles(postTmpl))
	blog := &Blog{Posts: make(map[string]*Post)}

	for _, file := range files {
		url := cleanURL(filepath.Base(file))

		slug, err := fileNameSlug(file)
		if err != nil {
			return nil, err
		}

		dateStr, err := fileNameDate(file)
		if err != nil {
			return nil, err
		}
		i, err := strconv.ParseInt(dateStr, 10, 64)
		if err != nil {
			return nil, err
		}
		date := time.Unix(i, 0)

		post := &Post{
			URL:  url,
			Path: file,
			Slug: slug,
			Date: date,
			tmpl: tmpl,
		}

		if err := post.Load(); err != nil {
			return nil, err
		}

		blog.Posts[url] = post
	}

	return blog, nil
}

// ParseGlob generates a Blog using the provided glob string
func ParseGlob(postTmpl string, glob string) (*Blog, error) {
	files, err := filepath.Glob(glob)
	if err != nil {
		return nil, err
	}

	return ParseFiles(postTmpl, files)
}

// Get looks up a post by cleaning up the provided URL string
func (b *Blog) Get(file string) *Post {
	return b.Posts[cleanURL(file)]
}

// Load reads a post's contents from the file system
func (p *Post) Load() error {
	bytes, err := ioutil.ReadFile(p.Path)
	if err != nil {
		return err
	}

	p.Content = string(blackfriday.Run(bytes))
	return nil
}

// String returns a fully generated blog post
func (p *Post) String() string {
	buf := &bytes.Buffer{}
	p.tmpl.Execute(buf, p)
	return buf.String()
}

func fileNameSlug(name string) (string, error) {
	match := fileNameSlugRe.FindStringSubmatch(name)
	if len(match) != 2 {
		return "", errors.New("invalid slug in file name")
	}
	return match[1], nil
}

func fileNameDate(name string) (string, error) {
	match := fileNameDateRe.FindStringSubmatch(name)
	if len(match) != 2 {
		return "", errors.New("invalid date in file name")
	}
	return match[1], nil
}

func cleanURL(name string) string {
	base := strings.TrimLeft(
		strings.TrimRight(strings.TrimRight(name, ".html"), ".md"), "/")
	return fileNameDateRe.ReplaceAllString(base, "")
}
