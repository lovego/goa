package mdbook

import (
	"fmt"
	"testing"
)

func TestBoolTest(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{
			name: "markdown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BookTest()
		})
	}
}

func BookTest() {
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
