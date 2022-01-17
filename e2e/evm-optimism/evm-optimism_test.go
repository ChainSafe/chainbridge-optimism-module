package main

import (
	"testing"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmtransaction"
	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/ChainSafe/chainbridge-core/e2e/evm"
	"github.com/ChainSafe/chainbridge-optimism-module/optimismclient"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/local"
	"github.com/stretchr/testify/suite"
)

const ETHEndpoint1 = "ws://localhost:8646"
const OptimismEndpoint1 = "ws://localhost:8550"
const VerifierEndpoint1 = "ws://localhost:8550"

// Funded optimism address: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
const fundedOptimismPk = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

var OptimismRelayerAddresses = []common.Address{
	common.HexToAddress("0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC"),
	common.HexToAddress("0x90F79bf6EB2c4f870365E785982E1f101E93b906"),
	common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8"),
}

// Alice key is used by the relayer, Eve key is used as admin and depositter
func TestRunE2ETests(t *testing.T) {
	ethClient, err := evmclient.NewEVMClientFromParams(ETHEndpoint1, local.EveKp.PrivateKey())
	if err != nil {
		panic(err)
	}

	kp, err := secp256k1.NewKeypairFromString(fundedOptimismPk)
	if err != nil {
		panic(err)
	}
	optimismClient, err := optimismclient.NewOptimismClientFromParams(OptimismEndpoint1, kp.PrivateKey(), true, VerifierEndpoint1)
	if err != nil {
		panic(err)
	}

	suite.Run(t, evm.SetupEVM2EVMTestSuite(
		evmtransaction.NewTransaction,
		evmtransaction.NewTransaction,
		ethClient,
		optimismClient,
		local.DefaultRelayerAddresses,
		OptimismRelayerAddresses,
	))
}
