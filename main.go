package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type SemVer struct {
	Major, Minor, Patch int
}

func NewSemVer(commitCount int) *SemVer {
	version := &SemVer{}
	version.Major = commitCount / 100
	version.Minor = (commitCount % 100) / 10
	version.Patch = commitCount % 10
	return version
}

func (v *SemVer) String() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func GetCommitCount(path string) (int, error) {
	cmd := exec.Command("git", "rev-list", "--count", "HEAD")
	cmd.Dir = path
	out, err := cmd.Output()
	if err != nil {
		return 0, errors.New(fmt.Sprintf("Cannot get commit count for git repo in %s, due to %s", path, err))
	}

	count, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return 0, errors.New(fmt.Sprintf("Cannot convert commit count to integer, due to %s", err))
	}

	return count, nil
}

func hasGitFolder(path string) bool {
	files, err := fs.ReadDir(os.DirFS("."), path)
	if err != nil {
		log.Printf("Cannot read %s due to %s", path, err)
		return false
	}

	for _, file := range files {
		if file.Name() == ".git" {
			return true
		}
	}
	return false
}

func SetVersionInChart(chartFile string, version string) error {
	cmd := exec.Command("sed", "-i", fmt.Sprintf("s/version: .*/version: %s/g", version), chartFile)
	err := cmd.Run()
	if err != nil {
		return errors.New(fmt.Sprintf("Cannot set version in %s, due to %s", chartFile, err))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <chart_file>", os.Args[0])
	}

	chartFile := os.Args[1]
	commitCount := 0
	repoPaths := make([]string, 0)

	log.Printf("Found path .")
	repoPaths = append(repoPaths, ".")

	// walk current dir
	files, err := fs.ReadDir(os.DirFS("."), ".")
	if err != nil {
		log.Fatalf("Cannot read current dir, due to %s", err)
	}

	for _, file := range files {
		if file.Type() == fs.ModeDir && hasGitFolder(file.Name()) {
			log.Printf("Found path %s", file.Name())
			repoPaths = append(repoPaths, file.Name())
		}
	}

	for _, path := range repoPaths {
		count, err := GetCommitCount(path)
		if err != nil {
			log.Print(err)
			continue
		}

		commitCount += count
	}

	version := NewSemVer(commitCount).String()
	log.Printf("Total commit count: %d", commitCount)
	log.Printf("Semantic version: %s", version)

	err = SetVersionInChart(chartFile, version)
	if err != nil {
		log.Fatalf("Cannot update chart file version due to %s", err)
	}
}
