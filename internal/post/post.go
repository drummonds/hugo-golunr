package post

import (
	"fmt"
	"strings"

	"github.com/gernest/front"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/spf13/afero"
	stripmd "github.com/writeas/go-strip-markdown"
)

type Post struct {
	URI     string   `json:"location"`
	Title   string   `json:"title"`
	Content string   `json:"text"`
	Tags    []string `json:"tags"`
}

func PathToPost(fs afero.Fs, path string) (post Post, err error) {
	buf, err := afero.ReadFile(fs, path)
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
func ParsePost(fs afero.Fs, path string) {
	defer wg.Done()
	post, err := PathToPost(fs, path)
	if err != nil {
		fmt.Printf("error parsing path %s", path)
		return
	}

	fmt.Printf("Parsed %s\n", post.URI)
	// The template needs to use the baseURL to form a compete URL.  This allows the same
	// json file to be used on different sites eg development and production.
	AddPost(post)

}
