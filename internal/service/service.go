package service

import (
	"context"
	"encoding/json"
	"github.com/lokeshgadesula/async-event-notification-service/internal/event"
	"github.com/redis/go-redis/v9"
	"sync"
)

type Broadcaster interface {
	Broadcast(context.Context, event.Event) error
}
type Service struct {
	R       *redis.Client
	B       Broadcaster
	Workers int
	Queue   string
	Dead    string
}

func (s Service) Run(ctx context.Context) error {
	jobs := make(chan event.Event, s.Workers*2)
	var wg sync.WaitGroup
	for i := 0; i < s.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for e := range jobs {
				if err := s.B.Broadcast(ctx, e); err != nil {
					e.Attempts++
					raw, _ := json.Marshal(e)
					if e.Attempts >= 3 {
						s.R.LPush(ctx, s.Dead, raw)
					} else {
						s.R.LPush(ctx, s.Queue, raw)
					}
				}
			}
		}()
	}
	defer func() { close(jobs); wg.Wait() }()
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			x, err := s.R.BRPop(ctx, 0, s.Queue).Result()
			if err != nil {
				if ctx.Err() != nil {
					return nil
				}
				return err
			}
			var e event.Event
			if json.Unmarshal([]byte(x[1]), &e) != nil {
				s.R.LPush(ctx, s.Dead, x[1])
				continue
			}
			jobs <- e
		}
	}
}
