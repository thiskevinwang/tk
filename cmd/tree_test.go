package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

// Run benchmarks with:
// go test ./cmd -bench=.

// func BenchmarkTreeWalk(b *testing.B) {
// 	os.Stdout = nil

// 	cwd, _ := os.Getwd()
// 	p := filepath.Join(cwd, "..")
// 	for n := 0; n < b.N; n++ {
// 		dfs_walk(p, p)
// 	}
// }

func BenchmarkTreeTk(b *testing.B) {
	os.Stdout = nil

	cwd, _ := os.Getwd()
	p := filepath.Join(cwd, "..")
	for n := 0; n < b.N; n++ {
		dfs(p, p, map[string]string{})
	}
}

func BenchmarkTreeNextJs(b *testing.B) {
	os.Stdout = nil

	cwd, _ := os.Getwd()
	p := filepath.Join(cwd, "../../next.js")
	for n := 0; n < b.N; n++ {
		dfs(p, p, map[string]string{})
	}
}
