package dict

import (
	"sort"
	"strconv"
	"strings"

	"github.com/angch/multibot/pkg/bothandler"
)

func init() {
	bothandler.RegisterCatchallHandler(DictHandler)
	myDict = NewMetaDictionary()
}

var myDict *MetaDictionary

func DictHandler(r bothandler.Request) string {
	input := r.Content
	args := strings.Split(input, " ")

	if len(args) == 0 || args[0] != "!dict" {
		return ""
	}

	if myDict == nil {
		return ""
	}
	w := myDict.All
	d := myDict
	sorttype := ""
	for _, arg := range args {
		if arg == "5" && w == d.All {
			w = d.Five
			continue
		}

		n, err := strconv.Atoi(arg)
		if err == nil && n > 0 {
			w = w.Len(n)
			continue
		}
		if len(arg) > 0 && arg[0] == '=' {
			w = w.Match(arg[1:])
		}
		if len(arg) > 0 && arg[0] == '+' {
			w = w.Contains(arg[1:])
		}
		if len(arg) > 0 && arg[0] == '-' {
			w = w.DoesNotContain(arg[1:])
		}
		if len(arg) > 0 && arg[0] == '~' {
			lc := CountLetters(strings.ToLower(arg[1:]))
			w = w.ContainsAll(lc)
		}
		if len(arg) > 0 && arg[0] == '|' {
			sorttype = arg[1:]
		}
	}

	o := make([]string, 0, len(w.Words))
	for k := range w.Words {
		o = append(o, k)
	}

	switch sorttype {
	case "len":
		sort.Slice(o, func(i, j int) bool { return len(o[i]) > len(o[j]) })
	default:
		sort.Strings(o)
	}
	if len(o) >= 10 {
		o = o[:9]
		o = append(o, "...")
	}

	return strings.Join(o, ", ")
}
