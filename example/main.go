package main

import (
	"fmt"
	"time"

	"github.com/alexrios/timer"
)

func main() {
	// Stopwatch example
	sw := &timer.Stopwatch{}
	sw.Start()
	time.Sleep(2 * time.Second)
	sw.Pause()
	fmt.Printf("Paused time: %s\n", timer.FormatDuration(sw.Elapsed()))
	time.Sleep(1 * time.Second)
	sw.Resume()
	time.Sleep(1 * time.Second)
	sw.Lap()
	sw.Stop()
	fmt.Printf("Total elapsed time: %s\n", timer.FormatDuration(sw.Elapsed()))

	// Display laps
	for i, lap := range sw.Laps() {
		fmt.Printf("Lap %d: %s\n", i+1, timer.FormatDuration(lap))
	}

	// Countdown example
	countdownDone := make(chan bool)
	countdown := timer.NewCountdown(3*time.Second, func() {
		fmt.Println("Countdown finished!")
		countdownDone <- true
	})
	countdown.Start()

	// Display countdown progress
	go func() {
		for progress := range countdown.Progress() {
			fmt.Printf("Countdown progress: %.2f%%\n", progress*100)
		}
	}()

	// Wait for countdown to finish or stop it early
	select {
	case <-countdownDone:
		// Countdown completed
	case <-time.After(1 * time.Second):
		countdown.Stop()
		fmt.Println("Countdown stopped early.")
	}
}
