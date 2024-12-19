package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	numPhilosophers = 5 // Количество философов и вилок (должно быть одинаково)
)

// Fork представляет вилку, которая используется для синхронизации доступа к ней.
type Fork struct {
	sync.Mutex
}

// Philosopher представляет философа, который ест и размышляет.
type Philosopher struct {
	id                int   // Идентификатор философа
	leftFork, rightFork *Fork // Вилки, которые использует философ
}

// dine — основная функция философа, которая запускает цикл еды и размышлений.
func (p Philosopher) dine(wg *sync.WaitGroup, done chan struct{}) {
	defer wg.Done() // Уменьшаем счётчик WaitGroup после завершения работы философа

	for {
		select {
		case <-done: // Если канал done закрыт, философ завершает работу
			fmt.Printf("Философ %d закончил обедать.\n", p.id)
			return
		default: // Иначе философ продолжает есть и размышлять
			p.think()
			p.eat()
		}
	}
}

// think — функция, которая имитирует процесс размышления философа.
func (p Philosopher) think() {
	fmt.Printf("Философ %d размышляет о великом.\n", p.id)
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) // Случайная пауза
}

// eat — функция, которая имитирует процесс еды философа.
func (p Philosopher) eat() {
	// Чтобы избежать взаимной блокировки (deadlock), философы с чётными ID берут левую вилку первой,
	// а философы с нечётными ID — правую вилку первой.
	if p.id%2 == 0 {
		p.leftFork.Lock()  // Философ берёт левую вилку
		p.rightFork.Lock() // Затем берёт правую вилку
	} else {
		p.rightFork.Lock() // Философ берёт правую вилку
		p.leftFork.Lock()  // Затем берёт левую вилку
	}

	// Философ ест
	fmt.Printf("Философ %d ест спагетти.\n", p.id)
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) // Случайная пауза

	// После еды философ кладёт вилки обратно
	p.leftFork.Unlock()
	p.rightFork.Unlock()
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Инициализация генератора случайных чисел

	// Создаём вилки (по одной на каждого философа)
	forks := make([]*Fork, numPhilosophers)
	for i := 0; i < numPhilosophers; i++ {
		forks[i] = &Fork{}
	}

	// Создаём философов и назначаем им вилки
	philosophers := make([]*Philosopher, numPhilosophers)
	for i := 0; i < numPhilosophers; i++ {
		philosophers[i] = &Philosopher{
			id:        i,                     // Идентификатор философа
			leftFork:  forks[i],              // Левую вилку берём по индексу i
			rightFork: forks[(i+1)%numPhilosophers], // Правую вилку берём по кругу
		}
	}

	var wg sync.WaitGroup // WaitGroup для ожидания завершения всех философов
	done := make(chan struct{}) // Канал для сигнала о завершении работы

	// Запускаем философов
	for _, philosopher := range philosophers {
		wg.Add(1) // Увеличиваем счётчик WaitGroup
		go philosopher.dine(&wg, done) // Запускаем философа в отдельной горутине
	}

	// Философы едят 5 секунд
	time.Sleep(5 * time.Second)
	close(done) // Закрываем канал done, чтобы сообщить философам о завершении

	wg.Wait() // Ожидаем завершения всех философов
	fmt.Println("Все философы закончили обедать.")
}