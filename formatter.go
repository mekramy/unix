package unix

import (
	"fmt"
	"strings"
)

const (
	RESET     = "\033[0m"
	BOLD      = "\033[1m"
	UNDERLINE = "\033[4m"
	STRIKE    = "\033[9m"
	ITALIC    = "\033[3m"
)

const (
	RED    = "\033[31m"
	GREEN  = "\033[32m"
	YELLOW = "\033[33m"
	BLUE   = "\033[34m"
	PURPLE = "\033[35m"
	CYAN   = "\033[36m"
	GRAY   = "\033[37m"
	WHITE  = "\033[37m"
)

// Printf applies ANSI escape codes to a given pattern string and prints the formatted string with the provided arguments.
// The pattern string can contain tags followed by content token for various styles and colors, which will be replaced by the corresponding ANSI codes.
// You can scape token with \@
//
// Supported styling patterns:
// B: BOLD
// U: UNDERLINE
// S: STRIKE
// I: ITALIC
//
// Supported color patterns:
// r: RED
// g: GREEN
// y: YELLOW
// b: BLUE
// p: PURPLE
// c: CYAN
// m: GRAY
// w: WHITE
//
// Example usage:
// PrintF("@Bg{Bold Green Text} and @r{Red %s}\n", "message")
//
// Arguments:
// - format: The string containing the standard go fmt format with styled tokens.
// - args: The arguments to be passed into the format string.
func PrintF(format string, args ...any) {
	var started bool
	var inside bool
	var result strings.Builder
	var token strings.Builder
	var content strings.Builder
	replacer := strings.NewReplacer(
		// styles
		"B", BOLD,
		"U", UNDERLINE,
		"S", STRIKE,
		"I", ITALIC,
		// colors
		"r", RED,
		"g", GREEN,
		"y", YELLOW,
		"b", BLUE,
		"p", PURPLE,
		"c", CYAN,
		"m", GRAY,
		"w", WHITE,
	)

	for i := 0; i < len(format); i++ {
		if format[i] == '@' && (i == 0 || format[i-1] != '\\') {
			started = true
			token.Reset()
			content.Reset()
		} else if started && !inside {
			if format[i] == '{' && (i == 0 || format[i-1] != '\\') {
				inside = true
				content.Reset()
			} else {
				token.WriteByte(format[i])
			}
		} else if started && inside {
			if format[i] == '}' && (i == 0 || format[i-1] != '\\') {
				started = false
				inside = false
				result.WriteString(replacer.Replace(token.String()))
				result.WriteString(content.String())
				result.WriteString(RESET)
			} else {
				content.WriteByte(format[i])
			}
		} else {
			result.WriteByte(format[i])
		}
	}
	fmt.Printf(strings.NewReplacer("\\@", "@").Replace(result.String()), args...)
}
