module github.com/ChainSafe/chainbridge-optimism-module

go 1.16

require (
	github.com/ChainSafe/chainbridge-core v0.0.0-20220113021106-23cb35063d74 // branch maxim/optimism-e2e (GOPROXY=direct go get github.com/ChainSafe/chainbridge-core@maxim/optimism-e2e)
	github.com/centrifuge/go-substrate-rpc-client v2.0.0+incompatible
	github.com/ethereum/go-ethereum v1.10.15
	github.com/rs/zerolog v1.26.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.10.1
	github.com/stretchr/testify v1.7.0
)
