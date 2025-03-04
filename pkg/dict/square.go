package dict

import (
	"log"
	"strings"
)

type Coord struct {
	X int
	Y int
}

var Dir8 []Coord = []Coord{
	{-1, -1},
	{0, -1},
	{1, -1},
	{-1, 0},
	// {0, 0},
	{1, 0},
	{-1, 1},
	{0, 1},
	{1, 1},
}

func (d *Dictionary) Square(square string) *Dictionary {
	d.UpdateCounts()

	if len(square) != 16 {
		// fmt.Println("Wrong count")
		log.Fatal("x")
		return d
	}
	square = strings.ToLower(square)
	out := make(map[string]bool)
	square_map := make(map[Coord]rune)
	square_rune := make(map[rune][]Coord)

	i, j := 0, 0
	for _, v := range square {
		i++
		if i >= 4 {
			j++
			i = 0
		}
		square_map[Coord{i, j}] = v
		_, ok := square_rune[v]
		if !ok {
			square_rune[v] = make([]Coord, 0)
		}
		square_rune[v] = append(square_rune[v], Coord{i, j})
	}
	// fmt.Println(square_rune)
	// fmt.Println(square_map)

a:
	for word := range d.Words {
		if len(word) > 16 {
			continue a
		}

		// Check if all the letters of this word is in the square
		for k, v := range d.Counts[word] {
			coords, ok := square_rune[rune(k)]
			if !ok {
				continue a
			}
			if len(coords) < int(v) {
				continue a
			}
		}

		firstletter := rune(word[0])

		for _, coord := range square_rune[firstletter] {
			done := make(map[Coord]bool)
			done[coord] = true
			path := make([]Coord, 0, 16)
			path = append(path, coord)

			for i := 1; i < len(word); i++ {
				coord2 := coord
				found := false
				for _, dir := range Dir8 {
					coord2.X = coord.X + dir.X
					coord2.Y = coord.Y + dir.Y
					if done[coord2] {
						continue
					}
					if square_map[coord2] != rune(word[i]) {
						continue
					}
					found = true
				}
				if found {
					path = append(path, coord2)
					coord = coord2
					done[coord] = true
					_ = path
				} else { // nolint

				}
			}

		}

		// _ = done
		// _ = coords

		// for _, coord := range coords {
		// 	done[coord] = true
		// }

		out[word] = true
	}

	return &Dictionary{Words: out}
}
