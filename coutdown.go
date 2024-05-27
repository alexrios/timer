package timer

import (
	"context"
	"fmt"
	"time"
)

// Countdown is a simple countdown timer with context support.
type Countdown struct {
	duration     time.Duration
	callback     func()
	stopChan     chan struct{}
	progressChan chan float64
}

// NewCountdown creates a new countdown timer.
func NewCountdown(d time.Duration, cb func()) *Countdown {
	return &Countdown{
		duration:     d,
		callback:     cb,
		stopChan:     make(chan struct{}),
		progressChan: make(chan float64),
	}
}

// Start starts the countdown timer.
// Start starts the countdown timer.
func (c *Countdown) Start(ctx context.Context) error {
	if c.duration <= 0 {
		return fmt.Errorf("duration must be greater than zero")
	}

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		startTime := time.Now()

		for {
			select {
			case <-ticker.C:
				elapsed := time.Since(startTime)
				progress := float64(elapsed) / float64(c.duration)
				fmt.Println("Progress:", progress)

				if progress >= 1 { // Check if countdown is complete
					fmt.Println("Countdown complete")
					c.callback()
					close(c.progressChan)
					return // Exit the loop
				}

				select {
				case c.progressChan <- progress:
					// Progress sent successfully
				default:
					// Handle the case where the channel is full (if it's unbuffered)
					fmt.Println("Progress channel full, dropping update")
				}

			case <-ctx.Done():
				// Timer stopped by context cancellation
				close(c.stopChan)
				return
			case <-c.stopChan:
				// Timer stopped manually
				return
			}
		}
	}()
	return nil
}

// Stop stops the countdown timer.
func (c *Countdown) Stop() {
	close(c.stopChan)
}

// Progress returns a channel for receiving progress updates.
func (c *Countdown) Progress() <-chan float64 {
	return c.progressChan
}
