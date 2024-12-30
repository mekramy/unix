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
// The pattern string can contain placeholders for various styles and colors, which will be replaced by the corresponding ANSI codes.
//
// Supported styling patterns:
// {R}, @: RESET
// {B}, @B: BOLD
// {U}, @U: UNDERLINE
// {S}, @S: STRIKE
// {I}, @I: ITALIC
//
// Supported color patterns:
// {r}, @r: RED
// {g}, @g: GREEN
// {y}, @y: YELLOW
// {b}, @b: BLUE
// {p}, @p: PURPLE
// {c}, @c: CYAN
// {m}, @m: GRAY
// {w}, @w: WHITE
//
// Example usage:
// Formatter("{B}Bold Text{R} and {r}Red Text{R}\n")
// Formatter("{g}Green Text{R} with arguments: %d, %s\n", 42, "example")
//
// Arguments:
// - pattern: The string containing the text and placeholders for styling and colors.
// - args: The arguments to be formatted into the pattern string.
func PrintF(pattern string, args ...any) {
	replacer := strings.NewReplacer(
		"{R}", RESET,
		"@", RESET,
		"{B}", BOLD,
		"@B", BOLD,
		"{U}", UNDERLINE,
		"@U", UNDERLINE,
		"{S}", STRIKE,
		"@S", STRIKE,
		"{I}", ITALIC,
		"@I", ITALIC,
		"{r}", RED,
		"@r", RED,
		"{g}", GREEN,
		"@g", GREEN,
		"{y}", YELLOW,
		"@y", YELLOW,
		"{b}", BLUE,
		"@b", BLUE,
		"{p}", PURPLE,
		"@p", PURPLE,
		"{c}", CYAN,
		"@c", CYAN,
		"{m}", GRAY,
		"@m", GRAY,
		"{w}", WHITE,
		"@w", WHITE,
	)
	fmt.Printf(replacer.Replace(pattern), args...)
}
