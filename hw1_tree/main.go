package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

const (
	branchStart = "├───"
	branchEnd   = "└───"
	tab         = "\t"
	vertTab     = "│\t"
)

func marker2(dirPath string, prefix string, isPrintFiles bool) (string, error) {
	files, errRead := ioutil.ReadDir(dirPath)
	if errRead != nil {
		return "", fmt.Errorf("Unable to get files list in " + dirPath)
	}

	if !isPrintFiles {
		var dirs []os.FileInfo
		for _, file := range files {
			if file.IsDir() {
				dirs = append(dirs, file)
			}
		}
		files = dirs
	}

	var tree string
	for idx, file := range files {
		isLastIndex := idx == len(files)-1

		branch := prefix
		if isLastIndex {
			branch += branchEnd
		} else {
			branch += branchStart
		}

		fileName := file.Name()
		tree += branch + fileName + fileSize(file) + "\n"

		if file.IsDir() {
			var newPrefix string
			nextDir := dirPath + string(os.PathSeparator) + file.Name()

			if isLastIndex {
				newPrefix = prefix + tab
			} else {
				newPrefix = prefix + vertTab
			}
			branches, errMark := marker2(nextDir, newPrefix, isPrintFiles)
			if errMark != nil {
				return "", fmt.Errorf("Unable to get files list in " + nextDir)
			}
			tree += branches
		}
	}

	return tree, nil
}

func dirTree(out io.Writer, rootPath string, printFiles bool) error {
	root, errAbs := filepath.Abs(rootPath)
	if errAbs != nil {
		return fmt.Errorf("Unable to get root path")
	}

	tree, err := marker2(root, "", printFiles)
	if err != nil {
		fmt.Errorf("Error on get files in " + root)
	}

	tree = tree[:len(tree)-1]

	fmt.Fprintln(out, tree)

	return nil
}

func fileSize(fileInfo os.FileInfo) string {
	if fileInfo.IsDir() {
		return ""
	}

	if size := fileInfo.Size(); size > 0 {
		return " (" + strconv.FormatInt(size, 10) + "b)"
	}

	return " (empty)"
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
