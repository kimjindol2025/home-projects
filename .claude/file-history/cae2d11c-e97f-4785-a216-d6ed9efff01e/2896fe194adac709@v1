/**
 * JIT Compiler & Executor 테스트
 */

import { Lexer } from './lexer';
import { Parser } from './parser';
import { BytecodeGenerator } from './bytecode_generator';
import { JITExecutor } from './jit_executor';

// ============================================================
// 헬퍼 함수
// ============================================================

function compileAndExecute(source: string): any {
  // 1. Lexing
  const lexer = new Lexer(source);
  const tokens = lexer.tokenize();

  // 2. Parsing
  const parser = new Parser(tokens);
  const ast = parser.parse();

  // 3. Bytecode Generation
  const generator = new BytecodeGenerator();
  const chunk = generator.compile(ast);

  // 4. JIT Execution
  const executor = new JITExecutor(chunk);
  const result = executor.execute();

  return result;
}

// ============================================================
// 테스트 케이스
// ============================================================

describe('JIT Compiler & Executor', () => {
  test('PUSH & CONST: 숫자 상수', () => {
    const source = `
      x = 42
      print(x)
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('ADD: 덧셈 연산', () => {
    const source = `
      x = 5
      y = 3
      z = x + y
      print(z)
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('SUB: 뺄셈 연산', () => {
    const source = `
      a = 10
      b = 4
      c = a - b
      print(c)
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('MUL: 곱셈 연산', () => {
    const source = `
      x = 6
      y = 7
      z = x * y
      print(z)
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('DIV: 나눗셈 연산', () => {
    const source = `
      a = 20
      b = 4
      c = a / b
      print(c)
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('MOD: 나머지 연산', () => {
    const source = `
      x = 17
      y = 5
      z = x % y
      print(z)
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('EQ: 같음 비교', () => {
    const source = `
      a = 5
      b = 5
      c = a == b
      print(c)
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('NE: 다름 비교', () => {
    const source = `
      a = 5
      b = 3
      c = a != b
      print(c)
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('LT: 미만 비교', () => {
    const source = `
      a = 3
      b = 5
      c = a < b
      print(c)
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('LE: 이하 비교', () => {
    const source = `
      a = 5
      b = 5
      c = a <= b
      print(c)
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('GT: 초과 비교', () => {
    const source = `
      a = 7
      b = 3
      c = a > b
      print(c)
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('GE: 이상 비교', () => {
    const source = `
      a = 5
      b = 5
      c = a >= b
      print(c)
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('AND: 논리 AND', () => {
    const source = `
      a = 1
      b = 1
      c = a && b
      print(c)
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('OR: 논리 OR', () => {
    const source = `
      a = 0
      b = 1
      c = a || b
      print(c)
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('NOT: 논리 NOT', () => {
    const source = `
      a = 1
      b = !a
      print(b)
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('if: 조건문', () => {
    const source = `
      x = 10
      if (x > 5) {
        print("big")
      } else {
        print("small")
      }
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('while: 루프', () => {
    const source = `
      i = 0
      while (i < 3) {
        print(i)
        i = i + 1
      }
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('ARRAY_CREATE: 배열 생성', () => {
    const source = `
      arr = [1, 2, 3]
      print(arr)
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('ARRAY_GET: 배열 인덱싱', () => {
    const source = `
      arr = [10, 20, 30]
      x = arr[1]
      print(x)
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('ARRAY_LEN: 배열 길이', () => {
    const source = `
      arr = [1, 2, 3, 4, 5]
      n = len(arr)
      print(n)
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('CALL: 함수 호출', () => {
    const source = `
      fn add(a, b) {
        return a + b
      }

      result = add(5, 3)
      print(result)
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('RETURN: 함수 반환', () => {
    const source = `
      fn square(x) {
        return x * x
      }

      y = square(7)
      print(y)
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('재귀 함수: 팩토리얼', () => {
    const source = `
      fn factorial(n) {
        if (n <= 1) {
          return 1
        } else {
          return n * factorial(n - 1)
        }
      }

      result = factorial(5)
      print(result)
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('재귀 함수: 피보나치', () => {
    const source = `
      fn fib(n) {
        if (n <= 1) {
          return n
        } else {
          return fib(n - 1) + fib(n - 2)
        }
      }

      result = fib(10)
      print(result)
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('BUILTIN: type 함수', () => {
    const source = `
      x = 42
      t = type(x)
      print(t)
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('BUILTIN: str 함수', () => {
    const source = `
      x = 42
      s = str(x)
      print(s)
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('BUILTIN: num 함수', () => {
    const source = `
      s = "123"
      n = num(s)
      print(n)
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('복합: 계산기', () => {
    const source = `
      fn add(a, b) { return a + b }
      fn sub(a, b) { return a - b }
      fn mul(a, b) { return a * b }

      print(add(10, 5))
      print(sub(10, 5))
      print(mul(10, 5))
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('복합: 배열 합계', () => {
    const source = `
      fn sum_array(arr) {
        total = 0
        i = 0
        while (i < len(arr)) {
          total = total + arr[i]
          i = i + 1
        }
        return total
      }

      numbers = [1, 2, 3, 4, 5]
      result = sum_array(numbers)
      print(result)
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('복합: 최댓값 찾기', () => {
    const source = `
      fn max_array(arr) {
        if (len(arr) == 0) {
          return null
        }
        max = arr[0]
        i = 1
        while (i < len(arr)) {
          if (arr[i] > max) {
            max = arr[i]
          }
          i = i + 1
        }
        return max
      }

      numbers = [5, 2, 8, 1, 9, 3]
      result = max_array(numbers)
      print(result)
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('복합: 배열 필터링', () => {
    const source = `
      fn count_even(arr) {
        count = 0
        i = 0
        while (i < len(arr)) {
          if (arr[i] % 2 == 0) {
            count = count + 1
          }
          i = i + 1
        }
        return count
      }

      numbers = [1, 2, 3, 4, 5, 6, 7, 8]
      result = count_even(numbers)
      print(result)
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });

  test('복합: 중첩 루프', () => {
    const source = `
      i = 0
      while (i < 3) {
        j = 0
        while (j < 3) {
          print(i)
          print(j)
          j = j + 1
        }
        i = i + 1
      }
    `;
    expect(() => compileAndExecute(source)).not.toThrow();
  });
});

// ============================================================
// 성능 벤치마크
// ============================================================

describe('JIT Performance', () => {
  test('성능: 1000번 루프', () => {
    const source = `
      i = 0
      while (i < 1000) {
        i = i + 1
      }
      print(i)
    `;

    const start = performance.now();
    compileAndExecute(source);
    const elapsed = performance.now() - start;

    console.log(`[Performance] 1000-loop: ${elapsed.toFixed(2)}ms`);
    expect(elapsed).toBeLessThan(100); // 100ms 이하
  });

  test('성능: 피보나치(15)', () => {
    const source = `
      fn fib(n) {
        if (n <= 1) {
          return n
        } else {
          return fib(n - 1) + fib(n - 2)
        }
      }
      print(fib(15))
    `;

    const start = performance.now();
    compileAndExecute(source);
    const elapsed = performance.now() - start;

    console.log(`[Performance] Fibonacci(15): ${elapsed.toFixed(2)}ms`);
    expect(elapsed).toBeLessThan(500); // 500ms 이하
  });

  test('성능: 배열 처리 (1000 원소)', () => {
    const source = `
      fn sum_array(arr) {
        total = 0
        i = 0
        while (i < len(arr)) {
          total = total + arr[i]
          i = i + 1
        }
        return total
      }

      arr = [1, 2, 3, 4, 5]
      i = 0
      while (i < 200) {
        arr = [1, 2, 3, 4, 5]
        sum_array(arr)
        i = i + 1
      }
      print("done")
    `;

    const start = performance.now();
    compileAndExecute(source);
    const elapsed = performance.now() - start;

    console.log(`[Performance] Array processing: ${elapsed.toFixed(2)}ms`);
    expect(elapsed).toBeLessThan(500);
  });
});
