package go_smartchan

import (
	"sync"
	"testing"
)

func TestDataFlow(t *testing.T) {
	sc := NewSmartChan(5)

	// Test writing to the SmartChan
	err := sc.Write("Test Data")
	if err != nil {
		t.Errorf("Unexpected error while writing to the channel: %v", err)
	}
	if sc.Count() != 1 {
		t.Errorf("Count should be 1, got %d instead", sc.Count())
	}

	// Test reading from the SmartChan
	data, err := sc.Read()
	if err != nil {
		t.Errorf("Unexpected error while reading from the channel: %v", err)
	}
	if data != "Test Data" {
		t.Errorf("Expected 'Test Data', got %v instead", data)
	}

	// Test closing the SmartChan
	sc.Close()
	if sc.CanWrite() {
		t.Errorf("Expected CanWrite to return false, got true instead")
	}

	// Test writing to the closed SmartChan
	err = sc.Write("More Data")
	if err == nil {
		t.Error("Expected an error when writing to a closed channel, got nil instead")
	}

	// Test reading from the closed SmartChan
	_, err = sc.Read()
	if err == nil {
		t.Error("Expected an error when reading from a closed channel, got nil instead")
	}
}

func TestConcurrentAccess(t *testing.T) {
	sc := NewSmartChan(100)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	// Concurrent writing to the SmartChan
	go func() {
		for i := 0; i < 100; i++ {
			if err := sc.Write(i); err != nil {
				t.Errorf("Unexpected error while writing to the channel: %v", err)
			}
		}
		wg.Done()
	}()

	// Concurrent reading from the SmartChan
	wg.Add(1)
	go func() {
		for i := 0; i < 100; i++ {
			data, err := sc.Read()
			if err != nil {
				t.Errorf("Unexpected error while reading from the channel: %v", err)
			}
			if data != i {
				t.Errorf("Expected %d, got %v instead", i, data)
			}
		}
		wg.Done()
	}()

	// Wait for both goroutines to finish
	wg.Wait()

	if sc.Count() != 100 {
		t.Errorf("Count should be 100, got %d instead", sc.Count())
	}
	sc.Close()
}
