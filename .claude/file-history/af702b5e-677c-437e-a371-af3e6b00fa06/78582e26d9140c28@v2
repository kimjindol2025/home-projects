/* FreeLang Native AOT Compiler - GNU Linker Script
   Phase 4: Memory layout definition for x86-64 Linux ELF */

ENTRY(_start)

SECTIONS
{
    /* Load address starts at 0x400000 */
    . = 0x400000;

    /* Program header */
    PHDRS
    {
        headers PT_PHDR PHDRS;
        text    PT_LOAD;
        data    PT_LOAD;
        dynamic PT_DYNAMIC;
    }

    /* ELF header and program headers */
    .elf_header :
    {
        *(.elf_header)
    }

    /* Executable code section */
    .text ALIGN(0x1000) : AT(ADDR(.text))
    {
        *(.text)
        *(.text.*)
        *(.rodata)
        *(.rodata.*)
        . = ALIGN(0x1000);
    } : text

    /* Initialized data section */
    .data ALIGN(0x1000) : AT(ADDR(.data))
    {
        *(.data)
        *(.data.*)
        . = ALIGN(8);
    } : data

    /* Uninitialized data section (BSS) */
    .bss ALIGN(0x1000) :
    {
        *(.bss)
        *(.bss.*)
        *(COMMON)
        . = ALIGN(0x1000);
    } : data

    /* Symbol table */
    .symtab ALIGN(8) :
    {
        *(.symtab)
    }

    /* String table */
    .strtab ALIGN(8) :
    {
        *(.strtab)
    }

    /* Symbol string table */
    .shstrtab ALIGN(8) :
    {
        *(.shstrtab)
    }

    /* Relocation sections */
    .rela.text ALIGN(8) :
    {
        *(.rela.text)
        *(.rela.text.*)
    }

    .rela.data ALIGN(8) :
    {
        *(.rela.data)
        *(.rela.data.*)
    }

    /* Section headers come at the end */
    /DISCARD/ : { *(.note.GNU-stack) }
}

/* Additional symbols for runtime */
PROVIDE(_start = 0x400000);
PROVIDE(__data_start = ADDR(.data));
PROVIDE(__data_end = ADDR(.data) + SIZEOF(.data));
PROVIDE(__bss_start = ADDR(.bss));
PROVIDE(__bss_end = ADDR(.bss) + SIZEOF(.bss));
PROVIDE(__heap_start = ADDR(.bss) + SIZEOF(.bss));
PROVIDE(__text_start = ADDR(.text));
PROVIDE(__text_end = ADDR(.text) + SIZEOF(.text));
