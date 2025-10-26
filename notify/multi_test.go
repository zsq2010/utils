package notify

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

type mockNotifier struct {
	sendFunc func(message Message) error
	calls    int
	mu       sync.Mutex
}

func (m *mockNotifier) Send(message Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls++
	if m.sendFunc != nil {
		return m.sendFunc(message)
	}
	return nil
}

func (m *mockNotifier) getCalls() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.calls
}

func TestNewMulti(t *testing.T) {
	mock1 := &mockNotifier{}
	mock2 := &mockNotifier{}

	multi := NewMulti(mock1, mock2)
	if multi == nil {
		t.Fatal("NewMulti returned nil")
	}
	if len(multi.notifiers) != 2 {
		t.Errorf("len(notifiers) = %d, want 2", len(multi.notifiers))
	}
	if multi.parallel {
		t.Error("NewMulti should create sequential notifier")
	}
}

func TestNewMultiParallel(t *testing.T) {
	mock1 := &mockNotifier{}
	mock2 := &mockNotifier{}

	multi := NewMultiParallel(mock1, mock2)
	if multi == nil {
		t.Fatal("NewMultiParallel returned nil")
	}
	if !multi.parallel {
		t.Error("NewMultiParallel should create parallel notifier")
	}
}

func TestMultiNotifier_Send_noNotifiers(t *testing.T) {
	multi := NewMulti()

	message := Message{Title: "Test", Body: "Body"}
	err := multi.Send(message)

	if err == nil {
		t.Error("expected error for no notifiers, got nil")
	}
	if !strings.Contains(err.Error(), "no notifiers") {
		t.Errorf("error message = %v, want to contain 'no notifiers'", err)
	}
}

func TestMultiNotifier_Send_sequential_success(t *testing.T) {
	mock1 := &mockNotifier{}
	mock2 := &mockNotifier{}
	mock3 := &mockNotifier{}

	multi := NewMulti(mock1, mock2, mock3)

	message := Message{Title: "Test", Body: "Body"}
	err := multi.Send(message)

	if err != nil {
		t.Errorf("Send failed: %v", err)
	}

	if mock1.getCalls() != 1 {
		t.Errorf("mock1 calls = %d, want 1", mock1.getCalls())
	}
	if mock2.getCalls() != 1 {
		t.Errorf("mock2 calls = %d, want 1", mock2.getCalls())
	}
	if mock3.getCalls() != 1 {
		t.Errorf("mock3 calls = %d, want 1", mock3.getCalls())
	}
}

func TestMultiNotifier_Send_sequential_stopOnError(t *testing.T) {
	mock1 := &mockNotifier{}
	mock2 := &mockNotifier{
		sendFunc: func(message Message) error {
			return fmt.Errorf("mock error")
		},
	}
	mock3 := &mockNotifier{}

	multi := NewMulti(mock1, mock2, mock3)

	message := Message{Title: "Test", Body: "Body"}
	err := multi.Send(message)

	if err == nil {
		t.Error("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "notifier 1 failed") {
		t.Errorf("error message = %v, want to contain 'notifier 1 failed'", err)
	}

	if mock1.getCalls() != 1 {
		t.Errorf("mock1 calls = %d, want 1", mock1.getCalls())
	}
	if mock2.getCalls() != 1 {
		t.Errorf("mock2 calls = %d, want 1", mock2.getCalls())
	}
	if mock3.getCalls() != 0 {
		t.Errorf("mock3 calls = %d, want 0 (should not be called after error)", mock3.getCalls())
	}
}

func TestMultiNotifier_Send_parallel_success(t *testing.T) {
	mock1 := &mockNotifier{
		sendFunc: func(message Message) error {
			time.Sleep(10 * time.Millisecond)
			return nil
		},
	}
	mock2 := &mockNotifier{
		sendFunc: func(message Message) error {
			time.Sleep(10 * time.Millisecond)
			return nil
		},
	}
	mock3 := &mockNotifier{
		sendFunc: func(message Message) error {
			time.Sleep(10 * time.Millisecond)
			return nil
		},
	}

	multi := NewMultiParallel(mock1, mock2, mock3)

	message := Message{Title: "Test", Body: "Body"}
	start := time.Now()
	err := multi.Send(message)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("Send failed: %v", err)
	}

	if duration > 50*time.Millisecond {
		t.Errorf("parallel send took %v, expected < 50ms (should run concurrently)", duration)
	}

	if mock1.getCalls() != 1 {
		t.Errorf("mock1 calls = %d, want 1", mock1.getCalls())
	}
	if mock2.getCalls() != 1 {
		t.Errorf("mock2 calls = %d, want 1", mock2.getCalls())
	}
	if mock3.getCalls() != 1 {
		t.Errorf("mock3 calls = %d, want 1", mock3.getCalls())
	}
}

func TestMultiNotifier_Send_parallel_collectErrors(t *testing.T) {
	mock1 := &mockNotifier{}
	mock2 := &mockNotifier{
		sendFunc: func(message Message) error {
			return fmt.Errorf("error from mock2")
		},
	}
	mock3 := &mockNotifier{
		sendFunc: func(message Message) error {
			return fmt.Errorf("error from mock3")
		},
	}

	multi := NewMultiParallel(mock1, mock2, mock3)

	message := Message{Title: "Test", Body: "Body"}
	err := multi.Send(message)

	if err == nil {
		t.Error("expected error, got nil")
	}

	errMsg := err.Error()
	if !strings.Contains(errMsg, "notifier 1") {
		t.Errorf("error message should contain 'notifier 1': %v", errMsg)
	}
	if !strings.Contains(errMsg, "notifier 2") {
		t.Errorf("error message should contain 'notifier 2': %v", errMsg)
	}

	if mock1.getCalls() != 1 {
		t.Errorf("mock1 calls = %d, want 1", mock1.getCalls())
	}
	if mock2.getCalls() != 1 {
		t.Errorf("mock2 calls = %d, want 1", mock2.getCalls())
	}
	if mock3.getCalls() != 1 {
		t.Errorf("mock3 calls = %d, want 1 (all notifiers should be called in parallel)", mock3.getCalls())
	}
}

func TestMultiNotifier_Send_parallel_partialFailure(t *testing.T) {
	mock1 := &mockNotifier{}
	mock2 := &mockNotifier{
		sendFunc: func(message Message) error {
			return fmt.Errorf("error from mock2")
		},
	}
	mock3 := &mockNotifier{}

	multi := NewMultiParallel(mock1, mock2, mock3)

	message := Message{Title: "Test", Body: "Body"}
	err := multi.Send(message)

	if err == nil {
		t.Error("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "notifier 1") {
		t.Errorf("error message = %v, want to contain 'notifier 1'", err)
	}

	if mock1.getCalls() != 1 {
		t.Errorf("mock1 calls = %d, want 1", mock1.getCalls())
	}
	if mock2.getCalls() != 1 {
		t.Errorf("mock2 calls = %d, want 1", mock2.getCalls())
	}
	if mock3.getCalls() != 1 {
		t.Errorf("mock3 calls = %d, want 1", mock3.getCalls())
	}
}
