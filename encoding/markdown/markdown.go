package markdown

import (
	"fmt"
	"strings"
)

func ToYAML(raw string) (string, error) {
	var content strings.Builder
	var err error
	fmt.Printf("contains %d lines\n", strings.Count(raw, "\n"))
	lines := strings.Split(raw, "\n")
	var inBody bool
	for i, line := range lines {
		// remove first delimiter
		if i == 0 {
			fmt.Printf("skipping line %d, %s\n", i, line)
		} else {
			if !inBody {
				// this is last delimiter
				// replace with 'body: |' and
				// indent the rest of the body by 2 spaces
				if line == "---" {
					content.WriteString("body_raw: |")
					content.WriteString("\n")
					inBody = true
				} else {
					content.WriteString(line)
					content.WriteString("\n")
				}
			} else {
				content.WriteString("  ")
				content.WriteString(line)
				content.WriteString("\n")
			}
		}

	}
	return content.String(), err
}

/*
--- first delimiter
key: value
key2: value
---
body text

===========

key: value
key2: value
body: |
  body text
  next line of body text
*/
