package unix_test

import (
	"testing"

	"github.com/mekramy/unix"
)

func TestCronJob(t *testing.T) {
	job := unix.NewCronJob("do some").
		Weekly(unix.Friday).
		SetTz(3, 30)

	if job.Compile() != "20 -3 * * 5" {
		t.Fatal("FAIL", job.Compile())
	}
}

func TestPrintF(t *testing.T) {
	tests := []struct {
		format string
		args   []any
	}{
		{"@B{Bold Text %s} @Um{%d}\n", []any{"test", 123}},
		{"@U{Underline Text%s} @Um{%d}\n", []any{"test", 123}},
		{"@S{Strike Text %s} @Um{%d}\n", []any{"test", 123}},
		{"@I{Italic Text %s} @Um{%d}\n", []any{"test", 123}},
		{"@r{Red Text %s} @Um{%d}\n", []any{"test", 123}},
		{"@g{Green Text %s} @Um{%d}\n", []any{"test", 123}},
		{"@y{Yellow Text %s} @Um{%d}\n", []any{"test", 123}},
		{"@b{Blue Text %s} @Um{%d}\n", []any{"test", 123}},
		{"@p{Purple Text %s} @Um{%d}\n", []any{"test", 123}},
		{"@c{Cyan Text %s} @UrS{%d}\n", []any{"test", 123}},
		{"@m{Gray Text %s} @Um{%d}\n", []any{"test", 123}},
		{"@w{White Text %s} \\@Um{%d Escaped Text}\n", []any{"test", 123}},
	}

	for _, tt := range tests {
		unix.PrintF(tt.format, tt.args...)
	}
}
