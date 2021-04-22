package echo

// +build cgo

// #cgo LDFLAGS: -L../../lib -luwu
// #include "../../lib/uwu.h"
import "C"

// (ꈍᴗꈍ)
import (
	"strings"
)

func EchoHandler(input string) string {
	i := strings.ToLower(input)
	if strings.Contains(i, "uwu") || strings.Contains(input, "(ꈍᴗꈍ)") {
		return strings.Replace(C.GoString(C.uwuify(C.CString(input))), "uwu", "(ꈍᴗꈍ)", 1)
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
