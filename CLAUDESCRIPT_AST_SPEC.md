# ClaudeScript AST 규격 (Abstract Syntax Tree)

**목적**: Claude가 이 형식으로 코드를 생성하면 자동으로 컴파일 및 실행 가능
**형식**: JSON
**버전**: 0.1.0

---

## 1. 전체 프로그램 구조

```json
{
  "type": "program",
  "version": "0.1.0",
  "definitions": [
    // 함수, 구조체, 타입 정의
  ],
  "instructions": [
    // 실행 명령
  ]
}
```

### 예시: 최소 프로그램

```json
{
  "type": "program",
  "version": "0.1.0",
  "definitions": [],
  "instructions": [
    {
      "type": "call",
      "function": "println",
      "args": [
        {
          "type": "literal",
          "value_type": "string",
          "value": "Hello, ClaudeScript!"
        }
      ]
    }
  ]
}
```

---

## 2. 타입 정의

### 2.1 기본 타입

```json
// 기본 타입
{
  "base": "i32"      // 32비트 정수
  "base": "i64"      // 64비트 정수
  "base": "f64"      // 64비트 실수
  "base": "string"   // 문자열
  "base": "bool"     // 불린
}
```

### 2.2 복합 타입

```json
// 배열
{
  "base": "Array",
  "element_type": {
    "base": "i32"
  }
}

// Optional (null 허용)
{
  "base": "Option",
  "element_type": {
    "base": "string"
  }
}

// 객체 (키-값)
{
  "base": "Object",
  "value_type": {
    "base": "i32"
  }
}

// Map (제너릭)
{
  "base": "Map",
  "key_type": {
    "base": "string"
  },
  "value_type": {
    "base": "i32"
  }
}
```

---

## 3. 변수 및 할당

### 3.1 변수 선언

```json
{
  "type": "var",
  "name": "x",
  "value_type": {
    "base": "i32"
  },
  "value": {
    "type": "literal",
    "value_type": "i32",
    "value": 42
  }
}
```

### 3.2 타입 추론 (타입 생략)

```json
{
  "type": "var",
  "name": "x",
  // value_type 생략 → 컴파일러가 literal 값으로 추론
  "value": {
    "type": "literal",
    "value_type": "i32",
    "value": 42
  }
}
```

### 3.3 변수 재할당

```json
{
  "type": "assign",
  "name": "x",
  "value": {
    "type": "binary_op",
    "op": "+",
    "left": {
      "type": "ref",
      "name": "x"
    },
    "right": {
      "type": "literal",
      "value_type": "i32",
      "value": 1
    }
  }
}
```

---

## 4. 리터럴 및 참조

### 4.1 리터럴

```json
// 정수 리터럴
{
  "type": "literal",
  "value_type": "i32",
  "value": 42
}

// 실수 리터럴
{
  "type": "literal",
  "value_type": "f64",
  "value": 3.14
}

// 문자열 리터럴
{
  "type": "literal",
  "value_type": "string",
  "value": "hello"
}

// 불린 리터럴
{
  "type": "literal",
  "value_type": "bool",
  "value": true
}

// 배열 리터럴
{
  "type": "literal_array",
  "element_type": {
    "base": "i32"
  },
  "values": [
    {
      "type": "literal",
      "value_type": "i32",
      "value": 1
    },
    {
      "type": "literal",
      "value_type": "i32",
      "value": 2
    }
  ]
}

// None (Optional)
{
  "type": "literal",
  "value_type": "none",
  "value": null
}
```

### 4.2 참조

```json
// 변수 참조
{
  "type": "ref",
  "name": "x"
}

// 배열 인덱스 접근
{
  "type": "index",
  "array": {
    "type": "ref",
    "name": "arr"
  },
  "index": {
    "type": "literal",
    "value_type": "i32",
    "value": 0
  }
}

// 객체 키 접근
{
  "type": "field",
  "object": {
    "type": "ref",
    "name": "obj"
  },
  "key": "name"
}
```

---

## 5. 연산

### 5.1 산술 연산

```json
{
  "type": "binary_op",
  "op": "+",  // +, -, *, /, %, ==, !=, <, >, <=, >=, &&, ||
  "left": {
    "type": "literal",
    "value_type": "i32",
    "value": 5
  },
  "right": {
    "type": "literal",
    "value_type": "i32",
    "value": 3
  }
}
```

### 5.2 단항 연산

```json
{
  "type": "unary_op",
  "op": "-",  // -, !
  "operand": {
    "type": "ref",
    "name": "x"
  }
}
```

---

## 6. 함수 정의

### 6.1 기본 함수

```json
{
  "type": "function",
  "name": "add",
  "params": [
    {
      "name": "a",
      "type": {
        "base": "i32"
      }
    },
    {
      "name": "b",
      "type": {
        "base": "i32"
      }
    }
  ],
  "return_type": {
    "base": "i32"
  },
  "body": [
    {
      "type": "return",
      "value": {
        "type": "binary_op",
        "op": "+",
        "left": {
          "type": "ref",
          "name": "a"
        },
        "right": {
          "type": "ref",
          "name": "b"
        }
      }
    }
  ]
}
```

