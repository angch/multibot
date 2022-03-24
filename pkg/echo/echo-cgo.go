//go:build cgo
// +build cgo

package echo

// #cgo LDFLAGS: -L../../lib/uwu/target/release -luwu
// #include "../../lib/uwu.h"
import "C"

// (ꈍᴗꈍ)
import (
	"strings"
)

func uwucheck(i string) string {
	if strings.Contains(i, "uwu") || strings.Contains(i, "(ꈍᴗꈍ)") {
		return strings.Replace(C.GoString(C.uwuify(C.CString(i))), "uwu", "(ꈍᴗꈍ)", 1)
	}
	return ""
}
