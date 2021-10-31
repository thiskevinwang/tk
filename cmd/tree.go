package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/logrusorgru/aurora"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(treeCmd)
}

const (
	SPACE  = "    "
	BRANCH = "│   "
	STEM   = "├──"
	CORNER = "└──"
)

var treeCmd = &cobra.Command{
	Use:   `tree`,
	Short: `List contents of directories in a tree-like format`,
	Long:  `List contents of directories in a tree-like format`,
	Run: func(treeCmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		// Default to current working directory
		dir := cwd
		// Or use path via 1st argument
		if len(args) >= 1 {
			dir = path.Join(dir, args[0])
		}

		fmt.Println(aurora.Blue(dir))

		fileCt, folderCt := dfs(dir, dir, map[string]string{})
		fmt.Println("Files:", fileCt, "Folders:", folderCt)
	},
}

// GIVEN a path like '/Users/kevin/repos/tk'
//
// RETURN true|false, if the path is the last entry amongst its siblings
func getIsLast(path string) bool {
	dir := filepath.Dir(path)   // Users/kevin/repos
	base := filepath.Base(path) // tk

	// TODO handle errors
	files, _ := os.ReadDir(dir)
	isLast := false
	for i, file := range files {
		if base == file.Name() {
			isLast = i == len(files)-1
			break
		}
	}

	return isLast
}

func dfs(dir string, root string, cache map[string]string) (int, int) {
	fileCt := 0
	folderCt := 0

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		filePath := path.Join(dir, f.Name())
		// // TODO make ignore-cases configurable
		switch name := f.Name(); {
		case
			strings.Contains(name, ".git"),         // dot files
			strings.Contains(name, "node_modules"): // node modules
			continue
		}

		relativePath := strings.Split(filePath, root)
		parts := strings.Split(relativePath[1], string(os.PathSeparator))

		res := ""
		abs := root

		for i, part := range parts[1:] {
			abs = path.Join(abs, part)

			if cache[abs] != "" {
				res += cache[abs]
			} else {
				isLast := getIsLast(abs)

				if i == len(parts[1:])-1 {
					if isLast {
						res += CORNER
					} else {
						res += STEM
					}
				} else {
					if isLast {
						res += SPACE
					} else {
						res += BRANCH
					}
					cache[abs] = res
				}
			}

		}

		fmt.Println(aurora.Gray(12, res), f.Name())

		if f.IsDir() {
			_fileCt, _folderCt := dfs(filePath, root, cache)
			fileCt += _fileCt
			folderCt += 1 + _folderCt
		} else {
			fileCt += 1
		}
	}

	return fileCt, folderCt
}

// This func has the same output as `dfs` but relies on
// filepath.Walkdir.
// This func is only used in benchmark tests.
func dfs_walk(dir string, root string) (int, int) {
	fileCt := 0
	folderCt := 0

	filepath.WalkDir(dir, func(_path string, d os.DirEntry, err error) error {
		switch p := _path; {
		case
			p == root,                           // avoid infinite loop on root
			strings.Contains(p, ".git"),         // dot files
			strings.Contains(p, "node_modules"): // node modules
			return err
		}

		if d.IsDir() {
			folderCt += 1
		} else {
			fileCt += 1
		}
		res := ""

		relativePath := strings.Split(_path, root)
		parts := strings.Split(relativePath[1], string(os.PathSeparator))

		abs := root
		for i, part := range parts[1:] {
			abs = path.Join(abs, part)
			isLast := getIsLast(abs)

			if i == len(parts[1:])-1 {
				if isLast {
					res += CORNER
				} else {
					res += STEM
				}
			} else {
				if isLast {
					res += SPACE
				} else {
					res += BRANCH
				}
			}
		}
		fmt.Println(res, d.Name())
		return err
	})

	return fileCt, folderCt
}
