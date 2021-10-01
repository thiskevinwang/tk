package cmd

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/manifoldco/promptui"
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

type Node struct {
	Level    int
	Children []*Node
	Parent   *Node
	Name     string
	Path     string
	IsLast   bool
}

func genNodes(root *Node, filename string, level int, filepath string) []*Node {
	children := []*Node{}

	files, err := os.ReadDir(filepath)
	if err != nil {
		return children
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

	for idx, f := range files {
		isLast := idx == len(files)-1

		child := &Node{
			Level:  level + 1,
			Parent: root,
			Name:   f.Name(),
			Path:   path.Join(filepath, f.Name()),
			IsLast: isLast,
		}
		child.Children = genNodes(child, f.Name(), level+1, path.Join(filepath, f.Name()))

		children = append(children, child)
	}

	return children
}

var treeCmd = &cobra.Command{
	Use:   `tree`,
	Short: `List contents of directories in a tree-like format`,
	Long:  `List contents of directories in a tree-like format`,
	Run: func(treeCmd *cobra.Command, args []string) {
		prompt := promptui.Select{
			Label: "Select Traversal Strategy",
			Items: []string{"Breadth First", "Depth First"},
		}

		_, result, err := prompt.Run()

		if err != nil {
			fmt.Println(aurora.Red(err.Error()), " Exiting...")
			return
		}

		fmt.Println(aurora.Green(result))

		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		// Default to current working directory
		dir := cwd
		// Or use path via 1st argument
		if len(args) >= 1 {
			dir = args[0]
		}

		// Generate in-memory tree
		root := &Node{
			Name:     dir,
			Level:    0,
			Path:     dir,
			Parent:   nil,
			Children: nil,
			IsLast:   true,
		}
		root.Children = genNodes(root, dir, 0, dir)

		switch result {
		case "Breadth First":
			bfs(root)
		case "Depth First":
			dfs(root, true)
		}
	},
}

// Not really search, just traversal
func bfs(root *Node) {
	// Instantiate queue
	queue := []*Node{}
	// Enqueue root node
	queue = append(queue, root)

	// While the queue isn't empty
	for len(queue) > 0 {
		// Dequeue first element
		el := queue[0]
		queue = queue[1:]
		// Enqueue its children
		queue = append(queue, el.Children...)

		fmt.Println(strings.Repeat(STEM, el.Level)+fmt.Sprintf("[%v]", el.Level)+">>", aurora.Blue(el.Name))

	}
}

// Also not search, just traversal
func dfs(root *Node, isLast bool) {
	// The symbol, preceding the filename
	tip := BRANCH
	if isLast {
		tip = CORNER
	}

	symbols := []string{}
	ptr := root

	// Traverse backwards to draw all preceding branches
	for ptr.Parent != nil {
		// Draw leftwards
		if ptr.Parent.IsLast {
			// If a parent is the last node out of its siblings,
			// draw empty space
			symbols = append([]string{SPACE}, symbols...)
		} else {
			// else draw a stem
			symbols = append([]string{STEM}, symbols...)
		}

		ptr = ptr.Parent
	}

	branch := strings.Join(symbols, "") + tip
	name := aurora.Yellow(root.Name)

	fmt.Println(branch, name)

	for i, next := range root.Children {
		isLast := i == len(root.Children)-1
		dfs(next, isLast)
	}
}
