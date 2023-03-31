package main

import (
	"fmt"
	"io/ioutil"
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

func ParsePost(path string) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error while reading file: ", path, err)
		return
	}

	m := front.NewMatter()
	m.Handle("---", front.YAMLHandler)
	f, body, err := m.Parse(strings.NewReader(string(buf)))

	post := Post{}
	if title, ok := f["title"]; ok {
		post.Title = stripmd.Strip(title.(string))
	}
	text := stripmd.Strip(body)
	text = strip.StripTags(text)

	post.Content = text
	uri := strings.Replace(strings.TrimSuffix(strings.TrimPrefix(path, "content"), ".md"), "index", "", -1)
	uri = strings.ToLower(uri)
	post.URI = strings.Replace(uri, " ", "-", -1)

	fmt.Printf("Parsed %s\n", post.URI)
	// The template needs to use the baseURL to form a compete URL.  This allows the same
	// json file to be used on different sites eg development and productio.
	mtx.Lock()
	posts = append(posts, post)
	mtx.Unlock()

	wg.Done()

}
