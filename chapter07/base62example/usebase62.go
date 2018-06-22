package main

import (
	"fmt"

	"github.com/dakaraj/go-api/chapter07/urlshortener/utils"
)

func main() {
	x := 100
	base62String := base62.ToBase62(x)
	fmt.Println(base62String)
	normalNumber := base62.FromBase62(base62String)
	fmt.Println(normalNumber)
}
