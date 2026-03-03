/**
 * Ownership/Borrow Checker 테스트
 *
 * 35개 테스트 (Week 6에서 구현)
 * - 소유권 규칙 검증
 * - 차용 규칙 검증
 * - 생명주기 검증
 * - 에러 감지 테스트
 */

import { Lexer } from './lexer';
import { Parser } from './parser';
import { OwnershipChecker, OwnershipError } from './ownership_checker';
import { BorrowChecker, BorrowError } from './borrow_checker';

// ============================================================
// Helper Functions
// ============================================================

function checkCode(code: string): {
  ownershipErrors: OwnershipError[];
  borrowErrors: BorrowError[];
} {
  try {
    const lexer = new Lexer(code);
    const tokens = lexer.tokenize();
    const parser = new Parser(tokens);
    const ast = parser.parse();

    const ownershipChecker = new OwnershipChecker();
    const borrowChecker = new BorrowChecker();

    const ownershipErrors = ownershipChecker.check(ast);
    const borrowErrors = borrowChecker.check(ast);

    return { ownershipErrors, borrowErrors };
  } catch (e) {
    throw new Error(`Parse error: ${e}`);
  }
}

// ============================================================
// Tests (Skeleton)
// ============================================================

describe('Ownership Rules', () => {
  describe('Rule 1: Single Owner', () => {
    test('simple variable ownership', () => {
      const code = `
        x = 42
      `;
      const { ownershipErrors } = checkCode(code);
      expect(ownershipErrors).toHaveLength(0);
    });

    test('move semantics: move ownership', () => {
      const code = `
        x = [1, 2, 3]
        y = x
      `;
      // Week 6: 구현 필요
      // x는 더 이상 소유하지 않음
    });

    test('copy semantics: number copy', () => {
      const code = `
        a = 42
        b = a
      `;
      // Week 6: 구현 필요
      // a와 b는 모두 유효
    });

    test('use after move error', () => {
      const code = `
        x = "hello"
        y = x
        print(x)
      `;
      // Week 6: 구현 필요
      // ownershipErrors.length > 0 (use_after_move)
    });
  });

  describe('Rule 2: Scope-Based Cleanup', () => {
    test('variable cleanup on scope end', () => {
      const code = `
        {
          x = 100
        }
      `;
      // Week 6: 구현 필요
      // x는 스코프 끝에서 자동 해제
    });

    test('nested scopes', () => {
      const code = `
        x = 1
        {
          y = 2
          {
            z = 3
          }
        }
      `;
      // Week 6: 구현 필요
      // 각 변수의 스코프 범위 검증
    });

    test('use after scope end', () => {
      const code = `
        {
          x = 42
        }
        print(x)
      `;
      // Week 6: 구현 필요
      // ownershipErrors.length > 0 (use_after_free)
    });
  });

  describe('Rule 3: Function Ownership', () => {
    test('function parameter ownership transfer', () => {
      const code = `
        fn take(x: int) {
          print(x)
        }
        a = 42
        take(a)
      `;
      // Week 6: 구현 필요
      // a는 함수로 이동
    });

    test('function return ownership', () => {
      const code = `
        fn get_value() -> int {
          x = 100
          return x
        }
        y = get_value()
      `;
      // Week 6: 구현 필요
      // y가 반환값의 소유권을 받음
    });

    test('reference parameter no transfer', () => {
      const code = `
        fn borrow(x: &int) {
          print(x)
        }
        a = 42
        borrow(&a)
        print(a)
      `;
      // Week 6: 구현 필요
      // a는 여전히 유효
    });
  });
});

