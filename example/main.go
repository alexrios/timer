package main

import (
	"context"
	"fmt"
	"time"

	"github.com/alexrios/timer/v2"
)

func main() {
	// Create a new stopwatch
	sw := &timer.Stopwatch{}
	ctx, cancel := context.WithCancel(context.Background())

	// Start the stopwatch
	if err := sw.Start(ctx); err != nil {
		fmt.Println("Error starting stopwatch:", err)
		return
	}

	// Simulate some work
	time.Sleep(2 * time.Second)

	// Record a lap
	if err := sw.Lap(); err != nil {
		fmt.Println("Error recording lap:", err)
		return
	}

	// Stop the stopwatch after some time
	time.Sleep(1 * time.Second)
	cancel() // This will stop the stopwatch

	// Get elapsed time
	elapsed, err := sw.Elapsed()
	if err != nil {
		fmt.Println("Error getting elapsed time:", err)
		return
	}

	fmt.Printf("Elapsed time: %v\n", elapsed)
}
