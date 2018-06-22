package base62

import (
	"strings"
)

const (
	base = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b    = 62
)

// ToBase62 converts int value into base62 encoded string
func ToBase62(num int) string {
	r := num % b
	res := string(base[r])
	div := num / b

	for div != 0 {
		r = div % b
		div = div / b
		res = string(base[r]) + res
	}

	return res
}

// FromBase62 converts base62 encoded string to int
func FromBase62(str string) int {
	res := 0
	for _, r := range str {
		res = (b * res) + strings.Index(base, string(r))
	}

	return res
}
