package fspress

import (
	"reflect"
	"testing"
)

func eq(t *testing.T, expecting, returned interface{}) {
	if !reflect.DeepEqual(expecting, returned) {
		t.Logf("expecting: %+v\n", expecting)
		t.Logf("returned: %+v\n", returned)
		t.Fatalf("assertion error")
	}
}

func TestCleanURL(t *testing.T) {
	eq(t, "one", cleanURL("1530415742-one.md"))
	eq(t, "one-two", cleanURL("1530415745-one-two.md"))
	eq(t, "one-two-three", cleanURL("1530415749-one-two-three.md"))
}

// func TestStringifyingPosts(t *testing.T) {
// 	tmpl := template.Must(template.New("").Parse("-{{.Content}}-"))
// 	post := &Post{Content: "hi", tmpl: tmpl}
// 	eq(t, "-hi-", post.String())
// }
//
// func TestParseGlobFindsAllFiles(t *testing.T) {
// 	blog := Must(ParseGlob("test/template.tmpl", "test/[0-9]*.md"))
// 	eq(t, 3, len(blog.Posts))
// }
//
// func TestGetUsesCleanURLs(t *testing.T) {
// 	blog := Must(ParseGlob("test/template.tmpl", "test/[0-9]*.md"))
//
// 	post1 := blog.Get("one")
// 	post2 := blog.Get("/one")
// 	post3 := blog.Get("/one.html")
//
// 	if post1 == nil {
// 		t.Fatal("not expecting nil")
// 	}
//
// 	eq(t, 3, len(blog.Posts))
// 	eq(t, post1, post2)
// 	eq(t, post2, post3)
// }
//
// func TestPostGeneration(t *testing.T) {
// 	blog := Must(ParseGlob("test/template.tmpl", "test/[0-9]*.md"))
// 	post := blog.Get("one")
// 	post.Content = "hi"
// 	eq(t, "-hi-", strings.TrimSpace(post.String()))
// }
