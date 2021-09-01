package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	reversePhrase := "Hello, OTUS!"
	fmt.Println(stringutil.Reverse(reversePhrase))
}
