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

func BooKTest() {
	book := Book{
		Title:    "book",
		BasePath: "base/guide/",
		PrefixChapter: &LineNodes{List: []*LineNode{
			{Line: Line{
				Title: "测试",
				Path:  "/test/test.md",
			}},
			{Line: Line{
				IsSeparators: true,
			}},
			{Line: Line{
				Title: "测试2",
				Path:  "/test/test2.md",
			}},
		}},
		PartChapters: []*PartChapter{
			{
				Title: "用户指导",
				Parts: &LineNodes{List: []*LineNode{
					{Line: Line{
						Title:    "安装",
						Path:     "/guide/install.md",
						IsNumber: true,
					}},
					{Line: Line{
						Title:    "简介",
						Path:     "/guide/install.md",
						IsNumber: true,
					}},
					{Line: Line{
						Title:        "版本",
						Path:         "/guide/version.md",
						IsNumber:     true,
						IsSeparators: false,
					},
						LineNodes: &LineNodes{List: []*LineNode{
							{Line: Line{
								Title:    "linux",
								Path:     "/guide/install.md",
								IsNumber: true,
							}},
							{Line: Line{
								Title:    "macos",
								Path:     "/guide/install.md",
								IsNumber: true,
							}},
						},
						},
					}},
				}},
			{
				Title: "erp saas项目",
				Parts: &LineNodes{List: []*LineNode{
					{Line: Line{
						Title:    "安装",
						Path:     "/guide/install.md",
						IsNumber: true,
					}},
					{Line: Line{
						Title:    "简介",
						Path:     "/guide/install.md",
						IsNumber: true,
					}},
					{Line: Line{
						IsSeparators: true,
					}},
					{Line: Line{
						Title:        "版本",
						Path:         "/guide/version.md",
						IsNumber:     true,
						IsSeparators: false,
					},
						LineNodes: &LineNodes{List: []*LineNode{
							{Line: Line{
								Title:    "linux",
								Path:     "/guide/install.md",
								IsNumber: true,
							},
							},
							{Line: Line{
								Title:    "macos",
								Path:     "/guide/install.md",
								IsNumber: true,
							}},
						}},
					},
				},
				}},
		},
		NumberedChapter: &LineNodes{List: []*LineNode{
			{Line: Line{
				Title:    "安装",
				Path:     "/guide/install.md",
				IsNumber: true,
			}},
			{Line: Line{
				Title:    "简介",
				Path:     "/guide/install.md",
				IsNumber: true,
			}},
			{Line: Line{
				IsSeparators: true,
			}},
			{Line: Line{
				Title:    "版本",
				Path:     "/guide/version.md",
				IsNumber: true,
			},
				LineNodes: &LineNodes{List: []*LineNode{
					{Line: Line{
						Title:    "linux",
						Path:     "/guide/install.md",
						IsNumber: true,
					},
					},
					{Line: Line{
						Title:    "macos",
						Path:     "/guide/install.md",
						IsNumber: true,
					}},
				}},
			},
		}},
		SuffixChapter: &LineNodes{List: []*LineNode{
			{Line: Line{
				Title: "测试",
				Path:  "/test/test.md",
			}},
			{Line: Line{
				IsSeparators: true,
			}},
			{Line: Line{
				Title: "测试2",
				Path:  "/test/test2.md",
			}},
		}},
	}

	buf := book.Markdown()

	fmt.Println(buf)

}
