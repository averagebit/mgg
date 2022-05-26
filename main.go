package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if err := generate(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func generate() error {
	files, err := getFiles()
	if err != nil {
		return err
	}

	cwd, _ := os.Getwd()
	pathSeparator := string(os.PathSeparator)
	mockFolder := filepath.Join(cwd, "mocks")
	mockPrefix := "mock_"
	mockPathSeparatorPrefix := pathSeparator + mockPrefix

	for _, src := range files {
		dest := strings.TrimPrefix(src, cwd)
		dest = strings.Replace(dest, pathSeparator, mockPathSeparatorPrefix, -1)
		dest = mockFolder + dest

		cmd := exec.Command("mockgen", "-source", src, "-destination", dest)
		if err := cmd.Run(); err != nil {
			return err
		}
		fmt.Printf("Generated '%s'", strings.TrimPrefix(dest, cwd))
	}

	return nil
}

func getFiles() ([]string, error) {
	files := []string{}
	cwd, _ := os.Getwd()

	err := filepath.Walk(cwd, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		isGoFile := strings.HasSuffix(path, ".go")
		isTestFile := strings.Contains(path, "_test")

		if isGoFile && !isTestFile {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				if strings.Contains(scanner.Text(), "interface {") {
					files = append(files, path)
					break
				}
			}
		}

		return nil
	})

	return files, err
}
