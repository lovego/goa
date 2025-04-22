package mdbook

import (
	"bytes"
	"fmt"
)

type PartChapter struct {
	Title string
	Parts *LineNodes
}

func (a *PartChapter) Markdown(buf *bytes.Buffer, level int, basePath string) {
	if len(a.Title) > 0 {
		buf.WriteString(fmt.Sprintf("# %s\n", a.Title))
	}
	a.Parts.Markdown(buf, level, basePath)
}

type Book struct {
	Title           string
	BasePath        string
	PrefixChapter   *LineNodes
	PartChapters    []*PartChapter
	NumberedChapter *LineNodes
	SuffixChapter   *LineNodes
	DraftChapters   *LineNodes
}

type PartChapters []*PartChapter

func (a *Book) Markdown() *bytes.Buffer {
	buf := new(bytes.Buffer)

	a.PrefixChapter.Markdown(buf, 0, a.BasePath)

	for _, chapter := range a.PartChapters {
		chapter.Markdown(buf, 0, a.BasePath)
	}

	a.NumberedChapter.Markdown(buf, 0, a.BasePath)
	a.SuffixChapter.Markdown(buf, 0, a.BasePath)

	return buf
}

func (a *Book) AddPart(title string) {
	part := &PartChapter{
		Title: title,
		Parts: &LineNodes{},
	}

	a.PartChapters = append(a.PartChapters, part)

}

func (a *Book) AddNode(basePath, fullPath, title string) {
	if len(a.PartChapters) == 0 {
		return
	}
	line := Line{
		Title:    title,
		BasePath: basePath,
		Path:     fullPath,
		IsNumber: true,
	}

	a.PartChapters[0].Parts.List = append(a.PartChapters[0].Parts.List, &LineNode{
		Line:      line,
		LineNodes: nil,
	})

}
