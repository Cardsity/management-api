package utils

import "regexp"

var blankRegex = regexp.MustCompile("(____+)+")

// Returns the amount of blanks found in a text. A blank counts as a blank when it consists of 4 or more underscores (_)
// that are not separated. The blank ends after the first character that is not an underscore. When no blank was found
// the blank count is 1. This is because it will likely be a question or something similar.
func GetBlankCount(text string) int {
	if c := len(blankRegex.FindAllString(text, -1)); c > 0 {
		return c
	} else {
		return 1
	}
}
