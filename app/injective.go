package app

import (
	"context"
	"fmt"
	"sync"
	"time"

	"bharvest.io/orchmon/client/grpc"
	"bharvest.io/orchmon/log"
	"bharvest.io/orchmon/server"
	"bharvest.io/orchmon/tg"
)

func (cfg *Config) CheckInjective(ctx context.Context) error {
	grpcQuery := grpc.NewInjective(cfg.Injective.GRPC)
	err := grpcQuery.Connect(ctx)
	if err != nil {
		log.Error(err)
		return err
	}
	defer func() {
		err := grpcQuery.Terminate(ctx)
		if err != nil {
			log.Error(err)
			return
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(4)

	// For get # pending batch size tx
	var pendingBatchSize int
	go func() {
		defer wg.Done()

		// Three re-checks with 2 seconds interval, because pending TXs can exist temporary
		for i:=0; i<3; i++ {
			pendingBatchSize, err = grpcQuery.GetLastPendingBatchSize(ctx, cfg.Injective.Wallet.PrintAcc())
			if err != nil {
				log.Error(err)
				return
			}

			if pendingBatchSize == 0 {
				break
			}

			time.Sleep(2*time.Second)
		}
	}()
	
	// For get # pending valset tx
	var pendingValsetSize int
	go func() {
		defer wg.Done()

		// Three re-checks with 2 seconds interval, because pending TXs can exist temporary
		for i:=0; i<3; i++ {
			pendingValsetSize, err = grpcQuery.GetLastPendingValsetSize(ctx, cfg.Injective.Wallet.PrintAcc())
			if err != nil {
				log.Error(err)
				return
			}

			if pendingValsetSize == 0 {
				break
			}

			time.Sleep(2*time.Second)
		}
	}()

	// For get event nonce
	var eventNonce uint64
	go func() {
		defer wg.Done()

		eventNonce, err = grpcQuery.GetLastEventNonce(ctx, cfg.Injective.Wallet.PrintAcc())
		if err != nil {
			log.Error(err)
			return
		}
	}()

	// For get observed event nonce
	var observedEventNonce uint64
	go func() {
		defer wg.Done()

		observedEventNonce, err = grpcQuery.GetLastObservedNonce(ctx)
		if err != nil {
			log.Error(err)
			return
		}
	}()

	wg.Wait()

	server.GlobalState.InjectivePeggo.PendingBatchSize = pendingBatchSize
	server.GlobalState.InjectivePeggo.PendingValsetSize = pendingValsetSize
	server.GlobalState.InjectivePeggo.Nonce = fmt.Sprintf("%d / %d", eventNonce, observedEventNonce)

	check := true

	if pendingBatchSize != 0 || pendingValsetSize != 0 {
		check = false
	}

	// Nonce could be different each other
	// even though it doesn't any issue in peggo.
	// So it have three chance for change status to false.
	if observedEventNonce != eventNonce {
		if server.GlobalState.InjectivePeggo.NonceFalseCount == 3 {
			check = false
		} else {
			server.GlobalState.InjectivePeggo.NonceFalseCount += 1
		}
	} else {
		// Set nonce false count to 0
		// if nonce is same each other.
		server.GlobalState.InjectivePeggo.NonceFalseCount = 0
	}

	server.GlobalState.InjectivePeggo.Status = check
	m := "Injective Peggo status: "
	if check {
		m += "🟢"
		log.Info(m)
	} else {
		m += "🛑"
		log.Info(m)
		tg.SendMsg(m)
	}
	
	return nil
}
