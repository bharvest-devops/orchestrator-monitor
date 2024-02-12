package app

import "bharvest.io/orchmon/wallet"

type Config struct {
	General struct {
		Network    string `toml:"network"`
		EthAPI     string `toml:"eth_api"`
		ListenPort int    `toml:"listen_port"`
		Period     int    `toml:"period"`
	} `toml:"general"`
	Injective struct {
		Enable       bool   `toml:"enable"`
		GRPC         string `toml:"grpc"`
		ValidatorAcc string `toml:"validator_acc"`
		Wallet       *wallet.Wallet
	} `toml:"injective"`
	Tg struct {
		Enable bool   `toml:"enable"`
		Token  string `toml:"token"`
		ChatID string `toml:"chat_id"`
	} `toml:"tg"`
}
