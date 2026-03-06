; FreeLang Native AOT Compiler - Bootstrap Assembly
; Phase 4: Minimal bootstrap for x86-64

[bits 64]

; ELF entry point (0x400000)
_start:
    ; Initialize RSP (stack pointer)
    mov rsp, 0x7fffff00        ; Set stack at high memory
    xor rbp, rbp               ; Clear base pointer

    ; Save argc/argv (passed by kernel)
    ; rdi = argc
    ; rsi = argv

    ; Call main
    call main

    ; Exit with return code
    mov rax, 60                ; SYS_exit
    mov rdi, 0                 ; exit code = 0
    syscall

    ; Unreachable, but safety
    jmp $

; main function stub
main:
    ; Function prologue
    push rbp
    mov rbp, rsp
    sub rsp, 0                 ; No local variables

    ; Hello world syscall (write)
    ; rax = 1 (SYS_write)
    ; rdi = 1 (stdout)
    ; rsi = message address
    ; rdx = message length

    mov rax, 1                 ; SYS_write
    mov rdi, 1                 ; stdout
    lea rsi, [rel msg]         ; Load message address (RIP-relative)
    mov rdx, msg_len           ; Message length
    syscall

    ; Return 0
    xor eax, eax

    ; Function epilogue
    leave
    ret

; Read-only data section
section .rodata
msg:
    db "FreeLang Native Binary!", 0x0A
msg_len equ $ - msg
