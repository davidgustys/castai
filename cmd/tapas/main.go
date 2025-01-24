package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var menu = map[string]int{
	"chopitos":            10,
	"chorizo":             5,
	"croquetas":           5,
	"patatas bravas":      7,
	"pimientos de padron": 5,
}

type Tapa struct {
	Name string
}

var visitors = []string{
	"Alice",
	"Bob",
	"Charlie",
	"Dave",
}

// getTapa returns a random tapa from the menu with remaining quantity
func getRandomTapa(mu *sync.Mutex) (string, error) {
	mu.Lock()
	defer mu.Unlock()

	// Count total remaining tapas for weighted random selection
	total := 0
	for _, quantity := range menu {
		total += quantity
	}

	if total == 0 {
		return "", fmt.Errorf("no more tapas")
	}

	// Random selection weighted by quantity
	randomIndex := rand.Intn(total)
	currentIndex := 0
	for tapa, quantity := range menu {
		currentIndex += quantity
		if randomIndex < currentIndex {
			menu[tapa]--
			return tapa, nil
		}
	}

	return "", fmt.Errorf("unexpected error selecting tapa")
}

func main() {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	fmt.Println("Tapas bar opened!")

	var mu sync.Mutex
	var wg sync.WaitGroup

	// Channel for prepared tapas, buffered to total capacity
	totalTapas := 0
	for _, quantity := range menu {
		totalTapas += quantity
	}
	preparedTapas := make(chan Tapa, totalTapas)

	// Chef goroutine
	go func() {
		for {
			tapa, err := getRandomTapa(&mu)
			if err != nil {
				fmt.Println("Chef: all tapas served!")
				close(preparedTapas)
				return
			}

			fmt.Printf("Chef: preparing %s...\n", tapa)
			time.Sleep(time.Second) // Fixed preparation time
			fmt.Printf("Chef: %s served!\n", tapa)
			preparedTapas <- Tapa{Name: tapa}
		}
	}()

	// Visitor goroutines
	for _, visitor := range visitors {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			for {
				tapa, ok := <-preparedTapas
				if !ok {
					fmt.Printf("%s: no more tapas to eat, leaving...\n", name)
					return
				}

				// Random eating duration between 0 and 2 seconds
				duration := time.Duration(rand.Float64() * 2 * float64(time.Second))
				fmt.Printf("%s: enjoying %s for %.1fs...\n", name, tapa.Name, duration.Seconds())
				time.Sleep(duration)
				fmt.Printf("%s: finished eating %s\n", name, tapa.Name)
			}
		}(visitor)
	}

	// Wait for all visitors to finish
	wg.Wait()
	fmt.Println("Tapas bar closed!")
}
