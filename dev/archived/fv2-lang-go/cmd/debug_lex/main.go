package main

import (
    "fmt"
    "fv2-lang/internal/lexer"
)

func main() {
    l, _ := lexer.New(`let s = ""`)
    tokens, _ := l.Tokenize()
    for _, tok := range tokens {
        fmt.Printf("Type=%v Text=%q\n", tok.Type, tok.Text)
    }
}
