package mdbook

var DefaultBook = Book{
	Title:         "默认标题",
	BasePath:      "",
	PrefixChapter: &LineNodes{},
	//PartChapters:  []*PartChapter{{}},
}

func SetBook(b Book) {
	DefaultBook = b
}
