package timer

import (
	"context"
	"fmt"
	"time"
)

// Stopwatch is a simple stopwatch with context support.
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
func (s *Stopwatch) Start(ctx context.Context) error {
	if s.running {
		return fmt.Errorf("stopwatch is already running")
	}
	s.start = time.Now()
	s.running = true
	go func() {
		<-ctx.Done()
		s.Stop()
	}()
	return nil
}

// Stop stops the stopwatch.
func (s *Stopwatch) Stop() {
	if !s.running {
		return
	}
	if s.paused {
		s.totalPauseDuration += time.Since(s.pauseStartTime)
		s.paused = false
	}
	s.end = time.Now()
	s.running = false
}

// Pause pauses the stopwatch.
func (s *Stopwatch) Pause() error {
	if !s.running {
		return fmt.Errorf("stopwatch is not running")
	}
	if s.paused {
		return fmt.Errorf("stopwatch is already paused")
	}
	s.pauseStartTime = time.Now()
	s.paused = true
	return nil
}

// Resume resumes the stopwatch.
func (s *Stopwatch) Resume() error {
	if !s.running {
		return fmt.Errorf("stopwatch is not running")
	}
	if !s.paused {
		return fmt.Errorf("stopwatch is not paused")
	}
	s.totalPauseDuration += time.Since(s.pauseStartTime)
	s.paused = false
	return nil
}

// Lap records a lap time.
func (s *Stopwatch) Lap() error {
	if !s.running {
		return fmt.Errorf("stopwatch is not running")
	}
	lapTime := time.Since(s.start) - s.totalPauseDuration
	s.laps = append(s.laps, lapTime)
	return nil
}

// Laps returns the recorded lap times.
func (s *Stopwatch) Laps() []time.Duration {
	return s.laps
}

// Elapsed returns the elapsed time.
func (s *Stopwatch) Elapsed() (time.Duration, error) {
	if s.running {
		return time.Since(s.start) - s.totalPauseDuration, nil
	}
	if s.start.IsZero() {
		return 0, fmt.Errorf("stopwatch has not been started")
	}
	return s.end.Sub(s.start) - s.totalPauseDuration, nil
}
