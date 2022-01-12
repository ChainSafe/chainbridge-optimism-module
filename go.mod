module github.com/ChainSafe/chainbridge-optimism-module

go 1.16

// To be removed when optimism config is inside core
replace github.com/ChainSafe/chainbridge-core => ../chainbridge-core

require (
	github.com/ChainSafe/chainbridge-core v0.0.0-20220110124723-abb0bf918502
	github.com/ethereum/go-ethereum v1.10.15
	github.com/rs/zerolog v1.26.1
	github.com/spf13/viper v1.10.1 // indirect
)
