package docs

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	webpath "path"
	"path/filepath"
	"strings"

	"github.com/lovego/fs"
)

const readme = "README.md"

// Group for docs generation. Every docs.Group has an one-to-one goa.RouterGroup.
// A group can be inline, that means its docs is generated in its parent's REAMED file.
type Group struct {
	// Readme file directory of this group, converted to hash if has special chars.
	// readmeDir is parent group's readmeDir if this group is inline.
	readmeDir string

	inlineHashPath string // Hash path prefix for child group or route if this group is inline.
	inlineDepth    int    // Indent depth in README file if this group is inline.
}

func (g *Group) Child(path, fullPath string, descs []string) Group {
	path = webpath.Clean(path)
	var inline = isInlinePath(path) || isInlinePath(fullPath)

	var hashPath = makeHashPath(path)
	var child Group
	if inline {
		child.readmeDir = g.readmeDir
		child.inlineHashPath = webpath.Join(g.inlineHashPath, hashPath)
		child.inlineDepth = g.inlineDepth + 1
	} else {
		child.readmeDir = filepath.Join(g.readmeDir, g.inlineHashPath, filepath.FromSlash(hashPath))
	}

	descs = cleanDescs(descs)
	if len(descs) > 0 {
		title := descs[0]
		if inline {
			g.LinkInReadme(title, "", descs[1:], false)
		} else {
			title += " (" + fullPath + ")"
			child.CreateReadme(title, descs[1:])
			g.LinkInReadme(title, webpath.Join(hashPath, readme), nil, false)
		}
	}
	return child
}

func (g *Group) Route(method, path, fullPath string, handler interface{}) {
	var route Route
	if !route.Parse(handler) {
		return
	}

	var hashPath = webpath.Join(g.inlineHashPath, makeHashPath(path))
	if isInlinePath(hashPath) {
		hashPath = method + ".md"
	} else {
		hashPath = webpath.Join(webpath.Dir(hashPath), method+"_"+webpath.Base(hashPath)+".md")
	}

	var file = filepath.Join(g.readmeDir, filepath.FromSlash(hashPath))
	mkdir(filepath.Dir(file))
	if err := ioutil.WriteFile(file, route.Doc(method, fullPath), 0666); err != nil {
		log.Panic(err)
	}

	title := route.Title() + " ： " + route.MethodPath(method, fullPath)
	g.LinkInReadme(title, hashPath, nil, true)
}

func (g *Group) CreateReadme(title string, descs []string) {
	buf := bytes.NewBufferString("# " + title + "\n")
	for _, desc := range descs {
		buf.WriteString(desc + "\n\n")
	}

	mkdir(g.readmeDir)
	if err := ioutil.WriteFile(filepath.Join(g.readmeDir, readme), buf.Bytes(), 0666); err != nil {
		log.Panic(err)
	}
}

func (g *Group) LinkInReadme(title, href string, desc []string, isRoute bool) {
	buf := bytes.NewBufferString("##")
	if g.inlineDepth > 0 {
		buf.WriteString(strings.Repeat("#", g.inlineDepth))
	} else if isRoute {
		buf.WriteString("#")
	}
	buf.WriteString(" ")
	if g.inlineDepth > 0 {
		// indent by full-width space
		buf.WriteString(strings.Repeat("　", g.inlineDepth))
	}

	if isInlinePath(href) {
		buf.WriteString(title)
	} else {
		u := url.URL{Path: webpath.Join(".", href)}
		buf.WriteString("[" + title + "](" + u.EscapedPath() + ")")
	}
	buf.WriteByte('\n')
	for _, line := range desc {
		buf.WriteString(line + "\n")
	}

	mkdir(g.readmeDir)
	if err := fs.Append(filepath.Join(g.readmeDir, readme), buf.Bytes()); err != nil {
		log.Panic(err)
	}
}

func (g *Group) Valid() bool {
	return g.readmeDir != ""
}

func (g *Group) SetDir(dir string) {
	g.readmeDir = filepath.Clean(strings.TrimSpace(dir))
	if fs.Exist(g.readmeDir) {
		if err := os.RemoveAll(g.readmeDir); err != nil {
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

func isInlinePath(path string) bool {
	switch path {
	case "", ".", "/": // thease also in factly makes an inline group.
		return true
	}
	return false
}

// convert path with sepecial chars to hash
func makeHashPath(path string) string {
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
		path = webpath.Join(components...)
	}
	return webpath.Clean(path)
}
