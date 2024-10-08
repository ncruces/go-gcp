package gmutex

import (
	"context"
	"math/rand"
	"time"
)

// Full Jitter from:
// https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/

const (
	backOffMin = 50 * time.Millisecond
	backOffMax = 30 * time.Second
)

type expBackOff struct {
	time time.Duration
}

type linBackOff struct {
	time time.Duration
}

func (b *linBackOff) wait(ctx context.Context) error {
	b.time += backOffMin
	if b.time < backOffMin {
		b.time = backOffMin
	}
	if b.time > backOffMax {
		b.time = backOffMax
	}
	return wait(ctx, time.Duration(rand.Int63n(int64(b.time))))
}

func (b *expBackOff) wait(ctx context.Context) error {
	b.time += b.time / 2
	if b.time < backOffMin {
		b.time = backOffMin
	}
	if b.time > backOffMax {
		b.time = backOffMax
	}
	return wait(ctx, time.Duration(rand.Int63n(int64(b.time))))
}

func wait(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		timer.Stop()
		return ctx.Err()
	}
}
