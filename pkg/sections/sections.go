package sections

import (
	"fmt"
	"regexp"
	"strings"
)

func Extract(source, section string, sections []string) (string, error) {
	nextSection := ""
	for i, sec := range sections {
		if sec == section && i+1 < len(sections) {
			nextSection = sections[i+1]
			break
		}
	}

	regexPattern := ""
	if nextSection == "" {
		regexPattern = `(?s)### ` + regexp.QuoteMeta(section) + `(.*?)$`
	} else {
		regexPattern = `(?s)### ` + regexp.QuoteMeta(section) + `(.*?)### ` + regexp.QuoteMeta(nextSection)
	}

	re := regexp.MustCompile(regexPattern)
	matches := re.FindStringSubmatch(source)
	if len(matches) < 2 {
		return "", fmt.Errorf("section %q not found", section)
	}
	return strings.TrimSpace(matches[1]), nil
}
