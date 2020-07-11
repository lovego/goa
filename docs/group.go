package docs

import (
	"bytes"
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

func (g *Group) Child(path, fullPath string, descs []string) Group {
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
		title := descs[0]
		if child.depth == 0 {
			title += " (" + fullPath + ")"
			child.CreateReadme(title, descs[1:])
		}
		g.LinkInReadme(title, path)
	}
	return child
}

func (g *Group) Route(method, path, fullPath string, handler interface{}) {
	var route Route
	if !route.Parse(handler) {
		return
	}

	path = pathPkg.Clean(path)
	switch path {
	case ".", "/":
		path = method + ".md"
	default:
		path = pathPkg.Join(pathPkg.Dir(path), method+"_"+pathPkg.Base(path)+".md")
	}

	file := filepath.Join(g.Dir, filepath.FromSlash(path))
	mkdir(filepath.Dir(file))
	if err := ioutil.WriteFile(file, route.Doc(method, fullPath), 0666); err != nil {
		log.Panic(err)
	}

	g.LinkInReadme(route.Title(method, fullPath), path)
}

func (g *Group) CreateReadme(title string, descs []string) {
	buf := bytes.NewBufferString("# " + title + "\n")
	for _, desc := range descs {
		buf.WriteString(desc + "\n\n")
	}

	mkdir(g.Dir)
	if err := ioutil.WriteFile(filepath.Join(g.Dir, readme), buf.Bytes(), 0666); err != nil {
		log.Panic(err)
	}
}

func (g *Group) LinkInReadme(title, href string) {
	buf := bytes.NewBufferString("##")
	if g.depth > 0 {
		buf.WriteString(strings.Repeat("#", g.depth))
	}
	buf.WriteString(" ")
	if g.depth > 0 {
		buf.WriteString(strings.Repeat("　", g.depth)) // use a full-width space
	}

	switch href {
	case ".", "/":
		buf.WriteString(title)
	default:
		href = pathPkg.Join(".", href)
		buf.WriteString("[" + title + "](" + href + ")")
	}
	buf.WriteByte('\n')

	mkdir(g.Dir)
	if err := fs.Append(filepath.Join(g.Dir, readme), buf.Bytes()); err != nil {
		log.Panic(err)
	}
}

func (g *Group) SetDir(dir string) {
	g.Dir = filepath.Clean(strings.TrimSpace(dir))
	if fs.Exist(g.Dir) {
		if err := os.RemoveAll(g.Dir); err != nil {
			log.Panic(err)
		}
	}
}

func cleanDescs(descs []string) (result []string) {
	for _, desc := range descs {
		if desc = strings.TrimSpace(desc); desc != "" {
			result = append(result, desc)
		}
	}
	return
}

func mkdir(dir string) {
	if fs.NotExist(dir) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Panic(err)
		}
	}
}
