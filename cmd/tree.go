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
	STEM   = "│   "
	BRANCH = "├──"
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
		fileCt, folderCt := dfs(dir, dir)
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

// GIVEN a path like "/Users/kevin/repos/tk/cmd/version.go"
//
// RETURN a string like `│   │       │       └── `
//
// Traverses a path upwards, and checks if each subsequent entry
// is the last amongst its siblings and draws a stem/space accordingly.
func drawBranch(path string, root string) string {
	if !filepath.HasPrefix(path, root) {
		log.Fatal("Invariant Violation: ", aurora.Underline(path), " does not extend ", aurora.Underline(root))
	}
	dir := filepath.Dir(path)

	tip := BRANCH
	if getIsLast(path) {
		tip = CORNER
	}
	symbols := []string{tip}

	// draw leftwards until we reach the root directory
	for dir != root {
		if getIsLast(dir) {
			symbols = append([]string{SPACE}, symbols...)
		} else {
			symbols = append([]string{STEM}, symbols...)
		}
		dir = filepath.Dir(dir)
	}

	return strings.Join(symbols, "")
}

// Traverse a given directory tree and print a tree-like output
//
// - dir is the path to each node
//
// - root is path to the root dir which dfs was first called on
func dfs(dir string, root string) (int, int) {
	fileCt := 0
	folderCt := 0

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		filePath := path.Join(dir, f.Name())

		// TODO make ignore-cases configurable
		// switch name := f.Name(); {
		// case
		// 	strings.HasPrefix(name, "."),           // dot files
		// 	strings.Contains(name, "node_modules"): // node modules
		// 	continue
		// }

		branch := aurora.Gray(12, drawBranch(filePath, root))
		fmt.Println(branch, f.Name())

		if f.IsDir() {
			folderCt += 1
			_fileCt, _folderCt := dfs(filePath, root)
			fileCt += _fileCt
			folderCt += _folderCt
		} else {
			fileCt += 1
		}
	}

	return fileCt, folderCt
}

// This func has the same output as `dfs` but relies on
// filepath.Walkdir.
// This func is only used in benchmark tests.
func dfs_walk(dir string, root string) {
	fileCt := 0
	folderCt := 0

	filepath.WalkDir(dir, func(_path string, d os.DirEntry, err error) error {
		if _path == root {
			return err
		}

		if d.IsDir() {
			folderCt += 1
		} else {
			fileCt += 1
		}

		filePath := _path

		branch := aurora.Gray(12, drawBranch(filePath, root))
		fmt.Println(branch, d.Name())
		return err
	})

	fmt.Println("Files: ", fileCt, "Folders: ", folderCt)
	return
}
