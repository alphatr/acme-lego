package main

import (
	"fmt"
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
		shortHash = githubSHA[:7]
	}

	ref := os.Getenv("GITHUB_REF")
	version := shortHash
	if ref != "" {
		version = strings.ReplaceAll(ref, "refs/tags/v", "")
	}

	currentTime := time.Now().Format(time.RFC3339)
	content := fmt.Sprintf(template, version, shortHash, currentTime)

	if err := os.WriteFile("./version.go", []byte(content), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write version.go: %v\n", err)
		os.Exit(1)
	}
}
