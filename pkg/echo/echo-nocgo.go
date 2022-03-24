//go:build !cgo
// +build !cgo

package echo

// (ꈍᴗꈍ)
import (
	"strings"
)

func uwucheck(i string) string {
	if strings.Contains(i, "uwu") || strings.Contains(i, "(ꈍᴗꈍ)") {
		return "(ꈍᴗꈍ)"
	}
}
