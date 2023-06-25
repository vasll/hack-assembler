# hack-assembler
A hack machine language assembler made for the [nand2tetris](https://www.nand2tetris.org) course with python and Go.

# Usage
```powershell
usage: hassembler [-h] file [outfile]

positional arguments:
  file        path of the input .asm file
  outfile     path of the output .bin file

options:
  -h, --help  show this help message and exit
```

# How does it work?
The hack platform from [nand2tetris](https://www.nand2tetris.org) has two basic A instructions and C instructions that need to be parsed into binary.\
The parsing of these instructions is done with different symbol tables consisting of a symbol and a binary value. Each line is read from the file and split into its various fields, then each field is parsed using the symbols table and finally written to the output file.

![image](https://github.com/vasll/hack-assembler/assets/67590845/d3d73e2f-3bc8-44dc-9fdc-7f071e33d3ed)

## A instructions
They always start with a '@' and change the value of the Address register, they can point to:
- A predefined symbol like a register `@R15, @R2, @THIS, @KBD, @SCREEN, ...`
- A label `@main_loop, ...`
- A memory address `@1254, @56, ...`
## C instructions
They can be of two types and perform computation on the ALU like adding, subtracting, anding, comparisons, flipping bits and so on:
- "dest=comp" `D=M, D=D-A, ...`
- "comp;jump" `0;JMP, D;JGT, ...`


