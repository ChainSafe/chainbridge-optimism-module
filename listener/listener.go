// Copyright 2021 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package listener

import (
	"context"
	"math/big"
	"time"

	"github.com/ChainSafe/chainbridge-core/blockstore"
	"github.com/ChainSafe/chainbridge-core/relayer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

var BlockRetryInterval = time.Second * 5

// Only 1 for Optimism as block confirmations check occur when events are indexed in the data transport layer
var BlockDelay = big.NewInt(5) //TODO: move to config

type DepositLogs struct {
	DestinationChainID uint8
	ResourceID         [32]byte
	DepositNonce       uint64
}

type ChainClient interface {
	LatestBlock() (*big.Int, error)
	FetchDepositLogs(ctx context.Context, address common.Address, startBlock *big.Int, endBlock *big.Int) ([]*DepositLogs, error)
	CallContract(ctx context.Context, callArgs map[string]interface{}, blockNumber *big.Int) ([]byte, error)
	IsRollupVerified(blockNumber uint64) (bool, error)
}

type EventHandler interface {
	HandleEvent(sourceID, destID uint8, nonce uint64, rID [32]byte) (*relayer.Message, error)
}

type EVMListener struct {
	chainReader   ChainClient
	eventHandler  EventHandler
	bridgeAddress common.Address
}

func NewEVMListener(chainReader ChainClient, handler EventHandler, bridgeAddress common.Address) *EVMListener {
	return &EVMListener{chainReader: chainReader, eventHandler: handler, bridgeAddress: bridgeAddress}
}

func (l *EVMListener) ListenToEvents(startBlock *big.Int, chainID uint8, kvrw blockstore.KeyValueWriter, stopChn <-chan struct{}, errChn chan<- error) <-chan *relayer.Message {
	// TODO: This channel should be closed somewhere!
	ch := make(chan *relayer.Message)
	go func() {
		for {
			select {
			case <-stopChn:
				return
			default:
				// Although L1 block confirmations are checked in the data-transport-layer,
				// this check is needed as to not infinitely loop without a bound if a continue statement below is hit. Our bound is the latest Optimism batch index
				// NOTE: If we wanted to do our own check we would most likely need another separate sync service for l1 which seems unnecessary
				head, err := l.chainReader.LatestBlock()
				//log.Debug().Msgf("head in listener: %v", head)
				if err != nil {
					log.Error().Err(err).Msg("Unable to get latest block")
					time.Sleep(BlockRetryInterval)
					continue
				}
				// Sleep if the difference is less than BlockDelay; (latest - current) < BlockDelay
				if big.NewInt(0).Sub(head, startBlock).Cmp(BlockDelay) == -1 {
					time.Sleep(BlockRetryInterval)
					continue
				}

				if verified, err := l.chainReader.IsRollupVerified(startBlock.Uint64()); err != nil {
					log.Error().Err(err).Msg("Error while checking whether chain is verified")
					log.Error().Msgf("Block Number: %v", startBlock)
					time.Sleep(BlockRetryInterval)
					continue
				} else if !verified {
					log.Error().Msg("Chain is not verified at current index")
					log.Error().Msgf("Block Number: %v", startBlock)
					time.Sleep(BlockRetryInterval)
					continue
				}
				log.Error().Msgf("Verified index: %v", startBlock)

				logs, err := l.chainReader.FetchDepositLogs(context.Background(), l.bridgeAddress, startBlock, startBlock)
				if err != nil {
					// Filtering logs error really can appear only on wrong configuration or temporary network problem
					// so i do no see any reason to break execution
					log.Error().Err(err).Str("ChainID", string(chainID)).Msgf("Unable to filter logs")
					continue
				}

				log.Debug().Msgf("Rollup is verified on L1, can handle deposit event on optimism")
				for _, eventLog := range logs {
					m, err := l.eventHandler.HandleEvent(chainID, eventLog.DestinationChainID, eventLog.DepositNonce, eventLog.ResourceID)
					if err != nil {
						errChn <- err
						log.Error().Err(err)
						return
					}
					log.Debug().Msgf("Resolved message %+v in block %s", m, startBlock.String())
					ch <- m
				}
				if startBlock.Int64()%20 == 0 {
					// Logging process every 20 bocks to exclude spam
					log.Debug().Str("block", startBlock.String()).Uint8("chainID", chainID).Msg("Queried block for deposit events")
				}
				// TODO: We can store blocks to DB inside listener or make listener send something to channel each block to save it.
				//Write to block store. Not a critical operation, no need to retry
				err = blockstore.StoreBlock(kvrw, startBlock, chainID)
				if err != nil {
					log.Error().Str("block", startBlock.String()).Err(err).Msg("Failed to write latest block to blockstore")
				}
				// Goto next block
				startBlock.Add(startBlock, big.NewInt(1))
			}
		}
	}()
	return ch
}