### 6.2 제너릭 함수

```json
{
  "type": "function",
  "name": "first",
  "generics": ["T"],  // 제너릭 타입 변수
  "params": [
    {
      "name": "items",
      "type": {
        "base": "Array",
        "element_type": {
          "base": "T"  // 제너릭 참조
        }
      }
    }
  ],
  "return_type": {
    "base": "Option",
    "element_type": {
      "base": "T"
    }
  },
  "body": [
    {
      "type": "condition",
      "test": {
        "type": "binary_op",
        "op": "==",
        "left": {
          "type": "call",
          "function": "length",
          "args": [
            {
              "type": "ref",
              "name": "items"
            }
          ]
        },
        "right": {
          "type": "literal",
          "value_type": "i32",
          "value": 0
        }
      },
      "then": [
        {
          "type": "return",
          "value": {
            "type": "literal",
            "value_type": "none",
            "value": null
          }
        }
      ],
      "else": [
        {
          "type": "return",
          "value": {
            "type": "some",
            "value": {
              "type": "index",
              "array": {
                "type": "ref",
                "name": "items"
              },
              "index": {
                "type": "literal",
                "value_type": "i32",
                "value": 0
              }
            }
          }
        }
      ]
    }
  ]
}
```

---

## 7. 제어 흐름

### 7.1 조건문 (if/else)

```json
{
  "type": "condition",
  "test": {
    "type": "binary_op",
    "op": ">",
    "left": {
      "type": "ref",
      "name": "x"
    },
    "right": {
      "type": "literal",
      "value_type": "i32",
      "value": 0
    }
  },
  "then": [
    {
      "type": "call",
      "function": "println",
      "args": [
        {
          "type": "literal",
          "value_type": "string",
          "value": "양수"
        }
      ]
    }
  ],
  "else": [
    {
      "type": "call",
      "function": "println",
      "args": [
        {
          "type": "literal",
          "value_type": "string",
          "value": "음수"
        }
      ]
    }
  ]
}
```

### 7.2 패턴 매칭 (match)

```json
{
  "type": "match",
  "value": {
    "type": "ref",
    "name": "maybe_value"
  },
  "cases": [
    {
      "pattern": "Some",
      "bind": "x",  // 바인딩된 변수 이름
      "body": [
        {
          "type": "call",
          "function": "println",
          "args": [
            {
              "type": "ref",
              "name": "x"
            }
          ]
        }
      ]
    },
    {
      "pattern": "None",
      "body": [
        {
          "type": "call",
          "function": "println",
          "args": [
            {
              "type": "literal",
              "value_type": "string",
              "value": "없음"
            }
          ]
        }
      ]
    }
  ]
}
```

### 7.3 반복문 (for)

```json
{
  "type": "for",
  "variable": "i",
  "range": {
    "start": {
      "type": "literal",
      "value_type": "i32",
      "value": 0
    },
    "end": {
      "type": "literal",
      "value_type": "i32",
      "value": 10
    }
  },
  "body": [
    {
      "type": "call",
      "function": "println",
      "args": [
        {
          "type": "ref",
          "name": "i"
        }
      ]
    }
  ]
}
```

### 7.4 while 반복

```json
{
  "type": "while",
  "condition": {
    "type": "binary_op",
    "op": "<",
    "left": {
      "type": "ref",
      "name": "count"
    },
    "right": {
      "type": "literal",
      "value_type": "i32",
      "value": 10
    }
  },
  "body": [
    {
      "type": "assign",
      "name": "count",
      "value": {
        "type": "binary_op",
        "op": "+",
        "left": {
          "type": "ref",
          "name": "count"
        },
        "right": {
          "type": "literal",
          "value_type": "i32",
          "value": 1
        }
      }
    }
  ]
}
```

---

## 8. 함수 호출

### 8.1 일반 호출

```json
{
  "type": "call",
  "function": "add",
  "args": [
    {
      "type": "literal",
      "value_type": "i32",
      "value": 5
    },
    {
      "type": "literal",
      "value_type": "i32",
      "value": 3
    }
  ]
}
```

### 8.2 반환값 저장

```json
{
  "type": "call",
  "function": "add",
  "args": [
    {
      "type": "literal",
      "value_type": "i32",
      "value": 5
    },
    {
      "type": "literal",
      "value_type": "i32",
      "value": 3
    }
  ],
  "assign_to": "result"
}
```

### 8.3 메서드 호출

```json
{
  "type": "method_call",
  "object": {
    "type": "ref",
    "name": "arr"
  },
  "method": "get",
  "args": [
    {
      "type": "literal",
      "value_type": "i32",
      "value": 0
    }
  ],
  "assign_to": "first_item"
}
```

---

