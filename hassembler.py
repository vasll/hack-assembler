""" Assembler for the hack machine from nand2tetris.org. """
from argparse import ArgumentParser


# Load the file from program args
argparser = ArgumentParser("hassembler")
argparser.add_argument("file", help="path of the .asm file", type=str)
args = argparser.parse_args()
binary_output = ""  # Save the output to be written later to a .bin file

# Symbols
predefined_symbols = {
    "R0": 0,    "R1": 1,    "R2": 2,    "R3": 3,  
    "R4": 4,    "R5": 5,    "R6": 6,    "R7": 7,   
    "R8": 8,    "R9": 9,    "R10": 10,  "R11": 11, 
    "R12": 12,  "R13": 13,  "R14": 14,  "R15": 15, 
    "SCREEN": 16384,    "KBD": 24576,   "SP":0,  
    "LCL": 1,   "ARG": 2,   "THIS": 3,  "THAT": 4,
} 
variable_symbols = {}
label_symbols = {}  # Will contain labels with this format: {"LABEL_NAME": "LINE_COUNT"}

# Functions
def remove_comments(line: str) -> str:
    """ Removes single-line comments from a string """
    if "//" in line:
        return line.split("//")[0]
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
    
    if address > 24576:
        print(f"ERROR: Address is too big at line {raw_line_count} -> {line}")
        exit()
    
    binary_instruction = f"0{address:015b}"  # Convert the int to 15-bit binary string, add a 0 at the start
    return binary_instruction


# Main program - Parsing of the file line-by-line
output = ""

with open(args.file, 'r') as file:  # Read file into "lines" variable
    lines = file.readlines()

# First iteration: Find labels and store them
instruction_number = 0

for line in lines:
    line = remove_comments(line.strip())    # Strip and remove comment
    if line == "":  # Skip empty lines
        continue

    if line.startswith("(") and line.endswith(")"):
        label = line.strip("()")
        label_symbols[label] = instruction_number

    instruction_number += 1

# Second iteration: The big parsing
raw_line_count = 0
instruction_number = 0

for line in lines:
    line = remove_comments(line.strip())    # Strip and remove comments
    if line == "":  # Skip empty lines
        raw_line_count += 1
        continue
    
    if line.startswith("@"):
        print(f"{parse_a_instruction(line)} -> {line}")
    else:
        print(f"? -> {line}")
        pass
    
    raw_line_count += 1
    instruction_number += 1
