package post

import (
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestLargeMarkdownFile(t *testing.T) {
	// Create a mock filesystem
	fs := afero.NewMemMapFs()

	// Create a large markdown file content
	var contentBuilder strings.Builder

	// Add frontmatter
	contentBuilder.WriteString(`---
title: "Large Test Post"
date: 2025-01-01
tags: ["test", "large"]
---

`)

	// Add a large body of content
	for i := 0; i < 1000; i++ {
		contentBuilder.WriteString(`## Section `)
		contentBuilder.WriteString(string(rune('A'+(i%26))) + fmt.Sprintf("-%d", (i/26)))
		contentBuilder.WriteString("\n\n")
		contentBuilder.WriteString("Lorem ipsum dolor sit amet, consectetur adipiscing elit. ")
		contentBuilder.WriteString("Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. ")
		contentBuilder.WriteString("Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris.\n\n")
	}
	contentBuilder.WriteString("And here is the haggis at the end.\n\n")

	content := contentBuilder.String()
	testPath := "/content/large-post.md"

	// Write the large file to mock filesystem
	err := afero.WriteFile(fs, testPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Try to parse the large file
	post, err := PathToPost(fs, testPath)
	if err != nil {
		t.Fatalf("Failed to parse large post: %v", err)
	}

	// Verify the post was parsed correctly
	if post.URI != "/content/large-post" {
		t.Errorf("Expected URI '/content/large-post', got '%s'", post.URI)
	}

	if !strings.Contains(post.Content, "Lorem ipsum") {
		t.Error("Post content was not parsed correctly")
	}

	if !strings.Contains(post.Content, "haggis") {
		t.Error("Post content did not contain the expected ending text")
	}

	// Test WordSet functionality
	originalLength := len(strings.Fields(post.Content))
	WordSet = true
	postWithWordSet, err := PathToPost(fs, testPath)
	if err != nil {
		t.Fatalf("Failed to parse large post with WordSet: %v", err)
	}
	wordSetLength := len(strings.Fields(postWithWordSet.Content))

	if wordSetLength >= originalLength {
		t.Errorf("WordSet did not reduce word count: original=%d, wordset=%d", originalLength, wordSetLength)
	}

	// Reset WordSet for other tests
	WordSet = false

	// Verify frontmatter was parsed
	if post.Title != "Large Test Post" {
		t.Errorf("Expected title 'Large Test Post', got '%s'", post.Title)
	}

	if len(post.Tags) != 2 || post.Tags[0] != "test" || post.Tags[1] != "large" {
		t.Errorf("Tags were not parsed correctly, got %v", post.Tags)
	}
}
