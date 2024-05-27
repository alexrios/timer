package timer

import (
	"context"
	"testing"
	"time"
)

func TestStopwatch(t *testing.T) {
	s := &Stopwatch{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := s.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start stopwatch: %v", err)
	}

	time.Sleep(1 * time.Second)
	s.Stop()

	elapsed, err := s.Elapsed()
	if err != nil {
		t.Fatalf("Failed to get elapsed time: %v", err)
	}
	if elapsed < 1*time.Second || elapsed > 2*time.Second {
		t.Errorf("Expected elapsed time around 1 second, got %v", elapsed)
	}
}

func TestCountdown(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		done := make(chan bool)
		c := NewCountdown(1*time.Second, func() { done <- true })

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		err := c.Start(ctx)
		if err != nil {
			t.Fatalf("Failed to start countdown: %v", err)
		}

		select {
		case <-done:
			// Success!
		case <-time.After(2 * time.Second):
			t.Fatal("Countdown did not trigger in expected time")
		}
	})

	t.Run("Cancellation", func(t *testing.T) {
		c := NewCountdown(5*time.Second, func() {
			t.Fatal("Callback should not be called")
		})

		ctx, cancel := context.WithCancel(context.Background())

		err := c.Start(ctx)
		if err != nil {
			t.Fatalf("Failed to start countdown: %v", err)
		}

		// Cancel after a short delay
		time.AfterFunc(1*time.Second, cancel)

		select {
		case <-c.stopChan:
			close(c.progressChan)
			// Ensure the progress channel is closed
			select {
			case _, ok := <-c.progressChan:
				if ok {
					t.Fatal("progressChan should be closed after cancellation")
				}
			default:
				// Channel is closed, as expected
			}
		case <-time.After(2 * time.Second):
			t.Fatal("Countdown did not stop after cancellation")
		}
	})
}

