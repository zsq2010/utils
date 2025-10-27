package notify

import (
	"fmt"
	"strings"
	"sync"
)

// MultiNotifier implements the Notifier interface and sends notifications to multiple channels.
type MultiNotifier struct {
	notifiers []Notifier
	parallel  bool
}

// NewMulti creates a new MultiNotifier that sends to all provided notifiers sequentially.
func NewMulti(notifiers ...Notifier) *MultiNotifier {
	return &MultiNotifier{
		notifiers: notifiers,
		parallel:  false,
	}
}

// NewMultiParallel creates a new MultiNotifier that sends to all provided notifiers in parallel.
func NewMultiParallel(notifiers ...Notifier) *MultiNotifier {
	return &MultiNotifier{
		notifiers: notifiers,
		parallel:  true,
	}
}

// Send sends the message to all configured notifiers.
// In sequential mode, it stops at the first error.
// In parallel mode, it collects all errors and returns them together.
func (m *MultiNotifier) Send(message Message) error {
	if len(m.notifiers) == 0 {
		return fmt.Errorf("no notifiers configured")
	}

	if m.parallel {
		return m.sendParallel(message)
	}

	return m.sendSequential(message)
}

// sendSequential sends the message to notifiers one by one, stopping at the first error.
func (m *MultiNotifier) sendSequential(message Message) error {
	for i, notifier := range m.notifiers {
		if err := notifier.Send(message); err != nil {
			return fmt.Errorf("notifier %d failed: %w", i, err)
		}
	}
	return nil
}

// sendParallel sends the message to all notifiers concurrently and collects errors.
func (m *MultiNotifier) sendParallel(message Message) error {
	var wg sync.WaitGroup
	errors := make([]error, len(m.notifiers))

	for i, notifier := range m.notifiers {
		wg.Add(1)
		go func(idx int, n Notifier) {
			defer wg.Done()
			errors[idx] = n.Send(message)
		}(i, notifier)
	}

	wg.Wait()

	var failedErrors []string
	for i, err := range errors {
		if err != nil {
			failedErrors = append(failedErrors, fmt.Sprintf("notifier %d: %v", i, err))
		}
	}

	if len(failedErrors) > 0 {
		return fmt.Errorf("multi-send failed: %s", strings.Join(failedErrors, "; "))
	}

	return nil
}
