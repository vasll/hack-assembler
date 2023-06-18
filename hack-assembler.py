""" Assembler for the hack machine from nand2tetris.org """

# 0. Load the file from program args
pass

# 1. Remove whitespaces and comments
pass

# 2. Translate instructions
def parse_a_instruction(line: str) -> str:
    """ Parses an A instruction into its 16-bit binary representation """
    address = int(line.strip('@'))  # Strip the @ from the str and convert to int
    binary_instruction = f"0{address:015b}"  # Convert the int to 15-bit binary string, add a 0 at the start
    return binary_instruction

