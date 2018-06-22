package main

import (
	"database/sql"
	"log"
	"math"
	"runtime"
	"time"
)

var interval = time.Second

func main() {
	log.Println("Running CPU benchmarks...")
	variance := bench(20, func() {
		fib(38)
	})
	log.Println("Variance:", variance)
	log.Println("Running memory benchmarks...")
	variance = bench(20, func() {
		alloc(1024 * 1024 / 4)
	})
	log.Println("Variance:", variance)
}

func fib(n int) int {
	if n < 2 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

func alloc(n int) {
	c := []sql.ColumnType{}
	for i := 0; i < n; i++ {
		c = append(c, sql.ColumnType{})
	}
	c = nil
	runtime.GC()
}

func stdev(times []time.Duration) time.Duration {
	var total time.Duration
	for _, t := range times {
		total += t
	}
	mean := total / time.Duration(len(times))
	var diff float64
	for i := range times {
		diff += math.Pow(float64(times[i]-mean), 2)
	}
	return time.Duration(math.Sqrt(diff / float64(len(times))))
}

func bench(n int, fn func()) time.Duration {
	times := []time.Duration{}
	dc := make(chan time.Duration)

	for i := 0; i < runtime.NumCPU(); i++ {
		go func(dc chan time.Duration) {
			for j := 0; j < n; j++ {
				start := time.Now()
				fn()
				dc <- time.Now().Sub(start)
				time.Sleep(interval)
			}
		}(dc)
	}

	for i := 0; i < runtime.NumCPU()*n; i++ {
		d := <-dc
		times = append(times, d)
	}

	return stdev(times)
}
