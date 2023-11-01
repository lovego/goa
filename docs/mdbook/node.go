package mdbook

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"
)

var pwd = func() string {
	dir, _ := os.Getwd()
	return dir
}()

type Line struct {
	Title        string
	BasePath     string
	Path         string
	IsNumber     bool
	level        int
	IsSeparators bool
}
type LineList []*Line

type LineNode struct {
	Line
	LineNodes *LineNodes
}

type LineNodes struct {
	List []*LineNode
}

func (a *LineNodes) Sort() {
	sort.Slice(a.List, func(i, j int) bool {
		return a.List[i].Line.Path < a.List[j].Line.Path
	})
}

func (a *LineNodes) Markdown(buf *bytes.Buffer, level int, basePath string) {
	if a == nil {
		return
	}

	a.Sort()

	for _, node := range a.List {
		node.level = level
		node.BasePath = basePath
		node.Markdown(buf)
	}
}

func (a *LineNodes) GetNodes(basePath string) (list *LineNodes) {
	if a == nil {
		return nil
	}

	node := a.GetNode(basePath)
	if node != nil {
		return a
	}
	return nil

}
func (a *LineNodes) GetNode(basePath string) *LineNode {
	if a == nil {
		return nil
	}

	for _, node := range a.List {
		if len(node.BasePath) > len(basePath) {
			return nil
		}
		if node.BasePath == basePath {
			return node
		}
		if node.LineNodes != nil {
			if n := node.LineNodes.GetNode(basePath); n != nil {
				return n
			}
		}
	}

	return nil
}

func (a *Line) Markdown(buf *bytes.Buffer) {
	if a == nil {
		return
	}

	buf.WriteString(a.LineString())
}

func (a *LineNode) Markdown(buf *bytes.Buffer) {
	if a == nil {
		return
	}

	a.Line.Markdown(buf)
	//for _, node := range a.LineNodes {
	//	node.Level = a.Level + 1
	//	node.Markdown(buf)
	//}

	a.LineNodes.Markdown(buf, a.level+1, a.BasePath)

}

func (a *Line) LineString() string {
	if a == nil {
		return ""
	}

	if a.IsSeparators {
		return "\n\n---\n\n\n"
	}

	if a.Title == "" && a.Path == "" {
		return ""
	}

	s := ""
	if a.IsNumber {
		s = "- "
	}

	return fmt.Sprintf("%s%s[%s](%s)\n", a.LevelSpace(), s, a.Title, a.GetPath())
}

func (a *Line) GetPath() string {
	return path.Join(a.BasePath, strings.TrimPrefix(a.Path, pwd))
}

func (a *Line) LevelSpace() string {
	p := strings.TrimPrefix(a.Path, pwd)

	a.level = len(strings.Split(p, "/"))

	if a.level == 0 {
		return ""
	}

	buf := new(bytes.Buffer)

	num := 2 * a.level

	for i := 0; i < num; i++ {
		buf.WriteByte(' ')
	}
	return buf.String()
}
