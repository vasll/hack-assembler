package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/docopt/docopt.go"
)

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
    "0": 42,    "1": 63,    "-1": 58,   "D": 12,
    "A": 48,    "!D": 13,   "!A": 49,   "-D": 15,
    "-A": 51,   "D+1": 31,  "A+1": 55,  "D-1": 14,
    "A-1": 50,  "D+A": 2,   "D-A": 19,  "A-D": 7,
    "D&A": 0,   "D|A": 21,
    // a=1
    "M": 112,   "!M": 113, "-M": 115,  "M+1": 119,
    "M-1": 114, "D+M": 66, "D-M": 83, "M-D": 71,
    "D&M": 64,  "D|M": 85,
}
var jumpSymbols map[string]int = map[string]int{
    "JGT": 1,   "JEQ": 2,   "JGE": 3,   "JLT": 4,
    "JNE": 5,   "JLE": 6,   "JMP": 7,   "": 0,
}
var destSymbols map[string]int = map[string]int{
    "M": 1,     "D": 2,     "MD": 3,    "A": 4,
    "AM": 5,    "AD": 6,    "AMD": 7,
}
var predefinedSymbols map[string]int = map[string]int{ // A instruction symbols
    "R0": 0,    "R1": 1,    "R2": 2,    "R3": 3,  
    "R4": 4,    "R5": 5,    "R6": 6,    "R7": 7,   
    "R8": 8,    "R9": 9,    "R10": 10,  "R11": 11, 
    "R12": 12,  "R13": 13,  "R14": 14,  "R15": 15, 
    "SCREEN": 16384,    "KBD": 24576,   "SP":0,  
    "LCL": 1,   "ARG": 2,   "THIS": 3,  "THAT": 4,
}
var variableSymbols map[string]int = map[string]int{}
var labelSymbols map[string]int = map[string]int{}

// Functions
/* Removes comments from a string and applies strings.TrimSpace() to it */
func removeComments(s string) string{
	if strings.ContainsAny(s, "//") {
		return strings.TrimSpace(strings.Split(s, "//")[0])
	}
	return strings.TrimSpace(s)
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
		return "", fmt.Errorf("AddressOverflow")
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
			return "", fmt.Errorf("SymbolNotFound")
		}
		if _, ok := compSymbols[comp]; !ok {
			return "", fmt.Errorf("SymbolNotFound")
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
			return "", fmt.Errorf("SymbolNotFound")
		}
		if _, ok := jumpSymbols[jump]; !ok {
			return "", fmt.Errorf("SymbolNotFound")
		}
		return fmt.Sprintf("111%07b000%03b", compSymbols[comp], jumpSymbols[jump]), nil
	}

	return "", fmt.Errorf("ParseError")
}

// Main
func main() {
	startTime := time.Now()	// Keep track of execution time
	
	// Load cli args
	args, _ := docopt.ParseDoc(usage)

	// Open file for reading
	file, file_err := os.Open(args["<file>"].(string))
	if file_err != nil {
		fmt.Printf("Error with opnening file.")
		os.Exit(-1)
	}
	scanner := bufio.NewScanner(file)

	// First pass: Find label symbols
	instructionNumber := 0
	
	for scanner.Scan() {	// Read file line by line
		line := removeComments(scanner.Text())	
		if line == "" {	continue }	// Skip empty lines
		if strings.HasPrefix(line, "(") && strings.HasSuffix(line, ")") {
			label := strings.ReplaceAll(line, "(", "")
			label = strings.ReplaceAll(label, ")", "")
			labelSymbols[label] = instructionNumber
		}
		instructionNumber += 1
	}
	file.Seek(0, 0)

	// Second pass: Find variable symbols
	var outInstructions []string = []string{}
	rawLineCount := 0
	errorCount := 0

	for scanner.Scan() {	// Read file line by line
		line := scanner.Text()	
		strpLine := removeComments(line)
		if strpLine == "" || strings.HasPrefix(strpLine, "(") && strings.HasSuffix(strpLine, ")"){
			rawLineCount += 1
			continue
		}

		if strings.HasPrefix(strpLine, "@"){
			instr, err := parseAinstruction(strpLine)
			
			if err != nil {
				errorCount += 1
				fmt.Printf("┌%s at line %d\n", err.Error(), rawLineCount)
				fmt.Printf("│ line %d: '%s'\n", rawLineCount, strings.TrimSpace(line))
				continue
			}

			outInstructions = append(outInstructions, instr)
		} else {
			instr, err := parseCinstruction(strpLine)
			
			if err != nil {
				errorCount += 1
				fmt.Printf("┌%s at line %d\n", err.Error(), rawLineCount)
				fmt.Printf("│ line %d: '%s'\n", rawLineCount, strings.TrimSpace(line))
				continue
			}

			outInstructions = append(outInstructions, instr)
		}
		rawLineCount += 1
	}
	file.Close()

	// Decide whether to write the output to file or quit if there are errors
	if errorCount > 0 {
		fmt.Printf("Found %d errors. Exiting\n", errorCount)
		os.Exit(-1)
	}

	// open output file
	outfile, outfile_err := os.Create(args["<outfile>"].(string))
	if outfile_err != nil {
		fmt.Printf("Error with opnening output file.")
		os.Exit(-1)
	}
	defer file.Close()

	for _, instruction := range outInstructions{
		outfile.WriteString(instruction)
		outfile.WriteString("\n")
	}

	fmt.Printf("File written to '%s'\n", args["<outfile>"].(string))
	fmt.Printf("Took %.7f seconds.\n", time.Since(startTime).Seconds())
}