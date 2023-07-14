package cmd

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/angch/multibot/pkg/dict"
	"github.com/spf13/cobra"
)

// dictCmd represents the dict command
var dictCmd = &cobra.Command{
	Use:   "dict",
	Short: "dict",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		d := dict.NewMetaDictionary()

		w := d.All
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
				if strings.Contains(arg, "/") {
					w2 := w.Multimatch(arg[1:])
					w = dict.DedupeDictionaries(w2)
				} else {
					w = w.Match(arg[1:])
				}
			}
			if len(arg) > 0 && arg[0] == '+' {
				w = w.Contains(arg[1:])
			}
			if len(arg) > 0 && arg[0] == '-' {
				w = w.DoesNotContain(arg[1:])
			}
			if len(arg) > 0 && arg[0] == '~' {
				lc := dict.CountLetters(strings.ToLower(arg[1:]))
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
		// if len(o) >= 10 {
		// 	o = o[:9]
		// 	o = append(o, "...")
		// }
		fmt.Println(strings.Join(o, "\n"))
	},
}

func init() {
	rootCmd.AddCommand(dictCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dictCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dictCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
