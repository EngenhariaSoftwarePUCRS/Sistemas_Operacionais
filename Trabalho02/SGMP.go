package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

const (
	DEFAULT_V_MEM_SIZE = 16
	DEFAULT_F_MEM_SIZE = 15
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

// Calculates 2 to the power of the given number and returns it
func pow2(n uint) uint {
	return uint(math.Pow(2, float64(n)))
}

// Converts bits to bytes
func bitsToBytes(bits uint) uint {
	return bits / 8
}

// Returns the given number of memory formatted as Kilo, Mega, Giga or Tera {unit}
// @param unit: bits or bytes
func formatMemory(memory uint, unit string) string {
	var ending string
	if unit == "bits" {
		ending = "b"
	} else if unit == "bytes" {
		ending = "B"
	} else {
		os.Exit(1)
	}

	if memory < pow2(10) {
		return fmt.Sprintf("%d %s", int(memory), unit)
	}
	if memory < pow2(20) {
		return fmt.Sprintf("%d K%s", memory/1024, ending)
	}
	if memory < pow2(30) {
		return fmt.Sprintf("%d M%s", memory/1024/1024, ending)
	}
	if memory < pow2(40) {
		return fmt.Sprintf("%d G%s", memory/1024/1024/1024, ending)
	}
	return fmt.Sprintf("%d T%s", memory/1024/1024/1024/1024, ending)
}

func main() {
	fmt.Println(RESET, BOLD)
	fmt.Println("=====", "Sistema Gerência de Memória Paginada", "=====")
	
	if len(os.Args) < 3 {
		fmt.Println("\nUseful arguments missing")
		fmt.Printf("Usage: go run SGMP.go %s <V = Virtual Memory Size> %s <F = Physical Memory Size> %s <P = Page Size>\n", CYAN, MAGENTA, GREEN)
		fmt.Println(RESET, BOLD)
		fmt.Printf("All sizes are in 2^n (bits). i.e. 2^%d = %s = %s\n",
			10, formatMemory(pow2(10), "bits"), formatMemory(bitsToBytes(pow2(10)), "bytes"))
		fmt.Println("Example: go run SGMP.go 64 32 12")
		fmt.Print(RESET, BOLD)
		fmt.Println("Using default values to fill non-provided arguments")
	}
	V_MEM_SIZE = getArg(1, DEFAULT_V_MEM_SIZE)
	F_MEM_SIZE = getArg(2, DEFAULT_F_MEM_SIZE)
	PAGE_SIZE = getArg(3, DEFAULT_PAGE_SIZE)

	fmt.Printf("%s\tVirtual Memory Size:\t2^%d = %s = %s\n", CYAN, V_MEM_SIZE,
		formatMemory(pow2(V_MEM_SIZE), "bits"), formatMemory(bitsToBytes(pow2(V_MEM_SIZE)), "bytes"))

	fmt.Printf("%s\tPhysical Memory Size:\t2^%d = %s = %s\n", MAGENTA, F_MEM_SIZE,
		formatMemory(pow2(F_MEM_SIZE), "bits"), formatMemory(bitsToBytes(pow2(F_MEM_SIZE)), "bytes"))

	fmt.Printf("%s\tPage Size:\t\t2^%d = %s = %s\n", GREEN, PAGE_SIZE,
		formatMemory(pow2(PAGE_SIZE), "bits"), formatMemory(bitsToBytes(pow2(PAGE_SIZE)), "bytes"))
	fmt.Println(RESET)
}
