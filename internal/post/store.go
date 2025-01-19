package post

import (
	"fmt"
	"sync"

	"github.com/spf13/afero"
)

var (
	mtx   sync.Mutex
	wg    sync.WaitGroup
	posts []Post
)

// InitStore initializes the post store with the given capacity
func InitStore(capacity int) {
	posts = make([]Post, 0, capacity)
}

// AddToParseQueue adds a post to the parsing queue
func AddToParseQueue(fs afero.Fs, path string) {
	wg.Add(1)
	go ParsePost(fs, path)
}

// WaitForParsing waits for all posts to be parsed
func WaitForParsing() {
	wg.Wait()
}

// GetAllPosts returns all parsed posts
func GetAllPosts() []Post {
	fmt.Printf("\nNumber of posts: %d\n", len(posts))
	return posts
}

// AddPost adds a post to the store thread-safely
func AddPost(p Post) {
	mtx.Lock()
	defer mtx.Unlock()
	posts = append(posts, p)
}
