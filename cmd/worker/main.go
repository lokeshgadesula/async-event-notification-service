package main

import (
	"context"
	"fmt"
	"github.com/lokeshgadesula/async-event-notification-service/internal/event"
	"github.com/lokeshgadesula/async-event-notification-service/internal/service"
	"github.com/redis/go-redis/v9"
	"os"
	"os/signal"
	"syscall"
)

type logBroadcaster struct{}

func (logBroadcaster) Broadcast(_ context.Context, e event.Event) error {
	fmt.Printf("notify %s: %s\n", e.Recipient, e.Message)
	return nil
}
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	r := redis.NewClient(&redis.Options{Addr: env("REDIS_ADDR", "redis:6379")})
	s := service.Service{R: r, B: logBroadcaster{}, Workers: 8, Queue: "events", Dead: "events:dead"}
	if err := s.Run(ctx); err != nil {
		panic(err)
	}
}
func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
