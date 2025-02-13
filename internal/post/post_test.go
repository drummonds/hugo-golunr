package post

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

// absTestPathToPost testing conversion success
func absNormalTestPathToPost(t *testing.T, fs afero.Fs, path string, want string) {
	post, err := PathToPost(fs, path)
	if post.URI != want || err != nil {
		t.Fatalf(`PathToPost(fs,"%s") = "%s", %v, wanted "%s"`, path, post.URI, err, want)
	}
}

// PathToPost testing conversion expecting error
func absErrorTestPathToPost(t *testing.T, fs afero.Fs, path string, want string) {
	post, err := PathToPost(fs, path)
	if err == nil {
		t.Fatalf(`PathToPost(fs,"%s") -> "%s", %v, Expected error with  %s`, path, post.URI, err, want)
	}
	error_as_string := fmt.Sprintf("%v", err)
	want_r := regexp.MustCompile(`\b` + want + `\b`)
	if !want_r.MatchString(error_as_string) {
		t.Fatalf(`PathToPost("%s") -> "%s", %s, Expected error with  %s`, path, post.URI, error_as_string, want)
	}
}

// Testing blank files
func TestPathToPostMockFS(t *testing.T) {
	AppFs := afero.NewMemMapFs()

	// Set up test files in virtual filesystem
	files := []string{
		"/test/pathtest.md",
		"/test/index.md",
		"/test/_index.md",
	}

	for _, f := range files {
		// Write empty file
		err := afero.WriteFile(AppFs, f, []byte(""), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", f, err)
		}
	}

	// absErrorTestPathToPost(t, AppFs, "", "no such file")
	absNormalTestPathToPost(t, AppFs, "/test/pathtest.md", "/test/pathtest")
	absNormalTestPathToPost(t, AppFs, "/test/index.md", "/test/")
	absNormalTestPathToPost(t, AppFs, "/test/_index.md", "/test/")
}

func abstractTestContents(t *testing.T, contents, result string, expectError bool) {
	AppFs := afero.NewMemMapFs()

	// Set up test files in virtual filesystem
	files := []string{
		"/test/index.md",
	}

	for _, f := range files {
		// Write empty file
		err := afero.WriteFile(AppFs, f, []byte(contents), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", f, err)
		}
	}
	if expectError {
		absErrorTestPathToPost(t, AppFs, "/test/index.md", result)
	} else {
		absNormalTestPathToPost(t, AppFs, "/test/index.md", result)
	}
}

func TestEmptyFile(t *testing.T) {
	abstractTestContents(t, "", "/test/", false)
}

// With new frontmatter parser, this should not fail
func TestRubbishFile(t *testing.T) {
	abstractTestContents(t, "rubbish", "/test/", false)
}

func TestEmptyRealFile(t *testing.T) {
	contents := `---
---
`
	abstractTestContents(t, contents, "/test/", false)
}

var haggisTest = `---
title: Test
date: 2025-01-01
---
This is a haggis test.
`

func TestRealFile(t *testing.T) {
	abstractTestContents(t, haggisTest, "/test/", false)
}

func abstractTestJson(t *testing.T, filePath string) {
	AppFs := afero.NewMemMapFs()

	// Set up test files in virtual filesystem
	err := afero.WriteFile(AppFs, filePath, []byte(haggisTest), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %+v", err)
	}

	post, err := PathToPost(AppFs, filePath)
	if err != nil {
		t.Fatalf("Failed to parse post: %v", err)
	}

	output, err := json.Marshal(post)
	if err != nil {
		t.Fatalf("Failed to parse post: %v", err)
	}
	if !strings.Contains(string(output), "haggis") {
		t.Fatalf("output does not include words haggis for file %s", filePath)
	}
}

func TestJson(t *testing.T) {
	abstractTestJson(t, "/index.md")
}

func TestJsonSubDir(t *testing.T) {
	abstractTestJson(t, "/Donkey/_index.md")
}
