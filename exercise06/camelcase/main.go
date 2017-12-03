package main

import (
	"fmt"
	"unicode"
)

func camelcase(input string) int {
	if input == "" {
		return 0
	}
	wc := 1
	for _, ch := range input {
		if unicode.IsUpper(ch) {
			wc++
		}
	}
	return wc
}

func main() {
	var input string
	fmt.Scanf("%s\n", &input)
	fmt.Println(camelcase(input))
}
