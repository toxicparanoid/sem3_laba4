package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// Количество потоков
const numGoroutines = 10

// Генерация случайного ASCII символа
func generateRandomASCII() byte {
	return byte(rand.Intn(94) + 33)
}

// StopWatch обертка для измерения времени выполнения функции
func StopWatch(name string, f func()) {
	start := time.Now()
	f()
	duration := time.Since(start)
	fmt.Printf("%s Time: %v\n", name, duration)
}

// Mutex
func testMutex(wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()
	mu.Lock()
	fmt.Printf("Mutex: %c\n", generateRandomASCII())
	mu.Unlock()
}

// Semaphore (канал с буфером)
func testSemaphore(wg *sync.WaitGroup, sem chan struct{}) {
	defer wg.Done()
	sem <- struct{}{}
	fmt.Printf("Semaphore: %c\n", generateRandomASCII())
	<-sem
}

// SemaphoreSlim (ограничение повторных попыток)
func testSemaphoreSlim(wg *sync.WaitGroup, sem chan struct{}, retries int) {
	defer wg.Done()
	for i := 0; i < retries; i++ {
		select {
		case sem <- struct{}{}:
			fmt.Printf("SemaphoreSlim: %c\n", generateRandomASCII())
			<-sem
			return
		default:
			time.Sleep(time.Millisecond * 10) // ожидание перед новой попыткой
		}
	}
}

// Barrier (WaitGroup)
func testBarrier(wg *sync.WaitGroup, barrier *sync.WaitGroup) {
	defer wg.Done()
	barrier.Done()
	barrier.Wait()
	fmt.Printf("Barrier: %c\n", generateRandomASCII())
}

// SpinLock с использованием атомарной операции
func testSpinLock(counter *int32) {
	for {
		if atomic.CompareAndSwapInt32(counter, 0, 1) {
			fmt.Printf("SpinLock: %c\n", generateRandomASCII())
			atomic.StoreInt32(counter, 0)
			break
		}
	}
}

// SpinWait (активное ожидание)
func testSpinWait() {
	spinCount := 0
	for spinCount < 1000 { // активное ожидание с контролем
		if spinCount % 100 == 0 { // добавляем паузы через каждые 100 итераций
			randomDelay := time.Duration(rand.Intn(10)+1) * time.Microsecond
			time.Sleep(randomDelay)
		}
		spinCount++
	}
	fmt.Printf("SpinWait: %c\n", generateRandomASCII())
}

// Monitor (Mutex + Cond)
func testMonitor(wg *sync.WaitGroup, mu *sync.Mutex, cond *sync.Cond) {
	defer wg.Done()
	mu.Lock()
	cond.Wait()
	fmt.Printf("Monitor: %c\n", generateRandomASCII())
	mu.Unlock()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	var wg sync.WaitGroup

	// Тест Mutex
	mu := &sync.Mutex{}
	StopWatch("Mutex", func() {
		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go testMutex(&wg, mu)
		}
		wg.Wait()
	})

	// Тест Semaphore
	sem := make(chan struct{}, 3) // Ограничение на 3 горутины
	StopWatch("Semaphore", func() {
		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go testSemaphore(&wg, sem)
		}
		wg.Wait()
	})

	// Тест SemaphoreSlim
	retries := 5
	StopWatch("SemaphoreSlim", func() {
		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go testSemaphoreSlim(&wg, sem, retries)
		}
		wg.Wait()
	})

	// Тест Barrier
	barrier := &sync.WaitGroup{}
	barrier.Add(numGoroutines)
	StopWatch("Barrier", func() {
		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go testBarrier(&wg, barrier)
		}
		wg.Wait()
	})

	// Тест SpinLock
	var counter int32
	StopWatch("SpinLock", func() {
		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				testSpinLock(&counter)
			}()
		}
		wg.Wait()
	})

	// Тест SpinWait
	StopWatch("SpinWait", func() {
		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				testSpinWait()
			}()
		}
		wg.Wait()
	})

	// Тест Monitor
	mu = &sync.Mutex{}
	cond := sync.NewCond(mu)
	StopWatch("Monitor", func() {
		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go testMonitor(&wg, mu, cond)
		}
		time.Sleep(time.Microsecond * 1000) // даём время горутинам заблокироваться
		cond.Broadcast()
		wg.Wait()
	})
}
