package echo

// #cgo LDFLAGS: -L../../lib -luwu
// #include "../../lib/uwu.h"
import "C"

import (
	"strings"

	"github.com/angch/discordbot/pkg/bothandler"
)

func init() {
	bothandler.RegisterCatchallHandler(UwuHandler)
}

func UwuHandler(input string) string {
	if strings.Contains(strings.ToLower(input), "uwu") {
		return C.GoString(C.uwuify(C.CString(input)))
	}

	return ""
}
