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
	sw.Stop()
	fmt.Printf("Elapsed time: %s\n", timer.FormatDuration(sw.Elapsed()))

	// Countdown example
	countdownDone := make(chan bool)
	countdown := timer.NewCountdown(3*time.Second, func() {
		fmt.Println("Countdown finished!")
		countdownDone <- true
	})
	countdown.Start()

	// Wait for countdown to finish or stop it early
	select {
	case <-countdownDone:
		// Countdown completed
	case <-time.After(1 * time.Second):
		countdown.Stop()
		fmt.Println("Countdown stopped early.")
	}
}
