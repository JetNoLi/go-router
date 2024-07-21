package utils

import (
	"regexp"
)

func GetStringOccurrences(regexString string, value string) int {
	regExp := regexp.MustCompile(regexString)

	matches := regExp.FindAllStringIndex(value, -1)

	return len(matches)
}
