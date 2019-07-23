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
	catalogPath  string
	postTmplPath string
	filesGlob    string
	Posts        map[string]*Post
}

func New(catalogPath, postTmplPath string, filesGlob string) *Blog {
	return &Blog{
		catalogPath:  catalogPath,
		postTmplPath: postTmplPath,
		filesGlob:    filesGlob,
		Posts:        make(map[string]*Post),
	}
}

func (blog *Blog) Load() error {
	files, err := filepath.Glob(blog.filesGlob)
	if err != nil {
		return err
	}

	tmpl, err := template.ParseFiles(blog.postTmplPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		url := cleanURL(filepath.Base(file))

		slug, err := fileNameSlug(file)
		if err != nil {
			return err
		}

		dateStr, err := fileNameDate(file)
		if err != nil {
			return err
		}
		i, err := strconv.ParseInt(dateStr, 10, 64)
		if err != nil {
			return err
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
			return err
		}

		blog.Posts[url] = post
	}

	return nil
}

// Get looks up a post by cleaning up the provided URL string
func (b *Blog) Get(file string) *Post {
	return b.Posts[cleanURL(file)]
}

// Post hold information about a blog post
type Post struct {
	URL     string
	Path    string
	Slug    string
	Title   string
	Content string
	Date    time.Time
	tmpl    *template.Template
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
