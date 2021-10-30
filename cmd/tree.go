package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

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

var isLastCache = map[string]bool{}

// GIVEN a path like '/Users/kevin/repos/tk'
//
// RETURN true|false, if the path is the last entry amongst its siblings
func getIsLast(path string) bool {
	if x := isLastCache[path]; x {
		return x
	}

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

	// cache the result
	isLastCache[path] = isLast
	return isLast
}

// path will be globally unique so memoizing is not necessary
func drawTip(path string) string {
	if getIsLast(path) {
		return CORNER
	} else {
		return BRANCH
	}
}

var branchCache = map[string]string{}

// GIVEN a path like "/Users/kevin/repos/tk/cmd/version.go"
//
// RETURN a string like `│   │       │       └── `
//
// Traverses a path upwards, and checks if each subsequent entry
// is the last amongst its siblings and draws a stem/space accordingly.
func drawBranch(path string, root string) string {
	// `path` will be globally unique
	//
	// memoizing the result is not necessarily useful, assuming each
	// function invocation is isolated
	//
	// caching the result does improve benchmark runs, but it is not
	// an accurate improvement in practice

	if !filepath.HasPrefix(path, root) {
		log.Fatal("Invariant Violation: ", aurora.Underline(path), " does not extend ", aurora.Underline(root))
	}

	// `branch` will be repeatedly calculated by files sharing the
	// same parent directory, so this calculation should be cached
	branch := ""
	originalDirPath := filepath.Dir(path)
	// starting from the file's directory,
	// draw leftwards until we reach the root directory
	dir := filepath.Dir(path)
	// try to reach from the cache
	if branchCache[dir] != "" {
		branch = branchCache[dir]
	} else {
		// loop until we reach the root directory
		for dir != root {
			if getIsLast(dir) {
				branch = SPACE + branch
			} else {
				branch = STEM + branch
			}
			dir = filepath.Dir(dir)
		}
		// cache result of the original dir, before mutations
		branchCache[originalDirPath] = branch
	}

	tip := drawTip(path)

	return branch + tip
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
