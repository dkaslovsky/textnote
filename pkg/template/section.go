package template

import (
	"fmt"
	"strings"
)

// section is a named section of a Template
type section struct {
	name     string
	contents []string
}

// newSection constructs a Section
func newSection(name string, contents ...string) *section {
	return &section{
		name:     name,
		contents: contents,
	}
}

func (s *section) appendContents(contents string) {
	s.contents = append(s.contents, contents)
}

func (s *section) deleteContents() {
	s.contents = []string{}
}

func (s *section) getNameString(prefix string, suffix string) string {
	return fmt.Sprintf("%s%s%s\n", prefix, s.name, suffix)
}

func (s *section) getBodyString() string {
	body := ""
	for _, content := range s.contents {
		if !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		body += content
	}
	return body
}
