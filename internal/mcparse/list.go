package mcparse

import (
	"strings"
)

// ParsePlayersNames returns a slice of nicknames and bool, if the string does not match the result of the list command.
func ParsePlayersNames(resp string) ([]string, bool) {
	const header = "There are "
	if !strings.HasPrefix(resp, header) {
		return nil, false
	}

	colon := strings.IndexRune(resp, ':')
	if colon == -1 {
		return nil, false
	}

	tail := strings.TrimSpace(resp[colon+1:])

	if tail == "" {
		return nil, true
	}

	parts := strings.FieldsFunc(tail, func(r rune) bool {
		return r == ',' || r == '\n' || r == '\r'
	})

	var names []string
	for _, p := range parts {
		if n := strings.TrimSpace(p); n != "" {
			names = append(names, n)
		}
	}
	return names, true
}
