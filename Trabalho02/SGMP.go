package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
)

const (
	DEFAULT_V_MEM_SIZE = 16
	DEFAULT_F_MEM_SIZE = 15
	DEFAULT_PAGE_SIZE = 12
	DEFAULT_SEED = 42
	DEFAULT_V_ADDRS_COUNT = 5

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

	// Physical Memory
	F_MEM []uint
	// Page Table
	PAGE_TABLE []int

	// Random Seed
	RAND_SEED uint
	// Virtual Addresses
	V_ADDRS []uint
)

// Returns the argument at the given index or the default value if it was not provided
// Exits the program if the argument is not a non-negative integer
func getArg(index int, defaul uint) uint {
	if len(os.Args) > index {
		value, err := strconv.Atoi(os.Args[index])
		if err != nil || value < 0 {
			fmt.Println(colorize(RED, fmt.Sprintf(
				"Invalid argument %d: %s. Must be a non-negative integer.", index, os.Args[index])))
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

// Returns a string formatted to be printed with color
func colorize(color string, text string) string {
	return fmt.Sprintf("%s%s%s", color, text, RESET)
}

// Handles the arguments passed to the program
// Exits the program if the arguments are invalid
func handleArgs() {
	fmt.Println()

	if len(os.Args) < 6 {
		fmt.Println("\nUseful arguments missing")
		fmt.Printf("Usage: go run SGMP.go %s <V = Virtual Memory Size> %s <F = Physical Memory Size> %s <P = Page Size> %s <Optional: S = Random Seed> %s <Optional: A = The Amount of Virtual Addresses to Generate if One, the Virtual Addresses if More>\n", CYAN, MAGENTA, GREEN, BLUE, YELLOW)
		fmt.Println(RESET, BOLD)
		fmt.Printf("All sizes are in 2^n (bits). i.e. 2^%d = %s = %s\n",
			10, formatMemory(pow2(10), "bits"), formatMemory(bitsToBytes(pow2(10)), "bytes"))
		fmt.Printf("Example: go run SGMP.go %s %s %s %s %s\n",
			colorize(CYAN, fmt.Sprint(DEFAULT_V_MEM_SIZE)),
			colorize(MAGENTA, fmt.Sprint(DEFAULT_F_MEM_SIZE)),
			colorize(GREEN, fmt.Sprint(DEFAULT_PAGE_SIZE)),
			colorize(BLUE, fmt.Sprint(DEFAULT_SEED)),
			colorize(YELLOW, fmt.Sprint(DEFAULT_V_ADDRS_COUNT)),
		)
		fmt.Print(BOLD)
		fmt.Println("Using default values to fill non-provided arguments")
	}
	
	V_MEM_SIZE = getArg(1, DEFAULT_V_MEM_SIZE)
	F_MEM_SIZE = getArg(2, DEFAULT_F_MEM_SIZE)
	PAGE_SIZE = getArg(3, DEFAULT_PAGE_SIZE)
	RAND_SEED = getArg(4, DEFAULT_SEED)
	randomizer := rand.New(rand.NewSource(int64(RAND_SEED)))

	fmt.Println(colorize(CYAN,
		fmt.Sprintf("\tVirtual Memory Size:\t2^%d = %s = %s", V_MEM_SIZE,
			formatMemory(pow2(V_MEM_SIZE), "bits"), formatMemory(bitsToBytes(pow2(V_MEM_SIZE)), "bytes")),
	))

	fmt.Println(colorize(MAGENTA,
		fmt.Sprintf("\tPhysical Memory Size:\t2^%d = %s = %s", F_MEM_SIZE,
			formatMemory(pow2(F_MEM_SIZE), "bits"), formatMemory(bitsToBytes(pow2(F_MEM_SIZE)), "bytes")),
	))

	fmt.Println(colorize(GREEN,
		fmt.Sprintf("\tPage Size:\t\t2^%d = %s = %s", PAGE_SIZE,
			formatMemory(pow2(PAGE_SIZE), "bits"), formatMemory(bitsToBytes(pow2(PAGE_SIZE)), "bytes")),
	))

	fmt.Println(RESET)
	if len(os.Args) > 6 {
		V_ADDRS = make([]uint, len(os.Args) - 5)
		for i := 5; i < len(os.Args); i++ {
			value, err := strconv.Atoi(os.Args[i])
			if err != nil || value < 0 {
				fmt.Println(colorize(RED, fmt.Sprintf(
					"Invalid argument %d: %s. Must be a non-negative integer.", i, os.Args[i])))
				os.Exit(1)
			}
			V_ADDRS[i-5] = uint(value)
		}
		fmt.Printf("Using the following arguments as virtual addresses: %d\n", V_ADDRS)
	} else {
		V_ADDRS_COUNT := getArg(5, DEFAULT_V_ADDRS_COUNT)
		fmt.Printf("Using %s as random seed to generate %s virtual addresses\n",
			colorize(BLUE, fmt.Sprint(RAND_SEED)), colorize(YELLOW, fmt.Sprint(V_ADDRS_COUNT)))
		V_ADDRS = make([]uint, V_ADDRS_COUNT)
		for i := 0; i < len(V_ADDRS); i++ {
			upper_limit := uint32(pow2(V_MEM_SIZE) - 1)
			V_ADDRS[i] = uint(randomizer.Uint32() % upper_limit)
		}
		fmt.Printf("Generated virtual addresses: %d\n", V_ADDRS)
	}
}

func setupTables() {
	fmt.Println()

	frames_count := pow2(F_MEM_SIZE - PAGE_SIZE)
	fmt.Printf("Page Table has %d page frames of %s (%s) each\n", frames_count,
		formatMemory(pow2(PAGE_SIZE), "bits"), formatMemory(bitsToBytes(pow2(PAGE_SIZE)), "bytes"))
	F_MEM = make([]uint, frames_count)
	PAGE_TABLE = make([]int, len(F_MEM))

	for i := 0; i < len(F_MEM); i++ {
		F_MEM[i] = 0
		PAGE_TABLE[i] = -1
	}

	fmt.Printf("Physical Memory and Page Table initialized with %ds and %ds, respectively\n", 0, -1)
}

func main() {
	fmt.Println(RESET, BOLD)
	fmt.Println("=====", "Sistema Gerência de Memória Paginada", "=====")

	handleArgs()
	setupTables()

	fmt.Println(RESET)
}
