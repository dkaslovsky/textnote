package template

import (
	"fmt"
	"strings"
)

// section is a named section of a Template
type section struct {
	name     string
	contents []contentItem
}

// newSection constructs a Section
func newSection(name string, items ...contentItem) *section {
	return &section{
		name:     name,
		contents: items,
	}
}

func (s *section) deleteContents() {
	s.contents = []contentItem{}
}

func (s *section) getNameString(prefix string, suffix string) string {
	return fmt.Sprintf("%s%s%s\n", prefix, s.name, suffix)
}

func (s *section) getContentString() string {
	str := ""
	for _, content := range s.contents {
		txt := content.string()
		if !strings.HasSuffix(txt, "\n") {
			txt += "\n"
		}
		str += txt
	}
	return str
}

type contentItem struct {
	header string
	text   string
}

func (ci contentItem) string() string {
	if ci.header != "" {
		return fmt.Sprintf("%s\n%s", ci.header, ci.text)
	}
	return ci.text
}
