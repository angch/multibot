package spacetraders

func removeEmptyStrings(words []string) []string {
	var ret []string
	for _, w := range words {
		if w != "" {
			ret = append(ret, w)
		}
	}
	return ret
}
