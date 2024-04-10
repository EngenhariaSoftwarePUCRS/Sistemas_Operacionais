package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

const (
	COOKS_COUNT = 1
	SAVAGES_COUNT = 5
	THREAD_COUNT = COOKS_COUNT + SAVAGES_COUNT
	SERVINGS_COUNT = 5

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
	servings = SERVINGS_COUNT
	mutex *Mutex = NewMutex()
	emptyPot = NewSemaphore(0)
	fullPot = NewSemaphore(0)
)

// Returns a random integer between 0 and max
func GetRandom(max int) int {
	return rand.Intn(max + 1)
}

// Waits for a random time in milliseconds with 10x a custom multiplier
func WaitRandom(multiplier int) {
	time.Sleep(time.Duration(GetRandom(10 * multiplier)) * time.Millisecond)
}

func Cook(id int) {
	for {
		emptyPot.Wait()
		putServingsInPot(id)
		servings = SERVINGS_COUNT
		fullPot.Signal()
		WaitRandom(COOKS_COUNT)
	}
}

func putServingsInPot(id int) {
	PrintCook(fmt.Sprint("Cook ", id, " is putting servings in pot...", "\n"))
}

func Savage(id int) {
	for {
		WaitRandom(SAVAGES_COUNT)
		mutex.Lock(id)
			if servings == 0 {
				wakeUpCook(id)
				emptyPot.Signal()
				fullPot.Wait()
			}
			servings--
			getServingFromPot(id)
		mutex.Unlock(id)
	}
}

func wakeUpCook(id int) {
	PrintAlert(fmt.Sprint("Savage ", id, "- Pot is empty, I'll wake up the cook...", "\n"))
}

func getServingFromPot(id int) {
	PrintSavage(fmt.Sprint("Savage ", id, " is serving..."))
}

func main() {
	runtime.GOMAXPROCS(THREAD_COUNT)
	fmt.Println(RESET)
	fmt.Println(BOLD, "Dining Savages")
	fmt.Println(MAGENTA, "\tServings Count:\t", SERVINGS_COUNT)
	fmt.Println(GREEN, "\t", "Cooks: ", "\t", COOKS_COUNT)
	fmt.Println(CYAN, "\t", "Savages: ", "\t", SAVAGES_COUNT)
	fmt.Println(RESET)

	for i := 0; i < COOKS_COUNT; i++ {
		go Cook(i)
	}

	for i := 0; i < SAVAGES_COUNT; i++ {
		go Savage(i + COOKS_COUNT)
	}

	// Stops the program after some delay, so the user can see the output
	<- time.After(250 * time.Millisecond)
}

func PrintCook(s string) {
	fmt.Println(GREEN + s + RESET)
}

func PrintSavage(s string) {
	fmt.Println(CYAN + s + RESET)
}

func PrintAlert(s string) {
	fmt.Println(RED + s + RESET)
}

type Mutex struct {
	number [THREAD_COUNT]int
}

func NewMutex() *Mutex {
	return &Mutex{
		number: [THREAD_COUNT]int{},
	}
}

func (m *Mutex) maxNumber() int {
	if len(m.number) == 0 {
		return -1
	}

	max := m.number[0]
	for i := 1; i < len(m.number); i++ {
		if m.number[i] > max {
			max = m.number[i]
		}
	}
	return max
}

func (m *Mutex) Lock(i int) {
	m.number[i] = 1 + m.maxNumber()
	for j := 0; j < THREAD_COUNT; j++ {
		for m.number[j] != 0 && (m.number[j] < m.number[i] || (m.number[j] == m.number[i] && j < i)) {
			// Wait until all threads with smaller numbers or with the same
			// number, but with higher priority, finish their work
		}
	}
}

func (m *Mutex) Unlock(i int) {
	m.number[i] = 0
}

type Semaphore struct {
	v    int           // valor do semaforo: negativo significa proc bloqueado
	fila chan struct{} // canal para bloquear os processos se v < 0
	sc   chan struct{} // canal para atomicidade das operacoes wait e signal
}

func NewSemaphore(init int) *Semaphore {
	s := &Semaphore{
		v:    init,                   // valor inicial de creditos
		fila: make(chan struct{}),    // canal sincrono para bloquear processos
		sc:   make(chan struct{}, 1), // usaremos este como semaforo para SC, somente 0 ou 1
	}
	return s
}

func (s *Semaphore) Wait() {
	s.sc <- struct{}{} // SC do semaforo feita com canal
	s.v--              // decrementa valor
	if s.v < 0 {       // se negativo era 0 ou menor, tem que bloquear
		<-s.sc               // antes de bloq, libera acesso
		s.fila <- struct{}{} // bloqueia proc
	} else {
		<-s.sc // libera acesso
	}
}

func (s *Semaphore) Signal() {
	s.sc <- struct{}{} // entra sc
	s.v++
	if s.v <= 0 { // tem processo bloqueado ?
		<-s.fila // desbloqueia
	}
	<-s.sc // libera SC para outra op
}
