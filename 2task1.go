package main

import (
	"fmt"
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
	numGoroutines := 3 // Количество горутин
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
			for _, domino := range fullSet[startIndex:endIndex] {
				found := false
				for _, existing := range existingSet {
					if existing.Left == domino.Left && existing.Right == domino.Right {
						found = true
						break
					}
				}
				if !found {
					missingDominosChan <- domino
				}
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

func main() {
	// Существующий набор костей домино
	existingSet := []Domino{
		{0, 0}, {0, 1}, {0, 2}, {0, 3},
		{1, 1}, {1, 2}, {1, 3},
		{2, 2}, {2, 3},
		{3, 3},
	}

	// Максимальное значение для костей домино
	maxValue := 5

	// Обработка без многозадачности
	processWithoutConcurrency(existingSet, maxValue)

	// Обработка с многозадачностью
	processWithConcurrency(existingSet, maxValue)
}