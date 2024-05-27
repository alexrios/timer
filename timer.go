package timer

import (
	"fmt"
	"time"
)

// Stopwatch is a simple stopwatch.
type Stopwatch struct {
	start   time.Time
	end     time.Time
	running bool
}

// Start starts the stopwatch.
func (s *Stopwatch) Start() {
	s.start = time.Now()
	s.running = true
}

// Stop stops the stopwatch.
func (s *Stopwatch) Stop() {
	s.end = time.Now()
	s.running = false
}

// Elapsed returns the elapsed time.
func (s *Stopwatch) Elapsed() time.Duration {
	if s.running {
		return time.Since(s.start)
	}
	return s.end.Sub(s.start)
}

// Countdown is a simple countdown timer.
type Countdown struct {
	duration time.Duration
	callback func()
	stopChan chan struct{}
}

// NewCountdown creates a new countdown timer.
func NewCountdown(d time.Duration, cb func()) *Countdown {
	return &Countdown{
		duration: d,
		callback: cb,
		stopChan: make(chan struct{}),
	}
}

// Start starts the countdown timer.
func (c *Countdown) Start() {
	go func() {
		select {
		case <-time.After(c.duration):
			c.callback()
		case <-c.stopChan:
			// Timer stopped
		}
	}()
}

// Stop stops the countdown timer.
func (c *Countdown) Stop() {
	close(c.stopChan)
}

// FormatDuration formats a duration in a human-readable way.
func FormatDuration(d time.Duration) string {
	days := d / (24 * time.Hour)
	d -= days * 24 * time.Hour
	hours := d / time.Hour
	d -= hours * time.Hour
	minutes := d / time.Minute
	d -= minutes * time.Minute
	seconds := d / time.Second

	if days > 0 {
		return fmt.Sprintf("%dd %02d:%02d:%02d", days, hours, minutes, seconds)
	}
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
