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

type stringSlice []string

func (i *stringSlice) String() string {
	return strings.Join(*i, ",")
}

func (i *stringSlice) Set(value string) error {
	for _, path := range strings.Split(value, ",") {
		*i = append(*i, path)
	}
	return nil
}

var helpMessage = `USAGE:
    mgg [OPTIONS]

OPTIONS:
    -h, --help      Prints this message
    -d, --dir       Directory to generate mocks in [default: 'mocks']
    -p, --prefix    Prefix to use for generated mocks [default: 'mock_']
    -i, --ignore    Paths to ignore when scanning for interfaces [default: ['']]
`

var flags struct {
	help   bool
	dir    string
	prefix string
	ignore stringSlice
}

func init() {
	flag.BoolVar(&flags.help, "help", false, "Prints this message")
	flag.BoolVar(&flags.help, "h", false, "Prints this message shorthand")
	flag.StringVar(&flags.dir, "dir", "mocks", "`Directory` to generate mocks in")
	flag.StringVar(&flags.dir, "d", "mocks", "`Directory` to generate mocks in shorthand")
	flag.StringVar(&flags.prefix, "prefix", "mock_", "`Prefix` to use for generated mocks")
	flag.StringVar(&flags.prefix, "p", "mock_", "`Prefix` to use for generated mocks")
	flag.Var(&flags.ignore, "ignore", "`Paths` to ignore when scanning for interfaces")
	flag.Var(&flags.ignore, "i", "`Paths` to ignore when scanning for interfaces")
}

func main() {
	flag.Parse()

	if flags.help {
		fmt.Println(helpMessage)
		return
	}

	if err := generate(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func generate() error {
	cwd, _ := os.Getwd()

	_, err := os.Stat(filepath.Join(cwd, "go.mod"))
	if errors.Is(err, os.ErrNotExist) {
		return errors.New("go.mod not found")
	}

	files, err := getFiles()
	if err != nil {
		return err
	}

	pathSeparator := string(os.PathSeparator)
	mockFolder := flags.dir
	mockPrefix := flags.prefix
	mockPathSeparatorPrefix := pathSeparator + mockPrefix

	for _, src := range files {
		dest := strings.TrimPrefix(src, cwd)
		dest = strings.Replace(dest, pathSeparator, mockPathSeparatorPrefix, -1)
		dest = filepath.Join(mockFolder, dest)
		src = strings.TrimPrefix(src, cwd)
		src = strings.TrimPrefix(src, pathSeparator)

		if !strings.Contains(src, flags.prefix) {
			cmd := exec.Command("mockgen", "-source", src, "-destination", dest)
			if err := cmd.Run(); err != nil {
				return err
			}
			fmt.Printf("Generated '%s'\n", strings.TrimPrefix(dest, cwd))
		}
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

		for _, ignorePath := range flags.ignore {
			if strings.Contains(path, ignorePath) {
				return nil
			}
		}

		if strings.HasSuffix(path, ".go") && !strings.Contains(path, "_test") {
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
