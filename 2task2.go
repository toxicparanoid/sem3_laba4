package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Структура для кости домино
type Domino struct {
	Left  int
	Right int
}

// Функция для генерации полного набора костей домино
func generateFullSet(maxValue int) []Domino {
	var fullSet []Domino
	for i := 0; i <= maxValue; i++ {
		for j := i; j <= maxValue; j++ {
			fullSet = append(fullSet, Domino{Left: i, Right: j})
		}
	}
	return fullSet
}

// Функция для поиска недостающих костей домино
func findMissingDominos(existingSet []Domino, fullSet []Domino) []Domino {
	existingMap := make(map[Domino]bool)
	for _, domino := range existingSet {
		existingMap[domino] = true
	}
	var missingDominos []Domino
	for _, domino := range fullSet {
		time.Sleep(1 * time.Millisecond)
		if !existingMap[domino] {
			missingDominos = append(missingDominos, domino)
		}
	}
	return missingDominos
}

// Обработка без многозадачности
func processWithoutConcurrency(existingSet []Domino, maxValue int) {
	start := time.Now()

	fullSet := generateFullSet(maxValue)
	missingDominos := findMissingDominos(existingSet, fullSet)

	duration := time.Since(start)

	fmt.Printf("Без многозадачности:\n")
	fmt.Printf("Недостающие кости домино: %v\n", missingDominos)
	fmt.Printf("Время обработки: %v\n\n", duration)
}

// Обработка с многозадачностью
func processWithConcurrency(existingSet []Domino, maxValue int) {
	start := time.Now()

	var wg sync.WaitGroup

	// Генерируем полный набор костей домино
	fullSet := generateFullSet(maxValue)

	// Канал для хранения результатов
	missingDominosChan := make(chan Domino, len(fullSet))

	// Разбиваем данные на части для многозадачности
	numGoroutines := 4 // Увеличиваем количество горутин
	chunkSize := len(fullSet) / numGoroutines

	wg.Add(numGoroutines)
	// Горутины для проверки наличия костей
	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			defer wg.Done()
			startIndex := i * chunkSize
			endIndex := (i + 1) * chunkSize
			if i == numGoroutines-1 { // Для последнего сегмента
				endIndex = len(fullSet)
			}
			// Локальный набор результатов
			localMissingDominos := findMissingDominos(existingSet, fullSet[startIndex:endIndex])
			for _, domino := range localMissingDominos {
				missingDominosChan <- domino
			}
		}(i)
	}

	// Ждем завершения всех горутин
	wg.Wait()
	close(missingDominosChan)

	// Собираем результаты
	var missingDominos []Domino
	for domino := range missingDominosChan {
		missingDominos = append(missingDominos, domino)
	}

	duration := time.Since(start)

	fmt.Printf("С многозадачностью (с несколькими горутинами):\n")
	fmt.Printf("Недостающие кости домино: %v\n", missingDominos)
	fmt.Printf("Время обработки: %v\n\n", duration)
}

// Функция для генерации случайной кости домино
func generateRandomDomino(maxValue int) Domino {
	left := rand.Intn(maxValue + 1)
	right := rand.Intn(maxValue + 1)
	return Domino{Left: left, Right: right}
}

func main() {
	// Инициализируем генератор случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Максимальное значение для костей домино
	maxValue := 10 // Увеличиваем максимальное значение для генерации большего набора

	// Генерируем случайный набор костей домино
	var existingSet []Domino
	for i := 0; i < 100; i++ { // Генерируем 100 случайных костей
		existingSet = append(existingSet, generateRandomDomino(maxValue))
	}

	// Обработка без многозадачности
	processWithoutConcurrency(existingSet, maxValue)

	// Обработка с многозадачностью
	processWithConcurrency(existingSet, maxValue)
}