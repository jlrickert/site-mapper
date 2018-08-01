package main

func add1(n int) int {
	return n + 1
}

func sub1(n int) int {
	return n - 1
}

type empty struct{}
type Semaphore chan empty

func (sem Semaphore) Acquire(n int) {
	e := empty{}
	for i := 0; i < n; i++ {
		sem <- e
	}
}

func (sem Semaphore) Release(n int) {
	for i := 0; i < n; i++ {
		<-sem
	}
}

func (sem Semaphore) Lock() {
	sem.Acquire(1)
}

func (sem Semaphore) Unlock() {
	sem.Release(1)
}

func (sem Semaphore) Signal() {
	sem.Release(1)
}

func (sem Semaphore) Wait(n int) {
	sem.Acquire(n)
}
