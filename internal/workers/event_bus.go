package workers

import (
	"github.com/google/uuid"
	"sync"
	"tahap2/internal/domain"
)

const (
	EventTypeTransfer = "transfer"
)

type EventBus struct {
	subscribers map[string][]chan interface{}
	mutex       sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]chan interface{}),
	}
}

// Subscribe to an event type
func (eb *EventBus) Subscribe(eventType string) chan interface{} {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	ch := make(chan interface{}, 1) // Buffered channel to avoid blocking
	eb.subscribers[eventType] = append(eb.subscribers[eventType], ch)
	return ch
}

// Publish an event to all subscribers
func (eb *EventBus) Publish(eventType string, data interface{}) {
	eb.mutex.RLock()
	defer eb.mutex.RUnlock()

	if subs, found := eb.subscribers[eventType]; found {
		for _, ch := range subs {
			ch <- data
		}
	}
}

// Close all channels
func (eb *EventBus) Close() {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()

	for _, subs := range eb.subscribers {
		for _, ch := range subs {
			close(ch)
		}
	}
}

type TransferParam struct {
	TransferInfo domain.Transaction
	TargetID     uuid.UUID
}
