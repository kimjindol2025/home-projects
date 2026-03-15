/* FreeLang OS Kernel 링커 스크립트 */

ENTRY(_start)

SECTIONS {
    /* 커널은 1MB부터 시작 */
    . = 1M;

    .boot :
    {
        /* Multiboot2 헤더 (맨 앞) */
        KEEP(*(.multiboot_header))
        *(.boot)
    }

    .text :
    {
        *(.text .text.*)
        . = ALIGN(4K);
    }

    .rodata :
    {
        *(.rodata .rodata.*)
        . = ALIGN(4K);
    }

    .data :
    {
        *(.data .data.*)
        . = ALIGN(4K);
    }

    .bss :
    {
        *(.bss .bss.*)
        *(COMMON)
        . = ALIGN(4K);
    }

    /DISCARD/ :
    {
        *(.eh_frame)
        *(.comment)
    }
}
