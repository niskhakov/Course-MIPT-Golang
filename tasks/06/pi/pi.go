package pi

import (
	"sync"
)

const (
	A = 100.0 // Сторона квадрата
)

func CalculatePi(concurrent, iterations int, gen RandomPointGenerator) float64 {
	circNumChan := make(chan int)
	var wg sync.WaitGroup

	wg.Add(concurrent)

	for i := 0; i < concurrent; i++ {
		go calculatePiWorker(&wg, getIterations(iterations, concurrent, i), gen, circNumChan)
	}

	nCirc := 0
	for i := 0; i < concurrent; i++ {
		nCirc += <-circNumChan
	}

	wg.Wait()

	piNum := float64(nCirc) / float64(iterations) * 4

	// fmt.Printf("Calculating PI: numbers in circle = %d, total = %d, pi = %f\n", nCirc, iterations, piNum)
	return piNum
}

type RandomPointGenerator interface {
	Next() (float64, float64)
}

func calculatePiWorker(wg *sync.WaitGroup, iterations int, gen RandomPointGenerator, circNumChan chan<- int) {
	defer wg.Done()

	a := A
	nCirc := 0 // Число точек, попавших в круг
	for i := 0; i < iterations; i++ {
		x, y := gen.Next()
		x = x * a
		y = y * a

		if isInCircle(x, y, a) {
			nCirc++
		}
	}

	circNumChan <- nCirc
}

func isInCircle(x, y, r float64) bool {
	return x*x+y*y <= r*r
}

// getIterations returns fair amount of iterations to execute for specific workerId, workerId starts with 0
// Example: we have 20 tasks and 3 workers -> (worker - number of tasks)
// 			w1 - 7 tasks, w2 - 7 tasks, w3 - 6 tasks
func getIterations(tasks, workers, workerId int) int {
	delta := tasks%workers - workerId
	if delta > 0 {
		return tasks/workers + 1
	}
	return tasks / workers
}
