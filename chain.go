// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only
package optimism

import (
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/bridge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmgaspricer"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/chainbridge-optimism-module/config"
	"github.com/ChainSafe/chainbridge-optimism-module/optimismclient"

	"github.com/ChainSafe/chainbridge-core/blockstore"
	"github.com/ChainSafe/chainbridge-core/chains/evm"
	"github.com/ChainSafe/chainbridge-core/chains/evm/listener"
	"github.com/ChainSafe/chainbridge-core/chains/evm/voter"
	"github.com/ChainSafe/chainbridge-core/config/chain"
	"github.com/ethereum/go-ethereum/common"
)

// OptimismChain is struct that aggregates all data required for
type OptimismChain struct {
	evm.EVMChain
	listener evm.EventListener
	writer   evm.ProposalVoter
	kvdb     blockstore.KeyValueReaderWriter
	config   *config.OptimismConfig
}

// SetupDefaultOptimismChain sets up an OptimismChain with all supported handlers configured
func SetupDefaultOptimismChain(rawConfig map[string]interface{}, txFabric calls.TxFabric, db blockstore.KeyValueReaderWriter) (*OptimismChain, error) {
	config, err := config.NewOptimismConfig(rawConfig)
	if err != nil {
		return nil, err
	}

	client, err := optimismclient.NewOptimismClient(config)
	if err != nil {
		return nil, err
	}

	gasPricer := evmgaspricer.NewLondonGasPriceClient(client, nil)
	t := transactor.NewSignAndSendTransactor(txFabric, gasPricer, client)
	bridgeContract := bridge.NewBridgeContract(client, common.HexToAddress(config.EVMConfig.Bridge), t)

	eventHandler := listener.NewETHEventHandler(*bridgeContract)
	eventHandler.RegisterEventHandler(config.EVMConfig.Erc20Handler, listener.Erc20EventHandler)
	eventHandler.RegisterEventHandler(config.EVMConfig.Erc721Handler, listener.Erc721EventHandler)
	eventHandler.RegisterEventHandler(config.EVMConfig.GenericHandler, listener.GenericEventHandler)
	evmListener := listener.NewEVMListener(client, eventHandler, common.HexToAddress(config.EVMConfig.Bridge))

	mh := voter.NewEVMMessageHandler(*bridgeContract)
	mh.RegisterMessageHandler(config.EVMConfig.Erc20Handler, voter.ERC20MessageHandler)
	mh.RegisterMessageHandler(config.EVMConfig.Erc721Handler, voter.ERC721MessageHandler)
	mh.RegisterMessageHandler(config.EVMConfig.GenericHandler, voter.GenericMessageHandler)

	evmVoter, err := voter.NewVoterWithSubscription(mh, client, bridgeContract)
	if err != nil {
		return nil, err
	}

	return NewOptimismChain(evmListener, evmVoter, db, config), nil
}

func NewOptimismChain(listener evm.EventListener, writer evm.ProposalVoter, kvdb blockstore.KeyValueReaderWriter, config *chain.OptimismConfig) *OptimismChain {
	return &OptimismChain{listener: listener, writer: writer, kvdb: kvdb, config: config}
}
