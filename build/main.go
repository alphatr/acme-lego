package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const template = `package main

// Build Variable
var (
	GitVersionTag = "%s"
	GitHash       = "%s"
	CurrentTime   = "%s"
)
`

func main() {
	githubSHA := os.Getenv("GITHUB_SHA")
	shortHash := "0000000"
	if githubSHA != "" {
		shortHash = string([]byte(githubSHA)[:7])
	}

	ref := os.Getenv("GITHUB_REF")
	version := shortHash
	if ref != "" {
		version = strings.Replace(ref, "refs/tags/v", "", -1)
	}

	currentTime := time.Now().Format(time.RFC3339)
	content := fmt.Sprintf(template, version, shortHash, currentTime)

	ioutil.WriteFile("./version.go", []byte(content), 0644)
}
