/**
 * FreeLang v4 StdLib — String Functions (13개)
 *
 * 1.  len(s)                → i32          문자 길이
 * 2.  upper(s)              → str          대문자 변환
 * 3.  lower(s)              → str          소문자 변환
 * 4.  trim(s)               → str          앞뒤 공백 제거
 * 5.  split(s, sep)         → arr[str]     구분자로 분리
 * 6.  join(arr, sep)        → str          배열 → 문자열
 * 7.  contains(s, sub)      → bool         부분 문자열 포함 여부
 * 8.  starts_with(s, pre)   → bool         접두사 확인
 * 9.  ends_with(s, suf)     → bool         접미사 확인
 * 10. replace(s, from, to)  → str          치환 (첫 번째)
 * 11. replace_all(s, f, t)  → str          치환 (전체)
 * 12. slice(s, start, end)  → str          부분 문자열
 * 13. repeat(s, n)          → str          반복
 * 14. index_of(s, sub)      → i32          위치 탐색 (-1 = 없음)
 */

import { Value, V, assertStr, assertI32, assertArr } from "./types";

// ============================================================
// 1. len
// ============================================================

export function str_len(s: Value): Value {
  return V.i32(assertStr(s, "len").length);
}

// ============================================================
// 2. upper
// ============================================================

export function str_upper(s: Value): Value {
  return V.str(assertStr(s, "upper").toUpperCase());
}

// ============================================================
// 3. lower
// ============================================================

export function str_lower(s: Value): Value {
  return V.str(assertStr(s, "lower").toLowerCase());
}

// ============================================================
// 4. trim
// ============================================================

export function str_trim(s: Value): Value {
  return V.str(assertStr(s, "trim").trim());
}

// ============================================================
// 5. split
// ============================================================

export function str_split(s: Value, sep: Value): Value {
  const str = assertStr(s, "split");
  const delim = assertStr(sep, "split");
  const parts = str.split(delim).map((p) => V.str(p));
  return V.arr(parts);
}

// ============================================================
// 6. join
// ============================================================

export function str_join(arr: Value, sep: Value): Value {
  const els  = assertArr(arr, "join");
  const delim = assertStr(sep, "join");
  const strs = els.map((e) => assertStr(e, "join"));
  return V.str(strs.join(delim));
}

// ============================================================
// 7. contains
// ============================================================

export function str_contains(s: Value, sub: Value): Value {
  return V.bool(assertStr(s, "contains").includes(assertStr(sub, "contains")));
}

// ============================================================
// 8. starts_with
// ============================================================

export function str_starts_with(s: Value, prefix: Value): Value {
  return V.bool(assertStr(s, "starts_with").startsWith(assertStr(prefix, "starts_with")));
}

// ============================================================
// 9. ends_with
// ============================================================

export function str_ends_with(s: Value, suffix: Value): Value {
  return V.bool(assertStr(s, "ends_with").endsWith(assertStr(suffix, "ends_with")));
}

// ============================================================
// 10. replace (첫 번째만)
// ============================================================

export function str_replace(s: Value, from: Value, to: Value): Value {
  const str  = assertStr(s, "replace");
  const f    = assertStr(from, "replace");
  const t    = assertStr(to, "replace");
  return V.str(str.replace(f, t));
}

// ============================================================
// 11. replace_all (전체)
// ============================================================

export function str_replace_all(s: Value, from: Value, to: Value): Value {
  const str = assertStr(s, "replace_all");
  const f   = assertStr(from, "replace_all");
  const t   = assertStr(to, "replace_all");
  return V.str(str.replaceAll(f, t));
}

// ============================================================
// 12. slice
// ============================================================

export function str_slice(s: Value, start: Value, end: Value): Value {
  const str = assertStr(s, "slice");
  const st  = assertI32(start, "slice");
  const en  = assertI32(end, "slice");
  // 음수 인덱스 지원: Python 스타일
  const len  = str.length;
  const normSt = st < 0 ? Math.max(0, len + st) : Math.min(st, len);
  const normEn = en < 0 ? Math.max(0, len + en) : Math.min(en, len);
  return V.str(str.slice(normSt, normEn));
}

// ============================================================
// 13. repeat
// ============================================================

export function str_repeat(s: Value, n: Value): Value {
  const str = assertStr(s, "repeat");
  const cnt = assertI32(n, "repeat");
  if (cnt < 0) throw new RangeError("repeat: 음수 반복 횟수");
  return V.str(str.repeat(cnt));
}

// ============================================================
// 14. index_of  (-1 = 없음)
// ============================================================

export function str_index_of(s: Value, sub: Value): Value {
  return V.i32(assertStr(s, "index_of").indexOf(assertStr(sub, "index_of")));
}

// ============================================================
// 추가: pad_start, pad_end, char_at, to_chars
// ============================================================

/** 왼쪽 패딩 */
export function str_pad_start(s: Value, width: Value, fill: Value): Value {
  return V.str(assertStr(s, "pad_start").padStart(assertI32(width, "pad_start"), assertStr(fill, "pad_start")));
}

/** 오른쪽 패딩 */
export function str_pad_end(s: Value, width: Value, fill: Value): Value {
  return V.str(assertStr(s, "pad_end").padEnd(assertI32(width, "pad_end"), assertStr(fill, "pad_end")));
}

/** 특정 위치 문자 */
export function str_char_at(s: Value, idx: Value): Value {
  const str = assertStr(s, "char_at");
  const i   = assertI32(idx, "char_at");
  if (i < 0 || i >= str.length) return V.none();
  return V.some(V.str(str[i]));
}

/** 문자 배열로 분해 */
export function str_to_chars(s: Value): Value {
  return V.arr(assertStr(s, "to_chars").split("").map((c) => V.str(c)));
}
