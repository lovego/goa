package docs

import (
	"fmt"
	"path/filepath"
)

func Group(parentGroupDir, path string, descs []string) string {
	groupDir := filepath.Join(parentGroupDir, filepath.FromSlash(path))

	return groupDir
}

func Route(groupDir, path, fullPath string, handler interface{}) {
	fmt.Println("vim-go")
}

func addLink(parentGroupDir, title, href string) {
}
