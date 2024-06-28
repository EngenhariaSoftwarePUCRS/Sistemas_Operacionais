package main

import (
	"fmt"
	"os"
	"strconv"
)

const (
	DEFAULT_V_MEM_SIZE = 64
	DEFAULT_F_MEM_SIZE = 32
	DEFAULT_PAGE_SIZE = 12

	// Console Editing
	BOLD = "\033[1m"
	RED = "\033[31m"
	GREEN = "\033[32m"
	YELLOW = "\033[33m"
	BLUE = "\033[34m"
	MAGENTA = "\033[35m"
	CYAN = "\033[36m"
	RESET = "\033[0m"
)

var (
	// Virtual Memory Size
	V_MEM_SIZE uint
	// Physical Memory Size
	F_MEM_SIZE uint
	// Page Size
	PAGE_SIZE uint

	// Virtual Memory
	V_MEM []uint
	// Physical Memory
	F_MEM []uint
	// Page Table
	PAGE_TABLE []uint
)

// Returns the argument at the given index or the default value if it was not provided
// Exits the program if the argument is not a non-negative integer
func getArg(index int, defaul uint) uint {
	if len(os.Args) > index {
		value, err := strconv.Atoi(os.Args[index])
		if err != nil || value < 0 {
			fmt.Printf(
				"%s Invalid argument %d: %d. Must be a non-negative integer. %s\n",
				RED, index, value, RESET,
			)
			os.Exit(1)
		}
		return uint(value)
	}
	return defaul
}

func main() {
	fmt.Println(RESET, BOLD)
	fmt.Println("=====", "Sistema Gerência de Memória Paginada", "=====")
	
	if len(os.Args) < 3 {
		fmt.Println("\nUseful arguments missing")
		fmt.Printf("Usage: go run SGMP.go %s <V = Virtual Memory Size> %s <F = Physical Memory Size> %s <P = Page Size>\n", CYAN, MAGENTA, GREEN)
		fmt.Print(RESET, BOLD)
		fmt.Printf("All sizes are in base 2 power. i.e. 2^x\n")
		fmt.Println("Example: go run SGMP.go 64 32 12")
		fmt.Print(RESET, BOLD)
		fmt.Println("Using default values to fill non-provided arguments")
	}
	V_MEM_SIZE = getArg(1, DEFAULT_V_MEM_SIZE)
	F_MEM_SIZE = getArg(2, DEFAULT_F_MEM_SIZE)
	PAGE_SIZE = getArg(3, DEFAULT_PAGE_SIZE)

	fmt.Println(CYAN, "\t", "Virtual Memory Size: ", "\t", V_MEM_SIZE)
	fmt.Println(MAGENTA, "\t", "Physical Memory Size: ", "", F_MEM_SIZE)
	fmt.Println(GREEN, "\t", "Page Size: ", "\t\t", PAGE_SIZE)
	fmt.Println(RESET)
}
