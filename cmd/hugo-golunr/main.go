package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/riesinger/hugo-golunr/internal/post"
	"github.com/spf13/afero"
)

var _fs = afero.NewOsFs()

// baseURL should be parsed from the config.toml file in the hugo repo
func main() {
	flag.BoolVar(&post.Verbose, "v", false, "enable verbose output")
	flag.BoolVar(&post.Verbose, "verbose", false, "enable verbose output")
	flag.BoolVar(&post.WordSet, "w", false, "Convert content to set of words so only appear once")
	flag.BoolVar(&post.WordSet, "wordset", false, "Convert content to set of words so only appear once")
	flag.Parse()

	fmt.Println("Version 1.6.0 2025-01-19")

	// Initialize the post store
	post.InitStore(100) // adjust capacity based on expected number of posts

	filepath.Walk("./content", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error while walking content directory: ", err)
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(info.Name(), ".md") {
			post.AddToParseQueue(_fs, path)
		}
		return nil
	})
	post.WaitForParsing()

	output, err := json.Marshal(post.GetAllPosts())
	if err != nil {
		fmt.Println("Could not marshal posts to JSON: ", err)
		return
	}
	err = afero.WriteFile(_fs, "static/search_index.json", output, 0644)
	if err != nil {
		fmt.Println("Could not write file: ", err)
		return
	}
	fmt.Println("Done!")
}
