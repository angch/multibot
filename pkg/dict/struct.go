package dict

import (
	"bufio"
	"os"
	"sort"
	"strings"
)

type LetterCounts map[byte]byte
type Dictionary struct {
	Words  map[string]bool
	Counts map[string]LetterCounts
}

func (d *Dictionary) UpdateCounts() {
	if d.Counts != nil {
		return
	}
	d.Counts = make(map[string]LetterCounts)
	for k := range d.Words {
		for i := 0; i < len(k); i++ {
			if _, ok := d.Counts[k]; !ok {
				d.Counts[k] = make(LetterCounts)
			}
			d.Counts[k][k[i]]++
		}
	}
}

type MetaDictionary struct {
	All  *Dictionary
	Five *Dictionary
}

func NewMetaDictionary() *MetaDictionary {
	w := make(map[string]bool)
	f := make(map[string]bool)

	filepath := "/usr/share/dict/words"
	file, err := os.Open(filepath)
	if err != nil {
		return nil
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

s:
	for scanner.Scan() {
		t := strings.ToLower(scanner.Text())

		for _, v := range t {
			if !(v >= 'a' && v <= 'z') {
				continue s
			}
		}
		w[t] = true

		if len(t) == 5 {
			f[t] = true
		}
	}
	meta := MetaDictionary{}
	all := &Dictionary{
		Words: w,
	}
	five := &Dictionary{
		Words: f,
	}
	meta.All = all
	meta.Five = five

	return &meta
}

func (d *Dictionary) Match(s string) *Dictionary {
	out := make(map[string]bool)
o:
	for k := range d.Words {
		if len(k) < len(s) {
			continue
		}

		for k2, v2 := range s {
			if v2 == '.' {
				continue
			}
			if v2 >= 'A' && v2 <= 'Z' {
				v3 := v2 - 'A' + 'a'
				if k[k2] == byte(v3) {
					continue o
				}
			} else {
				if k[k2] != byte(v2) {
					continue o
				}
			}
		}
		out[k] = true
	}
	return &Dictionary{Words: out}
}

func (d *Dictionary) Contains(s string) *Dictionary {
	out := make(map[string]bool)
o:
	for k := range d.Words {
		if len(k) < len(s) {
			continue
		}

	o2:
		for _, v2 := range s {
			if v2 >= 'A' && v2 <= 'Z' {
				v3 := v2 - 'A' + 'a'
				for _, k3 := range k {
					if k3 == v3 {
						continue o
					}
				}
				continue o2
			} else {
				for _, k3 := range k {
					if k3 == v2 {
						continue o2
					}
				}
				continue o
			}
		}
		out[k] = true
	}
	return &Dictionary{Words: out}
}

func (d *Dictionary) DoesNotContain(s string) *Dictionary {
	out := make(map[string]bool)
o:
	for k := range d.Words {
		if len(k) < len(s) {
			continue
		}

		// o2:
		for _, v2 := range s {
			if v2 == '.' {
				continue
			}

			for _, k3 := range k {
				if k3 == v2 {
					continue o
				}
			}
		}
		out[k] = true
	}
	return &Dictionary{Words: out}
}

func (d *Dictionary) Len(n int) *Dictionary {
	out := make(map[string]bool)
	for k := range d.Words {
		if len(k) == n {
			out[k] = true
		}
	}
	return &Dictionary{Words: out}
}

func CountLetters(s string) LetterCounts {
	count := make(LetterCounts)
	for _, c := range s {
		count[byte(c)]++
	}
	return count
}

// ContainsAll returns a new Dictionary containing all words that contain all premutation of the words
func (d *Dictionary) ContainsAll(count LetterCounts) *Dictionary {
	d.UpdateCounts()
	out := make(map[string]bool)
a:
	for k, v := range d.Counts {
		for k2, v2 := range v {
			if byte(v2) > count[k2] {
				continue a
			}
		}
		out[k] = true
	}

	return &Dictionary{Words: out}
}

func (d *Dictionary) Multimatch(pattern string) []*Dictionary {
	lr := strings.Split(pattern, "/")
	if len(lr) != 2 {
		return []*Dictionary{d}
	}
	templates := strings.Split(lr[0], ",")
	letters := lr[1]
	for _, template := range templates {
		for _, letter := range template {
			if letter != '.' {
				letters += string(letter)
			}
		}
	}
	return d.MultimatchContainsAll(templates, letters)
}

func (lc *LetterCounts) Sub(lc2 LetterCounts) LetterCounts {
	lc3 := make(LetterCounts)
	for k, v := range *lc {
		lc3[k] = v
	}
	for k, v := range lc2 {
		lc3[k] -= v
	}
	return lc3
}

func ContainsAllDictionary(dicts []*Dictionary, count LetterCounts) []*Dictionary {
	if len(dicts) == 1 {
		k := dicts[0].ContainsAll(count)
		x := []*Dictionary{}
		for k2 := range k.Words {
			x = append(x, &Dictionary{Words: map[string]bool{k2: true}})
		}
		return x
	}
	out := []*Dictionary{}

	for _, d := range dicts {
		o := d.ContainsAll(count)
		o.UpdateCounts()
		for k, v := range o.Counts {
			o2 := ContainsAllDictionary(dicts[1:], count.Sub(v))
			if len(o2) > 0 {
				for _, v2 := range o2 {
					d2 := &Dictionary{}
					d2.Words = make(map[string]bool)
					for w := range v2.Words {
						d2.Words[w] = true
					}
					d2.Words[k] = true
					out = append(out, d2)
				}
			}
		}
	}
	return out
}

func (d *Dictionary) MultimatchContainsAll(templates []string, letters string) []*Dictionary {
	count := CountLetters(letters)

	perm := d.ContainsAll(count)
	matches := make([]*Dictionary, len(templates))

	for k, template := range templates {
		matches[k] = perm.Match(template)
	}

	o := ContainsAllDictionary(matches, count)
	return o
}

func DedupeDictionaries(ds []*Dictionary) *Dictionary {
	solutions := make(map[string]bool)
	for _, d := range ds {
		words := []string{}
		for k := range d.Words {
			words = append(words, k)
		}
		sort.Strings(words)
		solutions[strings.Join(words, ",")] = true
	}
	return &Dictionary{Words: solutions}
}
