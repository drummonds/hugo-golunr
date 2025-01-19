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

	if tags, ok := f["tags"]; ok {
		if tagSlice, ok := tags.([]interface{}); ok {
			post.Tags = make([]string, len(tagSlice))
			for i, tag := range tagSlice {
				if str, ok := tag.(string); ok {
					post.Tags[i] = str
				}
			}
		}
	}

	text := stripmd.Strip(body)
	text = strip.StripTags(text)
	if WordSet {
		// Convert text to unique words
		words := strings.Fields(text)
		uniqueWords := make(map[string]struct{})
		for _, word := range words {
			uniqueWords[word] = struct{}{}
		}

		// Convert back to space-separated string
		var uniqueWordList []string
		for word := range uniqueWords {
			uniqueWordList = append(uniqueWordList, word)
		}
		text = strings.Join(uniqueWordList, " ")
	}
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
	if Verbose && path == "content/humphrey/_index.md" {
		fmt.Printf("   path: %s\n", path)
	}
	post, err := PathToPost(fs, path)
	if err != nil {
		fmt.Printf("error parsing path %s", path)
		return
	}

	if Verbose {
		fmt.Printf("Processing: %s\n", path)
		fmt.Printf("   url: %s\n", post.URI)
	}
	if Verbose && path == "content/humphrey/_index.md" {
		fmt.Printf("   title: %s\n", post.Title)
		fmt.Printf("   tags: %s\n", post.Tags)
		fmt.Printf("   content: %s\n", post.Content)
	}
	// The template needs to use the baseURL to form a compete URL.  This allows the same
	// json file to be used on different sites eg development and production.
	AddPost(post)

}
