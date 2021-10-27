package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

// Run benchmarks with:
// go test ./cmd -bench=.

func BenchmarkDfsWalk(b *testing.B) {
	os.Stdout = nil

	cwd, _ := os.Getwd()
	p := filepath.Join(cwd, "..")
	for n := 0; n < b.N; n++ {
		dfs_walk(p, p)
	}
}

func BenchmarkDfs(b *testing.B) {
	os.Stdout = nil

	cwd, _ := os.Getwd()
	p := filepath.Join(cwd, "..")
	for n := 0; n < b.N; n++ {
		dfs(p, p)
	}
}
