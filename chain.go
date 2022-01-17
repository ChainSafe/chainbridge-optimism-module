// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only
package optimism

import (
	"github.com/ChainSafe/chainbridge-core/chains/evm"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/bridge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmgaspricer"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor/signAndSend"
	"github.com/ChainSafe/chainbridge-core/store"
	"github.com/ChainSafe/chainbridge-optimism-module/config"
	"github.com/ChainSafe/chainbridge-optimism-module/optimismclient"

	"github.com/ChainSafe/chainbridge-core/chains/evm/listener"
	"github.com/ChainSafe/chainbridge-core/chains/evm/voter"
	"github.com/ethereum/go-ethereum/common"
)

// SetupDefaultOptimismChain sets up an instance of EVMChain with optimism client setup that
// tracks only verified tx batches.
func SetupDefaultOptimismChain(rawConfig map[string]interface{}, txFabric calls.TxFabric, db *store.BlockStore) (*evm.EVMChain, error) {
	config, err := config.NewOptimismConfig(rawConfig)
	if err != nil {
		return nil, err
	}

	client, err := optimismclient.NewOptimismClient(config)
	if err != nil {
		return nil, err
	}

	gasPricer := evmgaspricer.NewLondonGasPriceClient(client, nil)
	t := signAndSend.NewSignAndSendTransactor(txFabric, gasPricer, client)
	bridgeContract := bridge.NewBridgeContract(client, common.HexToAddress(config.Bridge), t)

	eventHandler := listener.NewETHEventHandler(*bridgeContract)
	eventHandler.RegisterEventHandler(config.Erc20Handler, listener.Erc20EventHandler)
	eventHandler.RegisterEventHandler(config.Erc721Handler, listener.Erc721EventHandler)
	eventHandler.RegisterEventHandler(config.GenericHandler, listener.GenericEventHandler)
	evmListener := listener.NewEVMListener(client, eventHandler, common.HexToAddress(config.Bridge))

	mh := voter.NewEVMMessageHandler(*bridgeContract)
	mh.RegisterMessageHandler(config.Erc20Handler, voter.ERC20MessageHandler)
	mh.RegisterMessageHandler(config.Erc721Handler, voter.ERC721MessageHandler)
	mh.RegisterMessageHandler(config.GenericHandler, voter.GenericMessageHandler)

	evmVoter, err := voter.NewVoterWithSubscription(mh, client, bridgeContract)
	if err != nil {
		return nil, err
	}

	return evm.NewEVMChain(evmListener, evmVoter, db, &config.EVMConfig), nil
}
