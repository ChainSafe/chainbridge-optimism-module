module github.com/ChainSafe/chainbridge-optimism-module

go 1.16

require (
	github.com/ChainSafe/chainbridge-core v0.0.0-20220117151815-684ae464fb8d // branch maxim/optimism-e2e (GOPROXY=direct go get github.com/ChainSafe/chainbridge-core@maxim/optimism-e2e)
	github.com/ethereum/go-ethereum v1.10.15
	github.com/mitchellh/mapstructure v1.4.3
	github.com/rs/zerolog v1.26.1
	github.com/spf13/cobra v1.6.0
	github.com/spf13/viper v1.10.1
	github.com/stretchr/testify v1.7.0
)