## 9. 에러 처리

### 9.1 try/catch

```json
{
  "type": "try",
  "body": [
    {
      "type": "call",
      "function": "risky_operation",
      "args": []
    }
  ],
  "catch": {
    "error_var": "err",
    "error_type": "RuntimeError",
    "body": [
      {
        "type": "call",
        "function": "println",
        "args": [
          {
            "type": "ref",
            "name": "err"
          }
        ]
      }
    ]
  },
  "finally": [
    {
      "type": "call",
      "function": "cleanup",
      "args": []
    }
  ]
}
```

### 9.2 throw

```json
{
  "type": "throw",
  "error_type": "RuntimeError",
  "message": {
    "type": "literal",
    "value_type": "string",
    "value": "Something went wrong"
  }
}
```

---

## 10. 람다/익명 함수

```json
{
  "type": "lambda",
  "params": [
    {
      "name": "x",
      "type": {
        "base": "i32"
      }
    }
  ],
  "return_type": {
    "base": "i32"
  },
  "body": [
    {
      "type": "return",
      "value": {
        "type": "binary_op",
        "op": "*",
        "left": {
          "type": "ref",
          "name": "x"
        },
        "right": {
          "type": "literal",
          "value_type": "i32",
          "value": 2
        }
      }
    }
  ]
}
```

---

## 11. 실제 예시: 팩토리얼

```json
{
  "type": "program",
  "version": "0.1.0",
  "definitions": [
    {
      "type": "function",
      "name": "factorial",
      "params": [
        {
          "name": "n",
          "type": {
            "base": "i32"
          }
        }
      ],
      "return_type": {
        "base": "i32"
      },
      "body": [
        {
          "type": "condition",
          "test": {
            "type": "binary_op",
            "op": "<=",
            "left": {
              "type": "ref",
              "name": "n"
            },
            "right": {
              "type": "literal",
              "value_type": "i32",
              "value": 1
            }
          },
          "then": [
            {
              "type": "return",
              "value": {
                "type": "literal",
                "value_type": "i32",
                "value": 1
              }
            }
          ],
          "else": [
            {
              "type": "return",
              "value": {
                "type": "binary_op",
                "op": "*",
                "left": {
                  "type": "ref",
                  "name": "n"
                },
                "right": {
                  "type": "call",
                  "function": "factorial",
                  "args": [
                    {
                      "type": "binary_op",
                      "op": "-",
                      "left": {
                        "type": "ref",
                        "name": "n"
                      },
                      "right": {
                        "type": "literal",
                        "value_type": "i32",
                        "value": 1
                      }
                    }
                  ]
                }
              }
            }
          ]
        }
      ]
    }
  ],
  "instructions": [
    {
      "type": "call",
      "function": "factorial",
      "args": [
        {
          "type": "literal",
          "value_type": "i32",
          "value": 5
        }
      ],
      "assign_to": "result"
    },
    {
      "type": "call",
      "function": "println",
      "args": [
        {
          "type": "binary_op",
          "op": "+",
          "left": {
            "type": "literal",
            "value_type": "string",
            "value": "5! = "
          },
          "right": {
            "type": "call",
            "function": "to_string",
            "args": [
              {
                "type": "ref",
                "name": "result"
              }
            ]
          }
        }
      ]
    }
  ]
}
```

**예상 출력**:
```
5! = 120
```

---

## 12. 검증 규칙 (컴파일러가 검사할 사항)

### 12.1 타입 검증
- [ ] 변수 재할당 시 타입 일치 확인
- [ ] 함수 호출 인자 타입 일치
- [ ] 이항 연산 양쪽 타입 호환성
- [ ] 배열 인덱스는 i32만 허용
- [ ] 암묵적 타입 변환 금지

### 12.2 범위 검증
- [ ] 함수 정의 후 호출 (앞서 정의 원칙)
- [ ] 루프 변수 범위 확인
- [ ] 함수 인자 개수 확인
- [ ] 재귀 깊이 제한 (스택 오버플로우 방지)

### 12.3 안전성 검증
- [ ] null 참조 금지 (Option 강제)
- [ ] 배열 인덱스 범위 확인 (또는 Option 반환)
- [ ] 함수 반환값 타입 확인
- [ ] break/continue는 루프 내에서만

---

## 13. 확장 가능성

이 AST는 향후 다음을 추가할 수 있도록 설계됨:

```json
// 구조체 정의 (미래)
{
  "type": "struct",
  "name": "Point",
  "fields": [
    {"name": "x", "type": {"base": "f64"}},
    {"name": "y", "type": {"base": "f64"}}
  ]
}

// 특성/인터페이스 (미래)
{
  "type": "trait",
  "name": "Drawable",
  "methods": [...]
}

// async 함수 (미래)
{
  "type": "async_function",
  "name": "fetch_data",
  "body": [...]
}
```

---

**작성**: 2026-03-07
**최종 검토**: AST 스키마 완성

