package go_smartchan

import (
	"sync"
	"testing"
)

func TestDataFlow(t *testing.T) {
	sc := NewSmartChan(5)

	err := sc.Write("Test Data")
	if err != nil {
		t.Errorf("Unexpected error while writing to the channel: %v", err)
	}
	if sc.Count() != 1 {
		t.Errorf("Count should be 1, got %d instead", sc.Count())
	}

	data, err := sc.Read()
	if err != nil {
		t.Errorf("Unexpected error while reading from the channel: %v", err)
	}
	if data != "Test Data" {
		t.Errorf("Expected 'Test Data', got %v instead", data)
	}

	sc.Close()
	if sc.CanWrite() {
		t.Errorf("Expected CanWrite to return false, got true instead")
	}

	err = sc.Write("More Data")
	if err == nil {
		t.Error("Expected an error when writing to a closed channel, got nil instead")
	}

	_, err = sc.Read()
	if err == nil {
		t.Error("Expected an error when reading from a closed channel, got nil instead")
	}
}

func TestConcurrentAccess(t *testing.T) {
	sc := NewSmartChan(100)
	wg1 := &sync.WaitGroup{}
	wg2 := &sync.WaitGroup{}

	wg1.Add(1)
	go func() {
		for i := 0; i < 100; i++ {
			if err := sc.Write(i); err != nil {
				t.Errorf("Unexpected error while writing to the channel: %v", err)
			}
		}
		wg1.Done()
	}()

	wg1.Wait()

	wg2.Add(1)
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
		wg2.Done()
	}()

	wg2.Wait()

	sc.Close()
}

func TestChan(t *testing.T) {
	sc := NewSmartChan(1)

	err := sc.Write(1)
	if err != nil {
		t.Fatalf("Unexpected error while writing to the channel: %v", err)
	}

	val, err := sc.Read()
	if err != nil {
		t.Fatalf("Unexpected error while reading from the channel: %v", err)
	}

	if val != 1 {
		t.Errorf("Expected 1, got %v instead", val)
	}
}
