/**
 * FreeLang v4 StdLib — Regex Functions (3개 + 확장)
 *
 * 1. regex_match(s, pattern, flags?)  → bool | arr[str]
 *    - flags 없으면 bool (패턴 일치 여부)
 *    - flags에 "g" 포함 시 전체 매치 배열 반환
 *
 * 2. regex_split(s, pattern, flags?)  → arr[str]
 *    - 정규식 구분자로 분리
 *
 * 3. regex_replace(s, pattern, repl, flags?) → str
 *    - 정규식 패턴 치환
 *    - flags에 "g" → 전체 치환, 없으면 첫 번째만
 *
 * + regex_test(s, pattern)    → bool
 * + regex_groups(s, pattern)  → arr[str] (캡처 그룹)
 * + regex_valid(pattern)      → bool
 */

import { Value, V, assertStr } from "./types";

// ============================================================
// 내부 헬퍼: 안전한 RegExp 생성
// ============================================================

function makeRegex(pattern: string, flags: string): RegExp {
  try {
    return new RegExp(pattern, flags);
  } catch (e: any) {
    throw new Error(`regex: 유효하지 않은 패턴 '${pattern}': ${e.message}`);
  }
}

// ============================================================
// 1. regex_match
// ============================================================

/**
 * regex_match(s, pattern)         → bool   (단순 일치 여부)
 * regex_match(s, pattern, "g")    → arr[str] (전체 매치 목록)
 * regex_match(s, pattern, "gi")   → arr[str] (대소문자 무시 + 전체)
 */
export function regex_match(s: Value, pattern: Value, flags: Value = V.str("")): Value {
  const str  = assertStr(s, "regex_match");
  const pat  = assertStr(pattern, "regex_match");
  const flgs = assertStr(flags, "regex_match");

  if (flgs.includes("g")) {
    // 전체 매치 모드
    const re      = makeRegex(pat, flgs);
    const matches = [...str.matchAll(re)].map((m) => V.str(m[0]));
    return V.arr(matches);
  } else {
    // 단순 일치 여부
    const re = makeRegex(pat, flgs);
    return V.bool(re.test(str));
  }
}

// ============================================================
// 2. regex_split
// ============================================================

export function regex_split(s: Value, pattern: Value, flags: Value = V.str("")): Value {
  const str  = assertStr(s, "regex_split");
  const pat  = assertStr(pattern, "regex_split");
  const flgs = assertStr(flags, "regex_split");

  const re = makeRegex(pat, flgs);
  return V.arr(str.split(re).map((p) => V.str(p)));
}

// ============================================================
// 3. regex_replace
// ============================================================

/**
 * regex_replace(s, pattern, repl)      → str   (첫 번째 매치 치환)
 * regex_replace(s, pattern, repl, "g") → str   (전체 치환)
 */
export function regex_replace(
  s:       Value,
  pattern: Value,
  repl:    Value,
  flags:   Value = V.str("")
): Value {
  const str  = assertStr(s, "regex_replace");
  const pat  = assertStr(pattern, "regex_replace");
  const r    = assertStr(repl, "regex_replace");
  const flgs = assertStr(flags, "regex_replace");

  const re = makeRegex(pat, flgs);
  return V.str(flgs.includes("g") ? str.replaceAll(re, r) : str.replace(re, r));
}

// ============================================================
// 4. regex_test — 단순 bool 테스트
// ============================================================

export function regex_test(s: Value, pattern: Value, flags: Value = V.str("")): Value {
  const str  = assertStr(s, "regex_test");
  const pat  = assertStr(pattern, "regex_test");
  const flgs = assertStr(flags, "regex_test");
  return V.bool(makeRegex(pat, flgs).test(str));
}

// ============================================================
// 5. regex_groups — 캡처 그룹 반환
// ============================================================

export function regex_groups(s: Value, pattern: Value, flags: Value = V.str("")): Value {
  const str  = assertStr(s, "regex_groups");
  const pat  = assertStr(pattern, "regex_groups");
  const flgs = assertStr(flags, "regex_groups");

  const re = makeRegex(pat, flgs);
  const m  = str.match(re);
  if (!m) return V.none();

  // m[0] = 전체 매치, m[1..] = 캡처 그룹
  const groups = m.slice(1).map((g) => (g !== undefined ? V.str(g) : V.none()));
  return V.some(V.arr(groups));
}

// ============================================================
// 6. regex_valid — 패턴 유효성 확인
// ============================================================

export function regex_valid(pattern: Value): Value {
  const pat = assertStr(pattern, "regex_valid");
  try {
    new RegExp(pat);
    return V.bool(true);
  } catch {
    return V.bool(false);
  }
}

// ============================================================
// 7. regex_count — 패턴 매치 횟수
// ============================================================

export function regex_count(s: Value, pattern: Value, flags: Value = V.str("g")): Value {
  const str  = assertStr(s, "regex_count");
  const pat  = assertStr(pattern, "regex_count");
  const flgs = assertStr(flags, "regex_count");

  // 항상 g 플래그 포함
  const fStr = flgs.includes("g") ? flgs : flgs + "g";
  const re   = makeRegex(pat, fStr);
  return V.i32([...str.matchAll(re)].length);
}
