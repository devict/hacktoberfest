package main

import (
	"sort"
)

// top gives the keys of the top n values in a map[string]int. Tied values will
// be included and counted separately but will not be sorted consistently.
func top(n int, data map[string]int) []string {

	ranked := make(map[int][]string)
	for k, v := range data {
		ranked[v] = append(ranked[v], k)
	}

	var keys []int
	for k := range ranked {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(keys)))

	var langs []string
	for _, rank := range keys {
		for _, l := range ranked[rank] {
			langs = append(langs, l)
		}
	}

	if len(langs) < n {
		n = len(langs)
	}

	return langs[:n]
}
