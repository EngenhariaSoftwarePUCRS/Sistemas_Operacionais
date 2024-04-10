package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

const (
	N          = 5 // Thread Count
	QUEUE_SIZE = 50
)

var (
	criticalSection  *Mutex = NewMutex()
	items  *Mutex = NewMutex()
	emptyQueue *Mutex = NewMutex()
	queue  Queue  = make(Queue, 0, QUEUE_SIZE)
)

func Producer(id int) {
	for {
		fmt.Println("\033[32m", "Pronto para produzir...", "\033[0m")
		item := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(99)
		emptyQueue.Lock(id)
		criticalSection.Lock(id)
		queue.Enqueue(item)
		fmt.Println("\033[32m",
			id, "- Adicionando: ", item, "\n",
			id, "- Fila Critical Section: ", queue, "\n",
			id, "- TimeStamp: ", CurrentTimeStamp(), "\n",
		"\033[0m")
		criticalSection.UnlockOthers(id)
		items.UnlockOthers(id)
	}
}

// WORKAROUND
// This is a way of showing the user the aproximate time the operation was done contrary to the time the operation was printed
// This is a workaround because, sometimes, a thread has already inserted/removed an item from the queue but the OS gives the CPU to another thread to remove another item, so printing goes out of order
func CurrentTimeStamp() string {
	milliseconds := time.Now().UnixMilli()
	milliseconds = milliseconds % 10_000_000
	return fmt.Sprint(milliseconds)
}

func Consumer(idProducer int, idConsumer int) {
	for {
		items.Lock(idConsumer)
		fmt.Println("\033[34m", idConsumer, "- Tentando consumir...", "\033[0m")
		criticalSection.Lock(idConsumer)
		previousQueue := queue
		item, success := queue.Dequeue()
		if success {
			fmt.Println("\033[31m",
				idConsumer, "- Anterior: ", previousQueue, "\n",
				idConsumer, "- Removido: ", item, "\n",
				idConsumer, "- New:", queue, "\n",
				idConsumer, "- TimeStamp: ", CurrentTimeStamp(), "\n",
			"\033[0m")
		}
		criticalSection.UnlockOthers(idConsumer)
		emptyQueue.Unlock(idProducer)
	}
}

func main() {
	runtime.GOMAXPROCS(N)
	fmt.Printf("Producer Consumer with Queue Size %d\n", QUEUE_SIZE)
	producerId := 0
	go Producer(producerId)
	for i := 1; i < N; i++ {
		go Consumer(producerId, i)
	}
	<-time.After(300 * time.Millisecond)
}

type Queue []int

func (q *Queue) Enqueue(item int) {
	if len(*q) == QUEUE_SIZE {
		return
	}
	*q = append(*q, item)
}

func (q *Queue) Dequeue() (head int, success bool) {
	if len(*q) == 0 {
		return -1, false
	}
	item := (*q)[0]
	*q = (*q)[1:]
	return item, true
}

type Mutex struct {
	// volatile bool entering[0..N] := 0
	entering [N]bool
	// volatile int number[0..N] := 0
	number [N]int
}

func NewMutex() *Mutex {
	return &Mutex{
		entering: [N]bool{},
		number:   [N]int{},
	}
}

func (m *Mutex) printEntering() {
	fmt.Print("Entering: ")
	fmt.Print("[")
	for i := 0; i < len(m.entering); i++ {
		fmt.Printf("%d ", m.entering[i])
	}
	fmt.Println("\b]")
}

func (m *Mutex) PrintNumber() {
	fmt.Print("Number: ")
	fmt.Print("[")
	for i := 0; i < len(m.number); i++ {
		fmt.Printf("%d ", m.number[i])
	}
	fmt.Println("\b]")
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
	m.entering[i] = true
	m.number[i] = 1 + m.maxNumber()
	m.entering[i] = false
	for j := 0; j < N; j++ {
		for m.entering[j] {
			// Wait until thread j receives its number:
		}
		for m.number[j] != 0 && (m.number[j] < m.number[i] || (m.number[j] == m.number[i] && j < i)) {
			// Wait until all threads with smaller numbers or with the same
			// number, but with higher priority, finish their work
		}
	}
}

func (m *Mutex) Unlock(i int) {
	m.number[i] = 0
}

func (m *Mutex) UnlockOthers(i int) {
	for j := 0; j < N; j++ {
		if j != i {
			m.number[j] = 0
		}
	}
}
