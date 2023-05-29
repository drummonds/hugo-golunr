package main

import (
	"fmt"
	"regexp"
	"testing"
)

// absTestPathToPost testing conversion success
func absNormalTestPathToPost(t *testing.T, path string, want string) {
	post, err := PathToPost(path)
	if post.URI != want || err != nil {
		t.Fatalf(`PathToPost("%s") = "%s", %v, wanted "%s", <nil>`, path, post.URI, err, want)
	}
}

// PathToPost testing conversion expecting error
func absErrorTestPathToPost(t *testing.T, path string, want string) {
	post, err := PathToPost(path)
	if err == nil {
		t.Fatalf(`PathToPost("%s") = "%s", %v, Expected error with  %s`, path, post.URI, err, want)
	}
	error_as_string := fmt.Sprintf("%v", err)
	want_r := regexp.MustCompile(`\b` + want + `\b`)
	if !want_r.MatchString(error_as_string) {
		t.Fatalf(`PathToPost("%s") = "%s", %s, Expected error with  %s`, path, post.URI, error_as_string, want)
	}
}

// Null case
func TestPathToPost(t *testing.T) {
	absErrorTestPathToPost(t, "", "no such file")
	absNormalTestPathToPost(t, "./test/pathtest.md", "./test/pathtest")
	absNormalTestPathToPost(t, "./test/index.md", "./test/")
	absNormalTestPathToPost(t, "./test/_index.md", "./test/")
}
