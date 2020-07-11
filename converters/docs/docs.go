package docs

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	pathPkg "path"
	"path/filepath"
	"strings"

	"github.com/lovego/fs"
)

type Group struct {
	Dir   string
	depth int // depth in readme file of parent group.
}

const readme = "README.md"

func (g *Group) SetDir(dir string) {
	g.Dir = filepath.Clean(strings.TrimSpace(dir))
	if err := os.Remove(g.Dir); err != nil {
		log.Panic(err)
	}
}

func (g *Group) Child(path string, descs []string) Group {
	path = pathPkg.Clean(path)
	child := Group{Dir: filepath.Join(g.Dir, filepath.FromSlash(path))}
	switch path {
	case ".", "/":
		child.depth = g.depth + 1
	default:
		child.depth = 0
	}

	descs = cleanDescs(descs)
	if len(descs) > 0 {
		if child.depth == 0 {
			child.CreateReadme(descs)
		}
		g.LinkInReadme(child.depth, descs[0], path)
	}
	return child
}

func (g *Group) CreateReadme(descs []string) {
	content := "# " + descs[0] + "\n"
	for i := 1; i < len(descs); i++ {
		content += descs[i] + "<br>\n"
	}

	if err := ioutil.WriteFile(filepath.Join(g.Dir, readme), []byte(content), 0666); err != nil {
		log.Panic(err)
	}
}

func (g *Group) LinkInReadme(depth int, title, href string) {
	content := "##"
	if depth >= 2 {
		content += strings.Repeat("#", depth-1)
	}
	content += " "

	switch href {
	case ".", "/":
		content += title
	default:
		href = filepath.Join(".", filepath.FromSlash(href))
		content += "[" + title + "](" + href + ")"
	}
	content += "\n"

	if err := fs.Append(filepath.Join(g.Dir, readme), []byte(content)); err != nil {
		log.Panic(err)
	}
}
func (g *Group) Route(path, fullPath string, handler interface{}) {
	fmt.Println("vim-go")
}

func cleanDescs(descs []string) (result []string) {
	for _, desc := range descs {
		if desc = strings.TrimSpace(desc); desc != "" {
			result = append(result, desc)
		}
	}
	return
}
