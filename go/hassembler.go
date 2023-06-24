package main

import "fmt"
import "time"
import "strings"
import "strconv"
// import "github.com/docopt/docopt.go"

/* Docopt args usage */
const usage = `
Usage: 
  hassembler <file> [<outfile>]
  hassembler -h | --help

Options:
  file        path of the input .asm file
  outfile     path of the output .bin file
  -h, --help  show this help message and exit`

/* Symbols */
var compSymbols map[string]int = map[string]int{ // C instruction symbols
	// a=0
	"0":    0b101010, "1": 0b111111, "-1": 0b111010, "D": 0b001100,
	"A":    0b110000, "!D": 0b001101, "!A": 0b110001, "-D": 0b001111,
	"-A":   0b110011, "D+1": 0b011111, "A+1": 0b110111, "D-1": 0b001110,
	"A-1":  0b110010, "D+A": 0b000010, "D-A": 0b010011, "A-D": 0b000111,
	"D&A":  0b000000, "D|A": 0b010101,
	// a=1
	"M":    0b111000, "!M": 0b111001, "-M": 0b111011, "M+1": 0b111111,
	"M-1":  0b111010, "D+M": 0b100001, "D-M": 0b101000, "M-D": 0b100011,
	"D&M":  0b100000, "D|M": 0b101101,
}
var jumpSymbols map[string]int = map[string]int{
	"JGT": 0b001, "JEQ": 0b010, "JGE": 0b011, "JLT": 0b100,
	"JNE": 0b101, "JLE": 0b0110, "JMP": 0b111, "": 0b0,
}
var destSymbols map[string]int = map[string]int{
	"M": 0b001, "D": 0b010, "MD": 0b011, "A": 0b100,
	"AM": 0b101, "AD": 0b110, "AMD": 0b111,
}
var predefinedSymbols map[string]int = map[string]int{ // A instruction symbols
	"R0": 0b0000, "R1": 0b0001, "R2": 0b0010, "R3": 0b0011,
	"R4": 0b0100, "R5": 0b0101, "R6": 0b0110, "R7": 0b0111,
	"R8": 0b1000, "R9": 0b1001, "R10": 0b1010, "R11": 0b1011,
	"R12": 0b1100, "R13": 0b1101, "R14": 0b1110, "R15": 0b1111,
	"SCREEN": 0b100000000000000, "KBD": 0b110000000000000, "SP": 0b0,
	"LCL": 0b001, "ARG": 0b010, "THIS": 0b011, "THAT": 0b100,
}
var variableSymbols map[string]int = map[string]int{}
var labelSymbols map[string]int = map[string]int{}

// Functions

/* Removes comments from a string and applies strings.TrimSpace() to it */
func removeComments(s* string) {
	if strings.ContainsAny(*s, "//") {
		*s = strings.TrimSpace(strings.Split(*s, "//")[0])
	}
}

/* Returns: 16-bit binary representation of an A instruction parsed from a string */
func parseAinstruction(line string) (string, error) {
	line = strings.ReplaceAll(line, "@", "")  // Remove @

	if val, ok := predefinedSymbols[line]; ok {  // Check if instruction points to a predefined symbol
		return fmt.Sprintf("0%015b", val), nil
	}
	if val, ok := labelSymbols[line]; ok {  // Check if instruction points to a label
		return fmt.Sprintf("0%015b", val), nil
	}
	if val, ok := variableSymbols[line]; ok {  // Check if instruction points to a variable
		return fmt.Sprintf("0%015b", val), nil
	}
	
	// If A instruction is not in symbols it's either a variable that we need to store or a simple int that points to an address
	address, err := strconv.Atoi(line)
	if err != nil {
		variableSymbols[line] = 16 + len(variableSymbols)
		return fmt.Sprintf("0%015b", variableSymbols[line]), nil
	}

	if address > 32768 {
		return "", fmt.Errorf("AddressOverflow: Address is bigger than 32768")
	}

	return fmt.Sprintf("0%015b", address), nil
}

/* Returns: 16-bit binary representation of a C instruction parsed from a string */
func parseCinstruction(line string) (string, error) {
	// Instruction of type: "dest=comp"
	if strings.ContainsAny(line, "="){
		instruction := strings.Split(line, "=")
		dest := instruction[0]
		comp := instruction[1]

		// Check if given dest and comp fields exists
		if _, ok := destSymbols[dest]; !ok {
			return "", fmt.Errorf("SymbolNotFound: Dest field doesn't exist")
		}
		if _, ok := compSymbols[comp]; !ok {
			return "", fmt.Errorf("SymbolNotFound: Comp field doesn't exist")
		}
		return fmt.Sprintf("111%07b%03b000", compSymbols[comp], destSymbols[dest]), nil
	}

	// Instruction of type: "comp;jump"
	if strings.ContainsAny(line, ";"){
		instruction := strings.Split(line, ";")
		comp := instruction[0]
		jump := instruction[1]

		// Check if given comp and jump fields exists
		if _, ok := compSymbols[comp]; !ok {
			return "", fmt.Errorf("SymbolNotFound: Comp field doesn't exist")
		}
		if _, ok := jumpSymbols[jump]; !ok {
			return "", fmt.Errorf("SymbolNotFound: Jump field doesn't exist")
		}
		return fmt.Sprintf("111%07b000%03b", compSymbols[comp], jumpSymbols[jump]), nil
	}

	return "", fmt.Errorf("ParseError: Couldn't parse C instruction")
}

// Main
func main() {
	startTime := time.Now()
	
	// Load cli args and file
	// TODO
	// args, _ := docopt.ParseDoc(usage)

	fmt.Printf("Took %.7f seconds.", time.Since(startTime).Seconds())
}