package service

import (
	"context"
	"github.com/lokeshgadesula/async-event-notification-service/internal/event"
	"sync"
	"testing"
)

type mem struct {
	mu sync.Mutex
	n  int
}

func (m *mem) Broadcast(_ context.Context, _ event.Event) error {
	m.mu.Lock()
	m.n++
	m.mu.Unlock()
	return nil
}
func TestConcurrentBroadcaster(t *testing.T) {
	m := &mem{}
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); _ = m.Broadcast(context.Background(), event.Event{}) }()
	}
	wg.Wait()
	if m.n != 1000 {
		t.Fatalf("got %d", m.n)
	}
}
func TestEventShape(t *testing.T) {
	e := event.Event{ID: "1", Type: "meeting.action"}
	if e.Type != "meeting.action" {
		t.Fatal("bad type")
	}
}
