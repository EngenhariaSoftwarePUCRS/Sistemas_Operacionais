package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	N          = 5 // Thread Count
	QUEUE_SIZE = 50
)

var (
	mutex  *Semaphore = NewSemaphore(1)
	items  *Semaphore = NewSemaphore(0)
	spaces *Semaphore = NewSemaphore(QUEUE_SIZE)
	queue  Queue      = make(Queue, 0, QUEUE_SIZE)

	// volatile bool entering[0..N] := 0
	entering [N]bool
	// volatile int number[0..N] := 0
	number [N]int
)

func Producer() {
	for {
		fmt.Println("\033[32m", "Produzindo...", "\033[0m")
		item := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(99)
		spaces.Wait()
		mutex.Wait()
		fmt.Println("\033[32m", "Fila Critical Section: ", queue, "\033[0m")
		queue.Enqueue(item)
		fmt.Println("\033[32m", "Adicionado: ", item, "\033[0m")
		mutex.Signal()
		items.Signal()
	}
}

func Consumer() {
	for {
		fmt.Println("\033[31m", "Consumindo...", "\033[0m")
		items.Wait()
		mutex.Wait()
		fmt.Println("\033[31m", "Fila Critical Section: ", queue, "\033[0m")
		item := queue.Dequeue()
		fmt.Println("\033[31m", "Removido: ", item, "\033[0m")
		mutex.Signal()
		spaces.Signal()
	}
}

func main() {
	fmt.Printf("Producer Consumer with Queue Size %d\n", QUEUE_SIZE)
	// for i := 0; i < 1; i++ {
	// 	go Producer()
	// }
	// for i := 0; i < 5; i++ {
	// 	go Consumer()
	// }
	bay(number)
	go func() { lock(1) }()
	go func() { lock(3) }()
	go func() { lock(2) }()
	go func() { unlock(1) }()
	go func() { unlock(3) }()
	bay(number)
	<-time.After(5 * time.Millisecond)
}

type Queue []int

func (q *Queue) Enqueue(item int) {
	if len(*q) == QUEUE_SIZE {
		return
	}
	*q = append(*q, item)
}

func (q *Queue) Dequeue() int {
	if len(*q) == 0 {
		return -1
	}
	item := (*q)[0]
	*q = (*q)[1:]
	return item
}

func bay(slice [N]int) {
	fmt.Print("Numbers: ")
	fmt.Print("[")
	for i := 0; i < len(number); i++ {
		fmt.Printf("%d ", number[i])
	}
	fmt.Println("\b]")
}

func max(slice [N]int) int {
	if len(slice) == 0 {
		return -1
	}

	max := slice[0]
	for i := 1; i < len(slice); i++ {
		if slice[i] > max {
			max = slice[i]
		}
	}
	return max
}

func lock(i int) {
	entering[i] = true
	number[i] = 1 + max(number)
	entering[i] = false
	for j := 0; j < N; j++ {
		for entering[j] {
			// Wait until thread j receives its number:
		}
		for number[j] != 0 && (number[j] < number[i] || (number[j] == number[i] && j < i)) {
			// Wait until all threads with smaller numbers or with the same
			// number, but with higher priority, finish their work
		}
	}
	fmt.Printf("lock(%d) - %d\n", i, max(number))
}

func unlock(i int) {
	number[i] = 0
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
