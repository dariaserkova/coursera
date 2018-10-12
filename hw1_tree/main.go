package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func sizer(size int64) string {
	if size == 0 {
		return "empty"
	}
	s := fmt.Sprintf("%d", size)
	return s + "b"
}

func splitter(path string, last bool) string {
	result := ""
	firstSep := "│"
	sep := "├───"
	if last {
		sep = "└───"
	}
	osSep := string(os.PathSeparator)
	slashes := strings.Count(path, osSep)
	switch slashes {
	case 0:
		result += sep + " " + path + "\n"
	case 1:
		result += firstSep + "\t" + sep + " " + filepath.Base(path) + "\n"
	default:
		result += firstSep + "\t" + splitter(strings.SplitAfterN(path, osSep, 2)[1], last)
	}
	return result
}

func printer(files []string) string {
	result := ""
	for i, item := range files {
		last := false
		if item == "." {
			continue
		}
		if i < len(files)-1 && filepath.Dir(item) != filepath.Dir(files[i+1]) {
			last = true
		}
		if i == len(files)-1 {
			last = true
		}
		result += splitter(item, last)
	}
	return result
}

func dirTree(out io.Writer, root string, printFiles bool) error {
	list := make([]string, 0, 25)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if printFiles {
			if !info.IsDir() {
				list = append(list, path+" ("+sizer(info.Size())+")")
			} else {
				list = append(list, path)
			}
		} else {
			if info.IsDir() {
				list = append(list, path)
			}
		}
		return nil

	})
	sort.Strings(list)
	io.WriteString(out, printer(list))
	return err
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
