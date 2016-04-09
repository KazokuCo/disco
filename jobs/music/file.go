package music

import (
	"regexp"
)

func FindPlayFilename(s string) string {
	playRE := regexp.MustCompile(`play ([^,]+)`)
	if match := playRE.FindStringSubmatch(s); match != nil {
		return match[1]
	}
	return ""
}
