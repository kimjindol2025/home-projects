# 진정한 FreeLang 자체호스팅 파이프라인 구현 계획

## 문제 / 목표
MEMORY.md 검증: "자체호스팅 완성"이 허위였음 확인.
진짜 목표: `hello.free` → ELF 바이너리 → `./hello` → "Hello World!" 실행까지.

## 현재 자산 (검증됨)

| 자산 | 경로 | 상태 |
|------|------|------|
| C 기반 FreeLang 인터프리터 | `/tmp/freelang-c-final/` | ✅ 빌드 가능 (make) |
| bootstrap-demo.fl (렉서) | `/data/data/com.termux/files/home/v2-freelang-ai/self-hosting/` | ✅ 존재 |
| lexer.fl, parser.fl | 동 경로 | ✅ 존재 |

## 핵심 발견 (문제)

1. `fl_stdlib_register()` 함수 - **주석만 있고 실제 등록 코드 없음**
2. **파일 I/O 함수 없음** (write_file, read_file)
3. FreeLang 코드에서 ELF 바이너리 파일을 저장할 수단이 없음

## 구현 전략 (4단계)

### Stage 0: C 인터프리터 빌드 확인
```bash
cd /tmp/freelang-c-final
make
./bin/fl run examples/hello_world.fl   # 검증
```

### Stage 1: stdlib에 write_bytes 기능 추가
파일: `/tmp/freelang-c-final/src/stdlib.c` (수정)

추가할 함수 (C):
```c
// write_bytes_file(filename, byte_array)
// 이진 파일 저장 - ELF 생성에 필요
fl_value_t fl_write_bytes_file(fl_value_t* args, size_t argc) {
    // args[0] = filename (string)
    // args[1] = byte array (array of integers 0-255)
    FILE* f = fopen(args[0].data.string_val, "wb");
    fl_array_t* arr = args[1].data.array_val;
    for (size_t i = 0; i < arr->size; i++) {
        uint8_t byte = (uint8_t)arr->elements[i].data.int_val;
        fwrite(&byte, 1, 1, f);
    }
    fclose(f);
    return fl_new_null();
}

// read_file(filename) -> string
fl_value_t fl_read_file(fl_value_t* args, size_t argc);
```

VM에 등록 위치: `/tmp/freelang-c-final/src/vm.c`
(이미 builtin 함수 처리 로직이 있음)

### Stage 2: hello.free 작성 (소스 파일)
```
// hello.free
fn main() {
  println("Hello, World!")
}
```

### Stage 3: FreeLang ELF 생성기 작성
파일: `self-elf-gen.fl` (FreeLang으로 작성)

이 파일이 핵심입니다. FreeLang 코드로 x86-64 ELF 바이너리를 생성합니다.

**최소한의 Linux ELF Hello World (약 136 bytes):**

```
ELF Header:     64 bytes
Program Header: 56 bytes
x86-64 코드:    ~28 bytes (syscall 기반)
문자열:         13 bytes ("Hello World!\n")
```

x86-64 기계어 (Hello World syscall):
```
48 c7 c0 01 00 00 00  ; mov rax, 1  (sys_write)
48 c7 c7 01 00 00 00  ; mov rdi, 1  (stdout)
48 8d 35 XX XX XX XX  ; lea rsi, [msg]
48 c7 c2 0d 00 00 00  ; mov rdx, 13 (길이)
0f 05                 ; syscall
48 c7 c0 3c 00 00 00  ; mov rax, 60 (sys_exit)
48 31 ff              ; xor rdi, rdi
0f 05                 ; syscall
48 65 6c 6c 6f ...    ; "Hello World!\n"
```

FreeLang으로 이 바이트 배열을 생성하고 `write_bytes_file("hello_world", bytes)`로 저장.

### Stage 4: 파이프라인 통합
파일: `self-compiler.fl` (FreeLang으로 작성)

```
입력:  hello.free (소스 파일)
처리:
  1. read_file("hello.free")  → 소스 문자열
  2. tokenize(source)          → 토큰 배열 (bootstrap-demo.fl 함수 재사용)
  3. 파싱                       → 함수 이름, println 인자 추출
  4. ELF 코드 생성              → 메시지를 ELF에 포함
  5. write_bytes_file(...)     → 실행 가능한 ELF 저장
출력:  hello_world (ELF 64-bit executable)
```

## 파일 수정 계획

### 수정 파일
1. `/tmp/freelang-c-final/src/stdlib.c`
   - `fl_write_bytes_file()` 함수 추가 (약 30줄)
   - `fl_read_file()` 함수 추가 (약 20줄)

2. `/tmp/freelang-c-final/src/vm.c`
   - builtin 함수 디스패치에 `write_bytes_file`, `read_file` 등록

### 신규 파일 (FreeLang)
3. `/tmp/freelang-c-final/hello.free` - 소스 입력
4. `/tmp/freelang-c-final/self-elf-gen.fl` - ELF 생성기
5. `/tmp/freelang-c-final/self-compiler.fl` - 통합 파이프라인

## 검증 기준 (절대 거짓 없음)

```
✅ Step 1: make 성공 → bin/fl 존재
✅ Step 2: ./bin/fl run examples/hello_world.fl → "Hello World!" 출력
✅ Step 3: ./bin/fl run self-elf-gen.fl → hello_world 파일 생성
✅ Step 4: file hello_world → "ELF 64-bit LSB executable, x86-64" 확인
✅ Step 5: chmod +x hello_world && ./hello_world → "Hello World!" 출력
✅ Step 6: ./bin/fl run self-compiler.fl hello.free → 동일 결과
```

**로그 캡처**: 모든 단계 출력을 기록.

## 빌드 환경 확인 필요
- Termux에 gcc 설치 여부 (또는 clang)
- make 명령 사용 가능 여부
- Linux ARM64 환경 (ELF 아키텍처 주의: x86-64 vs aarch64)

⚠️ **중요**: Termux는 ARM64 (aarch64)이므로 x86-64 ELF가 아닌 **aarch64 ELF**를 생성해야 실행 가능합니다.
