// +build !cgo

package echo

// (ꈍᴗꈍ)
import (
	"strings"
)

func EchoHandler(input string) string {
	i := strings.ToLower(input)
	if strings.Contains(i, "uwu") || strings.Contains(input, "(ꈍᴗꈍ)") {
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
