package main

import (
	"fmt"
	"github.com/kimjindol2025/freelang-go/pkg/lexer"
)

func main() {
	input := "let x = 5;"
	l := lexer.New(input)
	fmt.Println("=== Token Stream ===")
	for {
		tok := l.NextToken()
		fmt.Printf("%v\n", tok)
		if tok.Type == "EOF" {
			break
		}
	}
}
