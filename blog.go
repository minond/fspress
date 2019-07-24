// Package fspress generates blog content from provided templates and file
// system content.
package fspress

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	blackfriday "github.com/russross/blackfriday/v2"
)

var fileNameDateRe = regexp.MustCompile("([0-9]+)-")

type Metadata struct {
	Title    string
	Abstract string
}

// Blog holds all information about this blog and its posts in memory.
type Blog struct {
	catalogPath  string
	postTmplPath string
	filesGlob    string
	Catalog      map[string]*Metadata
	Posts        map[string]*Post
}

func New(catalogPath, postTmplPath string, filesGlob string) *Blog {
	return &Blog{
		catalogPath:  catalogPath,
		postTmplPath: postTmplPath,
		filesGlob:    filesGlob,
	}
}

func (blog *Blog) Load() error {
	if err := blog.loadCatalog(); err != nil {
		return err
	}
	return blog.loadPosts()
}

func (blog *Blog) loadCatalog() error {
	catalog := make(map[string]*Metadata)

	handle, err := os.Open(blog.catalogPath)
	if err != nil {
		return err
	}

	reader := csv.NewReader(handle)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, entry := range records[1:] {
		catalog[entry[0]] = &Metadata{
			Title:    entry[1],
			Abstract: entry[2],
		}
	}

	blog.Catalog = catalog
	return nil
}

func (blog *Blog) loadPosts() error {
	posts := make(map[string]*Post)

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
		post, err := blog.generatePost(file, tmpl)
		if err != nil {
			return err
		}
		posts[url] = post
	}

	blog.Posts = posts
	return nil
}

func (blog *Blog) generatePost(file string, tmpl *template.Template) (*Post, error) {
	url := cleanURL(filepath.Base(file))

	metadata, ok := blog.Catalog[url]
	if !ok {
		return nil, fmt.Errorf("%s (%s) is not in the catalog (%s)",
			file, url, blog.catalogPath)
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
		URL:      url,
		Path:     file,
		Date:     date,
		Metadata: metadata,
		tmpl:     tmpl,
	}

	if err := post.Load(); err != nil {
		return nil, err
	}

	return post, nil
}

// Get looks up a post by cleaning up the provided URL string
func (blog *Blog) Get(file string) *Post {
	return blog.Posts[cleanURL(file)]
}

// Post hold information about a blog post
type Post struct {
	URL      string
	Path     string
	Content  string
	Date     time.Time
	Metadata *Metadata
	tmpl     *template.Template
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
