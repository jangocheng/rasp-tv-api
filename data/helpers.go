package data

import "strings"

func sqlEscape(str string) string {
	return strings.Replace(str, "'", "''", -1)
}
