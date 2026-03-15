package main

import (
	"fmt"
	"os"

	"github.com/freelang-ai/freelang-go/internal/lexer"
)

func main() {
	// 테스트 코드
	code := `
// 변수 선언
let x = 10
let y = 20.5
let name = "hello"
let isValid = true

// 함수 선언
fn add(a: number, b: number) -> number {
    return a + b
}

// 조건문
if x > 5 {
    y = y + 1
} else {
    y = y - 1
}

// 루프
for i in range(0, 10) {
    x = x + i
}

// 문자열 포함
let msg = "hello \"world\""

// 주석
// 라인 주석

/* 블록
   주석 */

// 정규식
let pattern = /[0-9]+/g
`

	// 렉서 생성
	lex := lexer.NewLexer(code)

	// 모든 토큰 출력
	fmt.Println("=== FreeLang Lexer (Go) ===\n")
	tokens := lex.Tokenize()

	for i, token := range tokens {
		fmt.Printf("[%3d] %s\n", i, token)
	}

	fmt.Printf("\nTotal tokens: %d\n", len(tokens))
}
