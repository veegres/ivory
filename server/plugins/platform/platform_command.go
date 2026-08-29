package platform

import (
	"strings"
	"unicode"
)

// SplitCommand turns a command as the user wrote it into the arguments it will
// run as, honoring single/double-quoted spans the way a real shell would - so
// a flag value like `-e SCOPE="my cluster"` stays one argument instead of
// being broken apart at the space inside the quotes, and the newlines that
// make a command readable collapse while those inside a quoted span survive as
// the statement separators they are.
//
// It is here rather than inside an adapter because a caller has to split
// before it interpolates: a value filled into an argument that has already
// been separated can never introduce an argument boundary, close a span the
// template author opened, or be read as syntax - which is what makes escaping
// unnecessary. Adapters receive the finished arguments.
//
// Every character it sees belongs to the template author, so backslashes
// follow ordinary shell rules: inside a single-quoted span a backslash is a
// literal character, so a post script's own `\"` survives tokenizing and
// reaches the inner `sh -c` still escaped.
func SplitCommand(value string) []string {
	fields := make([]string, 0)
	var current strings.Builder
	hasToken := false
	var quote rune
	runes := []rune(value)
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		switch {
		case r == '\\' && quote != '\'' && i+1 < len(runes) && escapesInQuote(runes[i+1], quote):
			i++
			current.WriteRune(runes[i])
			hasToken = true
		case quote != 0:
			if r == quote {
				quote = 0
			} else {
				current.WriteRune(r)
			}
		case r == '\'' || r == '"':
			quote = r
			hasToken = true
		case unicode.IsSpace(r):
			if hasToken {
				fields = append(fields, current.String())
				current.Reset()
				hasToken = false
			}
		default:
			current.WriteRune(r)
			hasToken = true
		}
	}
	if hasToken {
		fields = append(fields, current.String())
	}
	return fields
}

// escapesInQuote reports whether a backslash before r actually escapes it in
// the given quote state. Unquoted, a backslash escapes anything; inside a
// double-quoted span it only escapes the few characters that span still
// interprets, and is an ordinary character before anything else.
func escapesInQuote(r rune, quote rune) bool {
	if quote == '"' {
		return r == '"' || r == '\\' || r == '$' || r == '`'
	}
	return true
}
