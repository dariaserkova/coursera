package main

import (
	"fmt"
	"sort"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func sizer(size int64) string {
	if size == 0 { return "empty" }
	s := fmt.Sprintf("%d", size)
	return s + "b"
}

func splitter(path, prevPath string) string {
	result := ""
	firstSep := "│"
	sep := "├───"
	if last {
		firstSep = ""
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
		result += firstSep + "\t" + splitter(strings.SplitAfterN(path, osSep, 2)[1], path)
	}
	return result
}

func printer(files []string) string {
	result := ""
	for i, item := range files {
		if item == "." { continue }
		if i > 0 {
			result += splitter(item, files[i-1])
		} else {
			result += splitter(item, item)
		}
	}
	return result
}


func dirTree(out io.Writer, root string, printFiles bool) error {
	list := make([]string, 0, 25)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if printFiles {
			if !info.IsDir() {
				list = append(list, path + " (" + sizer(info.Size()) + ")")
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
	//fmt.Println(printer(list))
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