describe('Borrow Rules', () => {
  describe('Rule 1: Shared Borrow', () => {
    test('multiple shared borrows allowed', () => {
      const code = `
        x = 100
        r1 = &x
        r2 = &x
        r3 = &x
      `;
      // Week 6: 구현 필요
      // borrowErrors.length === 0
    });

    test('shared borrow read access only', () => {
      const code = `
        x = 100
        r = &x
        print(r)
      `;
      // Week 6: 구현 필요
      // r로 읽기 가능
    });

    test('cannot modify through shared borrow', () => {
      const code = `
        x = 100
        r = &x
        r = 200
      `;
      // Week 6: 구현 필요
      // borrowErrors.length > 0 (modification through shared borrow)
    });
  });

  describe('Rule 2: Mutable Borrow', () => {
    test('single mutable borrow allowed', () => {
      const code = `
        x = 100
        m = &mut x
        m = 200
      `;
      // Week 6: 구현 필요
      // borrowErrors.length === 0
    });

    test('multiple mutable borrows not allowed', () => {
      const code = `
        x = [1, 2, 3]
        m1 = &mut x
        m2 = &mut x
      `;
      // Week 6: 구현 필요
      // borrowErrors.length > 0 (multiple_mutable)
    });

    test('exclusive mutable access', () => {
      const code = `
        x = 100
        m = &mut x
        print(x)
      `;
      // Week 6: 구현 필요
      // borrowErrors.length > 0 (use_during_borrow)
    });
  });

  describe('Rule 3: Shared + Mutable Conflict', () => {
    test('shared and mutable cannot coexist', () => {
      const code = `
        x = 100
        s = &x
        m = &mut x
      `;
      // Week 6: 구현 필요
      // borrowErrors.length > 0 (shared_mutable_conflict)
    });

    test('mutable then shared conflict', () => {
      const code = `
        x = 100
        m = &mut x
        s = &x
      `;
      // Week 6: 구현 필요
      // borrowErrors.length > 0 (shared_mutable_conflict)
    });
  });
});

describe('Lifetime Validation', () => {
  test('reference lifetime validation', () => {
    const code = `
      fn first(arr: &[int]) -> &int {
        return &arr[0]
      }
      a = [1, 2, 3]
      first(&a)
    `;
    // Week 6: 구현 필요
  });

  test('dangling reference prevention', () => {
    const code = `
      fn bad_ref() -> &int {
        x = 10
        return &x
      }
    `;
    // Week 6: 구현 필요
    // lifetime mismatch 에러
  });
});

describe('Error Cases (감지 테스트)', () => {
  test('detect: use after move', () => {
    const code = `
      x = "hello"
      y = x
      print(x)
    `;
    // Week 6: 구현 필요
    // expect(ownershipErrors).toContainEqual(
    //   expect.objectContaining({ type: 'use_after_move' })
    // );
  });

  test('detect: multiple mutable borrows', () => {
    const code = `
      x = [1, 2, 3]
      m1 = &mut x
      m2 = &mut x
    `;
    // Week 6: 구현 필요
    // expect(borrowErrors).toContainEqual(
    //   expect.objectContaining({ type: 'multiple_mutable' })
    // );
  });

  test('detect: shared and mutable conflict', () => {
    const code = `
      x = 100
      s = &x
      m = &mut x
    `;
    // Week 6: 구현 필요
    // expect(borrowErrors).toContainEqual(
    //   expect.objectContaining({ type: 'shared_mutable_conflict' })
    // );
  });

  test('detect: use after free', () => {
    const code = `
      {
        x = 42
      }
      print(x)
    `;
    // Week 6: 구현 필요
    // expect(ownershipErrors).toContainEqual(
    //   expect.objectContaining({ type: 'use_after_free' })
    // );
  });

  // ... 11개 더 (Week 6에 구현)
});

describe('Complex Scenarios', () => {
  test('scenario 1: function with reference parameters', () => {
    const code = `
      fn modify(arr: &mut [int]) {
        arr[0] = 42
      }
      a = [1, 2, 3]
      modify(&mut a)
      print(a)
    `;
    // Week 6: 구현 필요
  });

  test('scenario 2: return borrowed reference', () => {
    const code = `
      fn get_element(arr: &[int], i: int) -> &int {
        return &arr[i]
      }
      a = [10, 20, 30]
      e = get_element(&a, 0)
      print(e)
    `;
    // Week 6: 구현 필요
  });

  test('scenario 3: nested borrows', () => {
    const code = `
      fn process(x: &int) {
        y = &x
        print(y)
      }
      a = 100
      process(&a)
    `;
    // Week 6: 구현 필요
  });
});
