package main

import (
	"fmt"
	"time"
)

const (
	/* Number of servings in the pot */
	M = 5
)

var (
	servings = M
	mutex = NewSemaphore(1)
	emptyPot = NewSemaphore(0)
	fullPot = NewSemaphore(0)
)

func putServingsInPot() {
	fmt.Println("\033[34m", "Putting servings in pot...", "\033[0m")
}

func Cook() {
	for {
		emptyPot.Wait()
		putServingsInPot()
		fullPot.Signal()
	}
}

func getServingFromPot() {
	fmt.Println("\033[33m", "Serving...", "\033[0m")
}

func eat() {
	fmt.Println("\033[32m", "Eating...", "\033[0m")
}

func Savage() {
	for {
		mutex.Wait()
		if servings == 0 {
			fmt.Println("\033[31m", "Pot is empty, waking up the cook...", "\033[0m")
			emptyPot.Signal()
			fullPot.Wait()
			servings = M
		}
		servings--
		getServingFromPot()
		mutex.Signal()
		eat()
	}
}

func main() {
	fmt.Println("Dining Savages")
	for i := 0; i < 1; i++ {
		go Cook()
	}
	for i := 0; i < 5; i++ {
		go Savage()
	}
	<- time.After(5 * time.Millisecond)
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
