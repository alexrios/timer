package timer

import (
	"fmt"
	"time"
)

// Stopwatch is a simple stopwatch.
type Stopwatch struct {
	start              time.Time
	end                time.Time
	running            bool
	paused             bool
	pauseStartTime     time.Time
	totalPauseDuration time.Duration
	laps               []time.Duration
}

// Start starts the stopwatch.
func (s *Stopwatch) Start() {
	s.start = time.Now()
	if s.paused {
		s.totalPauseDuration += time.Since(s.pauseStartTime)
		s.paused = false
	} else {
		s.start = time.Now()
	}
	s.running = true
}

// Stop stops the stopwatch.
func (s *Stopwatch) Stop() {
	if s.paused {
		s.totalPauseDuration += time.Since(s.pauseStartTime)
		s.paused = false
	}
	s.end = time.Now()
	s.running = false
}

// Pause pauses the stopwatch.
func (s *Stopwatch) Pause() {
	if s.running && !s.paused {
		s.pauseStartTime = time.Now()
		s.paused = true
	}
}

// Resume resumes the stopwatch.
func (s *Stopwatch) Resume() {
	if s.running && s.paused {
		s.totalPauseDuration += time.Since(s.pauseStartTime)
		s.paused = false
	}
}

// Lap records a lap time.
func (s *Stopwatch) Lap() {
	if s.running {
		lapTime := time.Since(s.start) - s.totalPauseDuration
		s.laps = append(s.laps, lapTime)
	}
}

// Laps returns the recorded lap times.
func (s *Stopwatch) Laps() []time.Duration {
	return s.laps
}

// Elapsed returns the elapsed time.
func (s *Stopwatch) Elapsed() time.Duration {
	if s.running {
		return time.Since(s.start) - s.totalPauseDuration
	}
	return s.end.Sub(s.start) - s.totalPauseDuration
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
