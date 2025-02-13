package main

import "strings"
import "slices"

func removeProfanity(raw string) string {
	l := strings.Split(raw, " ")

	naughtyList := []string{"kerfuffle", "sharbert", "fornax"}
	for i, word := range l {
		if slices.Contains(naughtyList, strings.ToLower(word)) {
			l[i] = "****"
		}
	}
	return strings.Join(l, " ")
}
