package config_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"
	"github.com/ChainSafe/chainbridge-core/config/chain"
	"github.com/ChainSafe/chainbridge-optimism-module/config"
	"github.com/stretchr/testify/suite"
)

type NewOptimismConfigTestSuite struct {
	suite.Suite
}

func TestRunNewOptimismConfigTestSuite(t *testing.T) {
	suite.Run(t, new(NewOptimismConfigTestSuite))
}

func (s *NewOptimismConfigTestSuite) SetupSuite()    {}
func (s *NewOptimismConfigTestSuite) TearDownSuite() {}
func (s *NewOptimismConfigTestSuite) SetupTest()     {}
func (s *NewOptimismConfigTestSuite) TearDownTest()  {}

func (s *NewOptimismConfigTestSuite) Test_FailedEVMConfigValidation() {
	_, err := config.NewOptimismConfig(map[string]interface{}{
		"id":       1,
		"endpoint": "ws://domain.com",
		"name":     "evm1",
		"from":     "address",
	})

	s.NotNil(err)
}

func (s *NewOptimismConfigTestSuite) Test_FailedOptimismDecode() {
	_, err := config.NewOptimismConfig(map[string]interface{}{
		"id":           1,
		"endpoint":     "ws://domain.com",
		"name":         "evm1",
		"from":         "address",
		"bridge":       "bridgeAddress",
		"verifyRollup": "invalid",
	})

	s.NotNil(err)
}

func (s *NewOptimismConfigTestSuite) Test_ValidConfig() {
	rawConfig := map[string]interface{}{
		"id":               1,
		"endpoint":         "ws://domain.com",
		"name":             "evm1",
		"from":             "address",
		"bridge":           "bridgeAddress",
		"verifyRollup":     true,
		"verifierEndpoint": "http://verifier",
	}

	actualConfig, err := config.NewOptimismConfig(rawConfig)

	id := new(uint8)
	*id = 1
	s.Nil(err)
	s.Equal(*actualConfig, config.OptimismConfig{
		EVMConfig: chain.EVMConfig{
			GeneralChainConfig: chain.GeneralChainConfig{
				Name:     "evm1",
				From:     "address",
				Endpoint: "ws://domain.com",
				Id:       id,
			},
			Bridge:             "bridgeAddress",
			Erc20Handler:       "",
			Erc721Handler:      "",
			GenericHandler:     "",
			GasLimit:           big.NewInt(consts.DefaultGasLimit),
			MaxGasPrice:        big.NewInt(consts.DefaultGasPrice),
			GasMultiplier:      big.NewFloat(consts.DefaultGasMultiplier),
			StartBlock:         big.NewInt(0),
			BlockConfirmations: big.NewInt(consts.DefaultBlockConfirmations),
			BlockRetryInterval: time.Duration(5) * time.Second,
		},
		VerifyRollup:     true,
		VerifierEndpoint: "http://verifier",
	})
}
