package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var mtx sync.Mutex
var wg sync.WaitGroup
var posts []Post

// baseURL should be parsed from the config.toml file in the hugo repo
func main() {

	fmt.Println("Version 1.2.0 2024-12-31")

	// Initialize posts with a reasonable capacity to reduce reallocations
	posts = make([]Post, 0, 100) // adjust capacity based on expected number of posts

	filepath.Walk("./content", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error while walking content directory: ", err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(info.Name(), ".md") {
			wg.Add(1)
			go ParsePost(path)
		}
		return nil
	})
	wg.Wait()
	output, err := json.Marshal(posts)
	if err != nil {
		fmt.Println("Could not marshal posts to JSON: ", err)
		return
	}
	err = os.WriteFile("static/search_index.json", output, 0644)
	if err != nil {
		fmt.Println("Could not write file: ", err)
		return
	}
	fmt.Println("Done!")
}
