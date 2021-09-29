package cmd

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/logrusorgru/aurora"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(treeCmd)
}

const (
	space  = "   "
	stem   = "│  "
	branch = "├──"
	corner = "└──"
)

// ├ ─ ─
// └ ─ ─
// │
func tree(filename string, level int, prefix []string, full string) {
	if filename == ".git" {
		return
	}

	fmt.Println(strings.Join(prefix, ""), aurora.Green(filename))

	files, err := os.ReadDir(full)
	if err != nil {
		// fmt.Println(filename)
		// log.Fatal(filename, err.Error())
		return
	}

	// sort by [..folder, ...file]
	sort.Slice(files, func(i, j int) bool {
		iIsDir := files[i].IsDir()
		jIsDir := files[j].IsDir()
		if !iIsDir && jIsDir {
			return false
		} else {
			return true
		}
	})

	for i, f := range files {
		// isDir := f.IsDir()

		var newprefix []string
		isLast := i == len(files)-1
		el := branch

		if level > 0 {
			if isLast {
				el = corner
			}
			// newprefix = append([]string{el}, prefix...)
			// newprefix = append(prefix, el)
			newprefix = []string{}
			for range prefix {
				newprefix = append(newprefix, stem)
			}
			newprefix = append(newprefix, el)

		} else {
			if isLast {
				el = corner
			}
			newprefix = []string{}
			for range prefix {
				newprefix = append(newprefix, stem)
			}
			newprefix = append(newprefix, el)
		}

		tree(f.Name(), level+1, newprefix, path.Join(full, f.Name()))
	}
}

var treeCmd = &cobra.Command{
	Use:   `tree`,
	Short: `List contents of directories in a tree-like format`,
	Long:  `List contents of directories in a tree-like format`,
	Run: func(treeCmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		dir := cwd
		if len(args) >= 1 {
			dir = args[0]
		}

		fmt.Println(dir)
		tree(dir, 0, []string{}, dir)
	},
}
