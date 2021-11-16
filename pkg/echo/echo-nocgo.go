//go:build !cgo
// +build !cgo

package echo

// (ꈍᴗꈍ)
import (
	"strings"
)

func EchoHandler(request bothandler.Request) string {
	i := strings.ToLower(request.Content)
	if strings.Contains(i, "uwu") || strings.Contains(i, "(ꈍᴗꈍ)") {
		return "(ꈍᴗꈍ)"
	}

	r, ok := echos[i]
	if ok {
		return r
	}

	for _, v := range fragments {
		if strings.Contains(i, v.From) {
			return v.To
		}
	}
	return ""
}
