package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gernest/front"
	strip "github.com/grokify/html-strip-tags-go"
	stripmd "github.com/writeas/go-strip-markdown"
)

type Post struct {
	URI     string   `json:"location"`
	Title   string   `json:"title"`
	Content string   `json:"text"`
	Tags    []string `json:"tags"`
}

func PathToPost(path string) (post Post, err error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error while reading file: ", path, err)
		return post, err
	}

	m := front.NewMatter()
	m.Handle("---", front.YAMLHandler)
	f, body, err := m.Parse(strings.NewReader(string(buf)))
	if err != nil {
		fmt.Println("Error while parsing file: ", path, err)
		return post, err
	}

	// post := Post{}
	if title, ok := f["title"]; ok {
		post.Title = stripmd.Strip(title.(string))
	}
	text := stripmd.Strip(body)
	text = strip.StripTags(text)

	post.Content = text
	uri := strings.ToLower(strings.TrimPrefix(path, "content"))
	uri = strings.TrimSuffix(uri, ".md")
	uri = strings.Replace(uri, "_index", "", 1)
	uri = strings.Replace(uri, "index", "", 1)
	post.URI = strings.Replace(uri, " ", "-", -1)

	return post, nil
}
func ParsePost(path string) {
	post, err := PathToPost(path)
	if err != nil {
		panic(fmt.Sprintf("error parsing path %s", path))
	}

	fmt.Printf("Parsed %s\n", post.URI)
	// The template needs to use the baseURL to form a compete URL.  This allows the same
	// json file to be used on different sites eg development and production.
	mtx.Lock()
	posts = append(posts, post)
	mtx.Unlock()

	wg.Done()

}
