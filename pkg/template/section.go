package template

import (
	"fmt"
	"strings"
)

// section is a named section of a Template
type section struct {
	name     string
	contents []item
}

// newSection constructs a Section
func newSection(name string, items ...item) *section {
	return &section{
		name:     name,
		contents: items,
	}
}

func (s *section) deleteContents() {
	s.contents = []item{}
}

func (s *section) getNameString(prefix string, suffix string) string {
	return fmt.Sprintf("%s%s%s\n", prefix, s.name, suffix)
}

func (s *section) getContentString() string {
	str := ""
	for _, content := range s.contents {
		txt := content.text
		if !strings.HasSuffix(txt, "\n") {
			txt += "\n"
		}
		str += txt
	}
	return str
}

type item struct {
	header string
	text   string
}
