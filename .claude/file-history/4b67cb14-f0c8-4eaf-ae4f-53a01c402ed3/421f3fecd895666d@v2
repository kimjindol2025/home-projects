; FreeLang OS Kernel - Multiboot2 부트로더
; x86-64 bare-metal 부팅 코드

[bits 32]

; Multiboot2 헤더
MULTIBOOT2_HEADER_MAGIC     equ 0xe85250d6
MULTIBOOT_ARCHITECTURE_I386 equ 0
MULTIBOOT2_HEADER_TAG_END   equ 0

; 섹션: Multiboot2 헤더
section .multiboot_header
align 8
multiboot_header_start:
    dd MULTIBOOT2_HEADER_MAGIC
    dd MULTIBOOT_ARCHITECTURE_I386
    dd (multiboot_header_end - multiboot_header_start)
    dd -(MULTIBOOT2_HEADER_MAGIC + MULTIBOOT_ARCHITECTURE_I386 + (multiboot_header_end - multiboot_header_start))

    ; 선택 태그 (없음)
    dw 0    ; tag type (end)
    dw 0    ; flags
    dd 8    ; size

multiboot_header_end:

; 섹션: .text (코드)
section .text
extern kernel_main

global _start
_start:
    mov esp, stack_top          ; 스택 포인터 설정
    mov edi, eax                ; multiboot magic 저장 (EDI)
    mov esi, ebx                ; multiboot info pointer (ESI)

    ; CPU 기능 확인
    call check_cpuid
    call check_long_mode

    ; 페이지 테이블 설정
    call setup_page_tables

    ; 페이징 활성화
    call enable_paging

    ; GDT 로드
    lgdt [gdt64.pointer]

    ; 64-bit 코드 세그먼트로 점프
    jmp gdt64.code:long_mode_start

; CPUID 확인
check_cpuid:
    pushfd
    pop eax
    mov ecx, eax
    xor eax, 0x00200000        ; ID 플래그 토글
    push eax
    popfd
    pushfd
    pop eax
    cmp eax, ecx
    je .no_cpuid                ; CPUID 미지원
    ret
.no_cpuid:
    mov al, "C"
    jmp error

; Long mode 확인
check_long_mode:
    mov eax, 0x80000000
    cpuid
    cmp eax, 0x80000001
    jb .no_long_mode

    mov eax, 0x80000001
    cpuid
    test edx, 1 << 29           ; LM 플래그
    jz .no_long_mode
    ret
.no_long_mode:
    mov al, "L"
    jmp error

; 페이지 테이블 설정 (4KB 페이지)
setup_page_tables:
    ; P4 테이블 초기화
    mov eax, p3_table
    or eax, 0b11                ; Present + Writable
    mov [p4_table], eax

    ; P3 테이블 초기화
    mov eax, p2_table
    or eax, 0b11
    mov [p3_table], eax

    ; P2 테이블: 2MB 페이지로 0 ~ 4GB 매핑
    mov ecx, 0                  ; 카운터
.map_p2:
    mov eax, 0x200000          ; 2MB 페이지 시작
    mul ecx
    or eax, 0b10000011         ; Present + Writable + 2MB 페이지
    mov [p2_table + ecx * 8], eax

    inc ecx
    cmp ecx, 512
    jne .map_p2

    ret

; 페이징 활성화
enable_paging:
    ; CR3 = P4 테이블 주소
    mov eax, p4_table
    mov cr3, eax

    ; CR4: PAE 활성화
    mov eax, cr4
    or eax, 1 << 5              ; PAE 플래그
    mov cr4, eax

    ; EFER: Long mode 활성화
    mov ecx, 0xC0000080
    rdmsr
    or eax, 1 << 8              ; LME 플래그
    wrmsr

    ; CR0: 페이징 활성화
    mov eax, cr0
    or eax, 1 << 31             ; PG 플래그
    or eax, 1 << 16             ; WP 플래그 (write-protect)
    mov cr0, eax

    ret

; 에러 처리: 오류 코드를 표시 후 정지
error:
    ; VGA 텍스트 모드: 0xB8000
    ; 형식: [배경색|글자색][ASCII 코드]
    mov dword [0xb8000], 0x4f524f45    ; "ER" (빨간 배경)
    mov dword [0xb8004], 0x4f3a4f52    ; "R:"
    mov dword [0xb8008], 0x4f204f20    ; "  "
    mov byte [0xb800a], al              ; 에러 코드

    hlt
    jmp error

; GDT (Global Descriptor Table)
gdt64:
    dq 0                        ; 널 디스크립터
.code: equ $ - gdt64
    dq (1<<43) | (1<<47) | (1<<53)  ; 코드 세그먼트 (64-bit)
.data: equ $ - gdt64
    dq (1<<44) | (1<<47)        ; 데이터 세그먼트
.pointer:
    dw $ - gdt64 - 1            ; 크기
    dq gdt64                    ; 주소

; BSS 섹션: 초기화되지 않은 데이터
section .bss
align 4096

; 페이지 테이블 (각 4KB)
p4_table:
    resq 512

p3_table:
    resq 512

p2_table:
    resq 512

; 스택 (16KB)
stack_bottom:
    resb 4096 * 4
stack_top:

; 64-bit 모드 진입 코드
section .text.64bit
[bits 64]

extern kernel_main

global long_mode_start
long_mode_start:
    ; 세그먼트 레지스터 초기화
    mov ax, 0                   ; 널 세그먼트
    mov ss, ax
    mov ds, ax
    mov es, ax
    mov fs, ax
    mov gs, ax

    ; 스택 포인터 다시 설정
    mov rsp, stack_top

    ; 커널 메인 함수 호출
    mov rdi, rsi                ; multiboot info pointer (RDI)
    call kernel_main

    ; 반복 정지
    hlt
    jmp $
