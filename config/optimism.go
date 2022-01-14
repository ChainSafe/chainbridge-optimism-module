package config

import (
	"github.com/ChainSafe/chainbridge-core/config/chain"
	"github.com/mitchellh/mapstructure"
)

type OptimismConfig struct {
	EVMConfig        chain.EVMConfig // Contains configuration of Optimism l2geth client which is used to for sequencing transactions to Optimism
	VerifyRollup     bool
	VerifierEndpoint string // This is the endpoint for an Optimism replica which is read-only and used only for verifying transactions
}

type RawOptimismConfig struct {
	chain.RawEVMConfig `mapstructure:",squash"`
	VerifyRollup       bool   `mapstructure:"verifyRollup"`
	VerifierEndpoint   string `mapstructure:"verifierEndpoint"`
}

// NewOptimismConfig decodes and validates an instance of an OptimismConfig from
// raw chain config
func NewOptimismConfig(chainConfig map[string]interface{}) (*OptimismConfig, error) {
	evmConfig, err := chain.NewEVMConfig(chainConfig)
	if err != nil {
		return nil, err
	}

	var c RawOptimismConfig
	err = mapstructure.Decode(chainConfig, &c)
	if err != nil {
		return nil, err
	}

	config := &OptimismConfig{
		EVMConfig:        *evmConfig,
		VerifyRollup:     c.VerifyRollup,
		VerifierEndpoint: c.VerifierEndpoint,
	}

	return config, nil
}
