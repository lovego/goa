package docs

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/url"
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
	path = cleanPath(path)
	child := Group{Dir: filepath.Join(g.Dir, filepath.FromSlash(path))}
	switch path {
	case ".", "/":
		child.depth = g.depth + 1
	default:
		child.depth = 0
		path = pathPkg.Join(path, readme)
	}

	descs = cleanDescs(descs)
	if len(descs) > 0 {
		title := descs[0]
		if child.depth == 0 {
			title += " (" + fullPath + ")"
			child.CreateReadme(title, descs[1:])
			g.LinkInReadme(title, path, nil, false)
		} else {
			g.LinkInReadme(title, path, descs[1:], false)
		}
	}
	return child
}

func (g *Group) Route(method, path, fullPath string, handler interface{}) {
	var route Route
	if !route.Parse(handler) {
		return
	}

	path = cleanPath(path)
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

	title := route.Title() + " ： " + route.MethodPath(method, fullPath)
	g.LinkInReadme(title, path, nil, true)
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

func (g *Group) LinkInReadme(title, href string, desc []string, isRoute bool) {
	buf := bytes.NewBufferString("##")
	if g.depth > 0 {
		buf.WriteString(strings.Repeat("#", g.depth))
	} else if isRoute {
		buf.WriteString("#")
	}
	buf.WriteString(" ")
	if g.depth > 0 {
		buf.WriteString(strings.Repeat("　" /* a full-width space */, g.depth))
	}

	switch href {
	case ".", "/":
		buf.WriteString(title)
	default:
		u := url.URL{Path: pathPkg.Join(".", href)}
		buf.WriteString("[" + title + "](" + u.EscapedPath() + ")")
	}
	buf.WriteByte('\n')
	for _, line := range desc {
		buf.WriteString(line + "\n")
	}

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

func cleanPath(path string) string {
	const chars = `\|:*?"<>` // characters invalid for windows file name
	if strings.ContainsAny(path, chars) {
		components := strings.Split(path, "/")
		for i, v := range components {
			if strings.ContainsAny(v, chars) {
				var bytesSlice []byte
				for _, b := range md5.Sum([]byte(v)) {
					bytesSlice = append(bytesSlice, b)
				}
				components[i] = base64.RawURLEncoding.EncodeToString(bytesSlice)
			}
		}
		path = pathPkg.Join(components...)
	}
	return pathPkg.Clean(path)
}
