package main

import "fmt"

func caesarcipher1(input string, key int) string {
	var output []rune
	for _, ch := range input {
		output = append(output, encodeAlpha(ch, key))
	}
	return string(output)
}

func encodeAlpha(ch rune, key int) rune {
	switch {
	case ch >= 'a' && ch <= 'z':
		return rotate(ch, 'a', key)
	case ch >= 'A' && ch <= 'Z':
		return rotate(ch, 'A', key)
	default:
		return ch
	}
}

func rotate(ch rune, base int, key int) rune {
	rebased := int(ch) - base
	rotated := (rebased + key) % 26
	return rune(rotated + base)
}

func main() {
	var input string
	var length, key int
	fmt.Scanf("%d\n%s\n%d\n", &length, &input, &key)
	fmt.Println(caesarcipher1(input, key))
}
