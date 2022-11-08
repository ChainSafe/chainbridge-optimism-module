module github.com/ChainSafe/chainbridge-optimism-module

go 1.16

require (
	github.com/ChainSafe/chainbridge-core v0.0.0-20220117151815-684ae464fb8d // branch maxim/optimism-e2e (GOPROXY=direct go get github.com/ChainSafe/chainbridge-core@maxim/optimism-e2e)
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/ethereum/go-ethereum v1.10.15
	github.com/mitchellh/mapstructure v1.5.0
	github.com/rs/zerolog v1.26.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.14.0
	github.com/stretchr/testify v1.8.1
)
