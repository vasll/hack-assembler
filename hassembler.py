""" Assembler for the hack machine from nand2tetris.org. """
from argparse import ArgumentParser


# Load the file from program args
argparser = ArgumentParser("hassembler")
argparser.add_argument("file", help="path of the .asm file", type=str)
args = argparser.parse_args()
binary_output = ""  # Save the output to be written later to a .bin file

# Symbols
comp_symbols = {    # C instruction symbols
    # a=0
    "0": 42,    "1": 63,    "-1": 58,   "D": 12,
    "A": 48,    "!D": 13,   "!A": 49,   "-D": 15,
    "-A": 51,   "D+1": 31,  "A+1": 55,  "D-1": 14,
    "A-1": 50,  "D+A": 2,   "D-A": 19,  "A-D": 7,
    "D&A": 0,   "D|A": 21,
    # a=1
    "M": 112,   "!M": 113, "-M": 115,  "M+1": 119,
    "M-1": 114, "D+M": 66, "D-M": 83, "M-D": 71,
    "D&M": 64,  "D|M": 85
}
jump_symbols = {
    "JGT": 1,   "JEQ": 2,   "JGE": 3,   "JLT": 4,
    "JNE": 5,   "JLE": 6,   "JMP": 7,   "": 0
}
dest_symbols = {
    "M": 1,     "D": 2,     "MD": 3,    "A": 4,
    "AM": 5,    "AD": 6,    "AMD": 7
}

predefined_symbols = {  # A instruction symbols
    "R0": 0,    "R1": 1,    "R2": 2,    "R3": 3,  
    "R4": 4,    "R5": 5,    "R6": 6,    "R7": 7,   
    "R8": 8,    "R9": 9,    "R10": 10,  "R11": 11, 
    "R12": 12,  "R13": 13,  "R14": 14,  "R15": 15, 
    "SCREEN": 16384,    "KBD": 24576,   "SP":0,  
    "LCL": 1,   "ARG": 2,   "THIS": 3,  "THAT": 4,
} 
variable_symbols = {}
label_symbols = {}

# Functions
def remove_comments(line: str) -> str:
    """ Removes single-line comments from a string """
    if "//" in line:
        return line.split("//")[0].strip()
    return line

def parse_a_instruction(line: str) -> str:
    """ Parses an A instruction into its 16-bit binary representation """
    instruction = line.strip('@')  # Strip @
    
    if instruction in predefined_symbols:  # Check if instruction points to a predefined symbol
        return f"0{predefined_symbols[instruction]:015b}"
    if instruction in label_symbols:       # Check if instruction points to a label
        return f"0{label_symbols[instruction]:015b}"
    if instruction in variable_symbols:    # Check if instruction points to a variable
        return f"0{variable_symbols[instruction]:015b}"

    # If A instruction is not in symbols it's either a variable that we need to store or a simple int that points to an address
    try:
        address = int(instruction)  # Check if it's an int, if it's not it means it's a variable
    except ValueError:
        variable_symbols[instruction] = 16 + len(variable_symbols)  # Assign a unique memory address starting at 16
        return f"0{variable_symbols[instruction]:015b}"
    
    if address > 32768:
        print(f"ERROR: Address is too big -> {line}")
        exit()
    
    binary_instruction = f"0{address:015b}"  # Convert the int to 15-bit binary string, add a 0 at the start
    return binary_instruction

def parse_c_instruction(line: str) -> str:
    """ Parses a C instruction into its 16-bit binary representation """
    # Instruction of type: "dest=comp"
    if "=" in line:
        instruction = line.split("=")
        dest = instruction[0]
        comp = instruction[1]

        if dest not in dest_symbols or comp not in comp_symbols:    # Check dest field
            print(f"ERROR: Wrong C instruction format -> {line}")
            exit()

        return f"111{comp_symbols[comp]:07b}{dest_symbols[dest]:03b}000"

    # Instruction of type: "comp;jump"
    if ";" in line: 
        instruction = line.split(";")
        comp = instruction[0]
        jump = instruction[1]

        if jump not in jump_symbols or comp not in comp_symbols:    # Check dest field
            print(f"ERROR: Wrong C instruction format -> {line}")
            exit()
        
        return f"111{comp_symbols[comp]:07b}000{jump_symbols[jump]:03b}"
    
    print(f"ERROR: Unknown C instruction -> {line}")
    exit()


# Main program - Parsing of the file line-by-line
output = ""

with open(args.file, 'r') as file:  # Read file into "lines" variable
    lines = file.readlines()

# First pass: Find label symbols
instruction_number = 0

for line in lines:
    line = remove_comments(line.strip())    # Strip and remove comment
    if line == "":  # Skip empty lines
        continue

    if line.startswith("(") and line.endswith(")"):
        label = line.strip("()")
        label_symbols[label] = instruction_number

    instruction_number += 1

# Second pass: Find variable symbols
raw_line_count = 0
instruction_number = 0

for line in lines:
    line = remove_comments(line.strip())    # Strip and remove comments
    if line == "" or line.startswith("(") and line.endswith(")"):  # Skip empty and label lines
        raw_line_count += 1
        continue
    
    if line.startswith("@"):
        print(f"{parse_a_instruction(line)} -> {line}")
    else:
        print(f"{parse_c_instruction(line)} -> {line}")
        pass
    
    raw_line_count += 1
    instruction_number += 1
