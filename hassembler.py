""" Assembler for the hack machine from nand2tetris.org. """
from argparse import ArgumentParser


def remove_comments(line: str) -> str:
    """ Removes single-line comments from a string """
    if "//" in line:
        return line.split("//")[0]
    return line

def parse_a_instruction(line: str) -> str:
    """ Parses an A instruction into its 16-bit binary representation """
    address = int(line.strip('@'))  # Strip the @ from the str and convert to int
    # TODO handle address overflow
    binary_instruction = f"0{address:015b}"  # Convert the int to 15-bit binary string, add a 0 at the start
    return binary_instruction


# Load the file from program args
argparser = ArgumentParser("hassembler")
argparser.add_argument("file", help="path of the .asm file", type=str)
args = argparser.parse_args()
binary_output = ""  # Save the output to be written later to a .bin file

# Parsing of the file line-by-line
with open(args.file, 'r') as file:
    while line := file.readline():
        line = remove_comments(line.strip())    # strip and remove comments
        if line == "":  # skip empty lines (also if line is many whitespaces and stripped it is empty)
            continue
        
        if line.startswith("@"):
            print(f"{parse_a_instruction(line)} - {line}")
        else:
            print(line)
