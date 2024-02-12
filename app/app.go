package app

import (
	"context"
	"sync"
	"time"

	"bharvest.io/orchmon/log"
)

func Run(ctx context.Context, cfg *Config) {
	ctx, cancel := context.WithTimeout(ctx, 25*time.Second)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)

	// Check Injective
	go func() {
		defer wg.Done()
		if !cfg.Injective.Enable {
			return
		}

		err := cfg.CheckInjective(ctx)
		if err != nil {
			log.Error(err)
			return
		}
	}()

	wg.Wait()

	return
}
