package menu

import (
	"strings"
	"unicode/utf8"

	"github.com/haochend413/muninx/internal/models"
)

// matchesQuery reports whether note's content contains query, case-insensitively.
func matchesQuery(n *models.Note, query string) bool {
	return strings.Contains(strings.ToLower(n.Content), strings.ToLower(query))
}

// snippetAround returns content starting a little before the first
// case-insensitive occurrence of query, so the match stays visible after the
// table truncates the cell to its column width. The table's own truncation
// already appends a trailing "…" when it clips the end; here we only need to
// add a leading "..." when we clip the front.
func snippetAround(content, query string, leftContext int) string {
	lower := strings.ToLower(content)
	byteIdx := strings.Index(lower, strings.ToLower(query))
	if byteIdx < 0 {
		return content
	}

	matchStart := utf8.RuneCountInString(content[:byteIdx])
	start := matchStart - leftContext
	if start <= 0 {
		return content
	}

	runes := []rune(content)
	return "..." + string(runes[start:])
}
