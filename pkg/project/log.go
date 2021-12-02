package project

import (
	"io"

	"github.com/kyokomi/emoji/v2"
)

func LogSuccess(w io.Writer, str string, function string) {
	if str != "" {
		emoji.Fprintf(w, ":white_check_mark: | %s | %s\n", function, str)
	}
}

func LogFail(w io.Writer, str string, function string) {
	if str != "" {
		emoji.Fprintf(w, ":x: | %s | %s\n", function, str)
	}
}

func LogInfo(w io.Writer, str string, function string) {
	if str != "" {
		emoji.Fprintf(w, ":information: | %s | %s\n", function, str)
	}
}
