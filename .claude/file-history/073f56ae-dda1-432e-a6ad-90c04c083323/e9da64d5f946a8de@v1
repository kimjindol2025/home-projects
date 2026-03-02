/**
 * FreeLang v4 StdLib — Functional Programming (4개 + 확장)
 *
 * 1. compose(...fns)         → (x) → y      우→좌 함수 합성
 * 2. curry(fn, arity)        → curried fn   커링 (인자 개수 지정)
 * 3. memoize(fn)             → fn           메모이제이션 (순수 함수 캐싱)
 * 4. partial(fn, ...args)    → fn           부분 적용
 *
 * + pipe(...fns)             좌→우 합성
 * + once(fn)                 최초 1회만 실행
 * + throttle(fn, ms)         스로틀링
 * + flip(fn)                 인수 순서 뒤집기
 */

import { Value } from "./types";

// FreeLang Value를 받아 Value를 반환하는 함수 타입
export type FlFn = (...args: Value[]) => Value;

// ============================================================
// 1. compose — 우→좌 합성  compose(f, g)(x) = f(g(x))
// ============================================================

export function compose(...fns: FlFn[]): FlFn {
  if (fns.length === 0) throw new Error("compose: 함수가 1개 이상 필요합니다");
  return (...args: Value[]): Value => {
    const reversed = [...fns].reverse();
    let result: Value = reversed[0](...args);
    for (let i = 1; i < reversed.length; i++) {
      result = reversed[i](result);
    }
    return result;
  };
}

// ============================================================
// 2. pipe — 좌→우 합성  pipe(f, g)(x) = g(f(x))
// ============================================================

export function pipe(...fns: FlFn[]): FlFn {
  if (fns.length === 0) throw new Error("pipe: 함수가 1개 이상 필요합니다");
  return (...args: Value[]): Value => {
    let result: Value = fns[0](...args);
    for (let i = 1; i < fns.length; i++) {
      result = fns[i](result);
    }
    return result;
  };
}

// ============================================================
// 3. curry — 커링 (arity개 인수를 하나씩 받음)
// ============================================================

export function curry(fn: FlFn, arity: number): FlFn {
  if (arity <= 0) throw new RangeError("curry: arity는 1 이상이어야 합니다");

  function curried(collected: Value[]): FlFn {
    return (...args: Value[]): Value => {
      const next = [...collected, ...args];
      if (next.length >= arity) {
        return fn(...next.slice(0, arity));
      }
      return curried(next) as unknown as Value; // 타입 편의상 캐스팅
    };
  }

  return curried([]);
}

// ============================================================
// 4. memoize — 순수 함수 결과 캐싱 (인수는 JSON 직렬화)
// ============================================================

export function memoize(fn: FlFn): FlFn {
  const cache = new Map<string, Value>();
  return (...args: Value[]): Value => {
    const key = JSON.stringify(args, (_k, v) => {
      // Map은 배열로 직렬화
      if (v instanceof Map) return Object.fromEntries(v);
      return v;
    });
    if (cache.has(key)) return cache.get(key)!;
    const result = fn(...args);
    cache.set(key, result);
    return result;
  };
}

/** 캐시 히트 수를 추적하는 memoize (테스트용) */
export function memoize_tracked(fn: FlFn): { fn: FlFn; hits: () => number; misses: () => number } {
  const cache = new Map<string, Value>();
  let hits = 0, misses = 0;
  const wrapped: FlFn = (...args: Value[]): Value => {
    const key = JSON.stringify(args, (_k, v) => (v instanceof Map ? Object.fromEntries(v) : v));
    if (cache.has(key)) { hits++; return cache.get(key)!; }
    misses++;
    const result = fn(...args);
    cache.set(key, result);
    return result;
  };
  return { fn: wrapped, hits: () => hits, misses: () => misses };
}

// ============================================================
// 5. partial — 부분 적용  partial(f, a, b)(c) = f(a, b, c)
// ============================================================

export function partial(fn: FlFn, ...partialArgs: Value[]): FlFn {
  return (...remainingArgs: Value[]): Value => fn(...partialArgs, ...remainingArgs);
}

// ============================================================
// 6. once — 최초 1회만 실행, 이후 첫 번째 결과 반환
// ============================================================

export function once(fn: FlFn): FlFn {
  let called = false;
  let result: Value | null = null;
  return (...args: Value[]): Value => {
    if (!called) {
      called = true;
      result = fn(...args);
    }
    return result!;
  };
}

// ============================================================
// 7. flip — 인수 첫 두 개 순서 뒤집기
// ============================================================

export function flip(fn: FlFn): FlFn {
  return (...args: Value[]): Value => {
    if (args.length < 2) return fn(...args);
    const [a, b, ...rest] = args;
    return fn(b, a, ...rest);
  };
}

// ============================================================
// 8. identity — 항등 함수
// ============================================================

export function identity(x: Value): Value {
  return x;
}

// ============================================================
// 9. constant — 상수 함수
// ============================================================

export function constant(x: Value): FlFn {
  return () => x;
}
