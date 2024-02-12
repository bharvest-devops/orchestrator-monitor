package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"bharvest.io/orchmon/app"
	"bharvest.io/orchmon/client/grpc"
	"bharvest.io/orchmon/log"
	"bharvest.io/orchmon/server"
	"bharvest.io/orchmon/tg"
	"bharvest.io/orchmon/wallet"
	"github.com/BurntSushi/toml"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	ctx := context.Background()

	cfgPath := flag.String("config", "", "Config file")
    flag.Parse()
	if *cfgPath == "" {
		panic("Error: Please input config file path with -config flag.")
	}

	f, err := os.ReadFile(*cfgPath)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	cfg := app.Config{}
	err = toml.Unmarshal(f, &cfg)
	if err != nil {
		log.Error(err)
		panic(err)
	}

	if cfg.Injective.Enable {
		err = PrePrepareForInjective(ctx, &cfg)
		if err != nil {
			log.Error(err)
			panic(err)
		}
	}

	tgTitle := fmt.Sprintf("🤖 Orchestrator Monitor for %s 🤖", cfg.General.Network)
	tg.SetTg(cfg.Tg.Enable, tgTitle, cfg.Tg.Token, cfg.Tg.ChatID)

	go server.Run(cfg.General.ListenPort)
	for {
		app.Run(ctx, &cfg)
		time.Sleep(time.Duration(cfg.General.Period) * time.Minute)
	}
}

func PrePrepareForInjective(ctx context.Context, cfg *app.Config) error {
	grpcQuery := grpc.NewInjective(cfg.Injective.GRPC)
	err := grpcQuery.Connect(ctx)
	if err != nil {
		log.Error(err)
		return err
	}
	defer func() {
		err = grpcQuery.Terminate(ctx)
		if err != nil {
			log.Error(err)
			return
		}
	}()

	tmpWallet, err := wallet.NewWallet(ctx, cfg.Injective.ValidatorAcc)
	if err != nil {
		log.Error(err)
		return err
	}

	// Get Orchestrator Address from Validator Address
	OrchAcc, err := grpcQuery.GetDelegateKeyByValidator(ctx, tmpWallet.PrintValoper())
	if err != nil {
		log.Error(err)
		return err
	}
	cfg.Injective.Wallet, err = wallet.NewWallet(ctx, OrchAcc)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}