func TestFormatDuration(t *testing.T) {
	d := 90061 * time.Second
	expected := "1d 01:01:01"
	result := FormatDuration(d)

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestStopwatchPauseResume(t *testing.T) {
	s := &Stopwatch{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := s.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start stopwatch: %v", err)
	}

	time.Sleep(1 * time.Second)
	err = s.Pause()
	if err != nil {
		t.Fatalf("Failed to pause stopwatch: %v", err)
	}

	time.Sleep(1 * time.Second)
	err = s.Resume()
	if err != nil {
		t.Fatalf("Failed to resume stopwatch: %v", err)
	}

	time.Sleep(1 * time.Second)
	s.Stop()

	elapsed, err := s.Elapsed()
	if err != nil {
		t.Fatalf("Failed to get elapsed time: %v", err)
	}
	if elapsed < 2*time.Second || elapsed > 3*time.Second {
		t.Errorf("Expected elapsed time around 2 seconds, got %v", elapsed)
	}
}

func TestStopwatchLap(t *testing.T) {
	s := &Stopwatch{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := s.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start stopwatch: %v", err)
	}

	time.Sleep(1 * time.Second)
	err = s.Lap()
	if err != nil {
		t.Fatalf("Failed to record lap: %v", err)
	}

	time.Sleep(1 * time.Second)
	err = s.Lap()
	if err != nil {
		t.Fatalf("Failed to record lap: %v", err)
	}

	s.Stop()

	laps := s.Laps()
	if len(laps) != 2 {
		t.Errorf("Expected 2 laps, got %d", len(laps))
	}
	if laps[0] < 1*time.Second || laps[0] > 2*time.Second {
		t.Errorf("Expected first lap time around 1 second, got %v", laps[0])
	}
	if laps[1] < 2*time.Second || laps[1] > 3*time.Second {
		t.Errorf("Expected second lap time around 2 seconds, got %v", laps[1])
	}
}

func TestCountdownProgress(t *testing.T) {
	done := make(chan bool)
	c := NewCountdown(1*time.Second, func() {
		done <- true
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := c.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start countdown: %v", err)
	}

	progressUpdates := make([]float64, 0)
	go func() {
		for progress := range c.Progress() {
			progressUpdates = append(progressUpdates, progress)
		}
	}()

	select {
	case <-done:
		// Success
	case <-time.After(2 * time.Second):
		t.Errorf("Countdown did not trigger in expected time")
	}

	if len(progressUpdates) == 0 {
		t.Errorf("Expected progress updates, but got none")
	}
}

func TestStopwatchContextCancellation(t *testing.T) {
	s := &Stopwatch{}
	ctx, cancel := context.WithCancel(context.Background())

	err := s.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start stopwatch: %v", err)
	}

	time.Sleep(500 * time.Millisecond)
	cancel()

	time.Sleep(500 * time.Millisecond) // Give some time for the context to cancel

	elapsed, err := s.Elapsed()
	if err != nil {
		t.Fatalf("Failed to get elapsed time: %v", err)
	}
	if elapsed >= 1*time.Second {
		t.Errorf("Expected elapsed time less than 1 second due to context cancellation, got %v", elapsed)
	}
}

func TestStopwatchErrorHandling(t *testing.T) {
	s := &Stopwatch{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the stopwatch
	err := s.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start stopwatch: %v", err)
	}

	// Try starting the stopwatch again
	err = s.Start(ctx)
	if err == nil {
		t.Error("Expected error when starting an already running stopwatch, but got none")
	}

	// Pause the stopwatch
	err = s.Pause()
	if err != nil {
		t.Fatalf("Failed to pause stopwatch: %v", err)
	}

	// Try pausing the stopwatch again
	err = s.Pause()
	if err == nil {
		t.Error("Expected error when pausing an already paused stopwatch, but got none")
	}

	// Resume the stopwatch
	err = s.Resume()
	if err != nil {
		t.Fatalf("Failed to resume stopwatch: %v", err)
	}

	// Try resuming the stopwatch again
	err = s.Resume()
	if err == nil {
		t.Error("Expected error when resuming a stopwatch that is not paused, but got none")
	}
}

func TestStopwatchElapsedBeforeStart(t *testing.T) {
	s := &Stopwatch{}

	_, err := s.Elapsed()
	if err == nil {
		t.Error("Expected error when getting elapsed time of a stopwatch that has not been started, but got none")
	}
}

func TestMultipleCountdowns(t *testing.T) {
	done1 := make(chan bool)
	done2 := make(chan bool)

	c1 := NewCountdown(1*time.Second, func() {
		done1 <- true
	})
	c2 := NewCountdown(2*time.Second, func() {
		done2 <- true
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := c1.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start first countdown: %v", err)
	}

	err = c2.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start second countdown: %v", err)
	}

	select {
	case <-done1:
		// Success
	case <-time.After(2 * time.Second):
		t.Errorf("First countdown did not trigger in expected time")
	}

	select {
	case <-done2:
		// Success
	case <-time.After(3 * time.Second):
		t.Errorf("Second countdown did not trigger in expected time")
	}
}

func TestCountdownProgressUpdates(t *testing.T) {
	done := make(chan bool)
	c := NewCountdown(1*time.Second, func() {
		done <- true
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := c.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start countdown: %v", err)
	}

	progressUpdates := make([]float64, 0)
	go func() {
		for progress := range c.Progress() {
			progressUpdates = append(progressUpdates, progress)
		}
	}()

	select {
	case <-done:
		// Success
	case <-time.After(2 * time.Second):
		t.Errorf("Countdown did not trigger in expected time")
	}

	if len(progressUpdates) == 0 {
		t.Errorf("Expected progress updates, but got none")
	}

	// Check that progress updates are in the correct range
	for _, progress := range progressUpdates {
		if progress < 0 || progress > 1 {
			t.Errorf("Expected progress between 0 and 1, but got %v", progress)
		}
	}
}
