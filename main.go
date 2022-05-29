package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var flags struct {
	remove bool
	dir    string
	prefix string
}

func init() {
	flag.BoolVar(&flags.remove, "remove", false, "remove old mock files before generating")
	flag.BoolVar(&flags.remove, "r", false, "remove old mock files before generating shorthand")
	flag.StringVar(&flags.dir, "dir", "mocks", "`directory` to generate mock files in")
	flag.StringVar(&flags.dir, "d", "mocks", "`directory` to generate mock files in shorthand")
	flag.StringVar(&flags.prefix, "prefix", "mock_", "`prefix` to use for mock files")
	flag.StringVar(&flags.prefix, "p", "mock_", "`prefix` to use for mock files shorthand")
}

func main() {
	flag.Parse()
	if err := generate(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func generate() error {
	cwd, _ := os.Getwd()

	_, err := os.Stat(filepath.Join(cwd, flags.dir))
	if flags.remove && !errors.Is(err, os.ErrNotExist) {
		if err := os.RemoveAll(filepath.Join(cwd, flags.dir)); err != nil {
			return err
		}
	}

	_, err = os.Stat(filepath.Join(cwd, "go.mod"))
	if errors.Is(err, os.ErrNotExist) {
		return errors.New("go.mod not found")
	}

	files, err := getFiles()
	if err != nil {
		return err
	}

	pathSeparator := string(os.PathSeparator)
	mockFolder := filepath.Join(cwd, flags.dir)
	mockPrefix := flags.prefix
	mockPathSeparatorPrefix := pathSeparator + mockPrefix

	for _, src := range files {
		dest := strings.TrimPrefix(src, cwd)
		dest = strings.Replace(dest, pathSeparator, mockPathSeparatorPrefix, -1)
		dest = mockFolder + dest

		src = strings.TrimPrefix(src, cwd)
		src = strings.TrimPrefix(src, pathSeparator)
		cmd := exec.Command("mockgen", "-source", src, "-destination", dest)
		if err := cmd.Run(); err != nil {
			return err
		}
		fmt.Printf("Generated '%s'\n", strings.TrimPrefix(dest, cwd))
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

		if strings.HasSuffix(path, ".go") &&
			!strings.Contains(path, "_test") &&
			!strings.Contains(path, "mock") {
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
