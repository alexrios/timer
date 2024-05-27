package timer

import (
	"testing"
	"time"
)

func TestStopwatch(t *testing.T) {
	s := &Stopwatch{}
	s.Start()
	time.Sleep(1 * time.Second)
	s.Stop()

	elapsed := s.Elapsed()
	if elapsed < 1*time.Second || elapsed > 2*time.Second {
		t.Errorf("Expected elapsed time around 1 second, got %v", elapsed)
	}
}

func TestCountdown(t *testing.T) {
	done := make(chan bool)
	c := NewCountdown(1*time.Second, func() {
		done <- true
	})
	c.Start()

	select {
	case <-done:
		// Success
	case <-time.After(2 * time.Second):
		t.Errorf("Countdown did not trigger in expected time")
	}
}

func TestFormatDuration(t *testing.T) {
	d := 90061 * time.Second
	expected := "1d 01:01:01"
	result := FormatDuration(d)

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
