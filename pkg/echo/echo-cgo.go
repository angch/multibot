//go:build cgo
// +build cgo

package echo

// #cgo LDFLAGS: -L../../lib/uwu/target/release -luwu
// #include "../../lib/uwu.h"
import "C"

// (ꈍᴗꈍ)
import (
	"strings"

	"github.com/angch/discordbot/pkg/bothandler"
)

func EchoHandler(request bothandler.Request) string {
	i := strings.ToLower(request.Content)
	if strings.Contains(i, "uwu") || strings.Contains(i, "(ꈍᴗꈍ)") {
		return strings.Replace(C.GoString(C.uwuify(C.CString(i))), "uwu", "(ꈍᴗꈍ)", 1)
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
