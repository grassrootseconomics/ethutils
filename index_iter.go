package ethutils

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lmittmann/w3/module/eth"
	"github.com/lmittmann/w3/w3types"
)

type BatchIterator struct {
	provider     *Provider
	index        common.Address
	entryCount   *big.Int
	currentIndex int64
	batchSize    int64
}

const defaultBatchSize = 500

func (p *Provider) NewBatchIterator(ctx context.Context, index common.Address) (*BatchIterator, error) {
	var entryCount big.Int

	if err := p.Client.CallCtx(
		ctx,
		eth.CallFunc(index, entryCountFunc).Returns(&entryCount),
	); err != nil {
		return nil, err
	}

	return &BatchIterator{
		provider:     p,
		index:        index,
		entryCount:   &entryCount,
		currentIndex: 0,
		batchSize:    defaultBatchSize,
	}, nil
}

func (iter *BatchIterator) Next(ctx context.Context) ([]common.Address, error) {
	if iter.currentIndex >= iter.entryCount.Int64() {
		return nil, nil
	}

	endIndex := iter.currentIndex + iter.batchSize
	if endIndex > iter.entryCount.Int64() {
		endIndex = iter.entryCount.Int64()
	}

	batchSize := endIndex - iter.currentIndex
	calls := make([]w3types.RPCCaller, batchSize)
	tokenAddresses := make([]common.Address, batchSize)

	for i := int64(0); i < batchSize; i++ {
		index := iter.currentIndex + i
		calls[i] = eth.CallFunc(iter.index, entrySig, new(big.Int).SetInt64(index)).Returns(&tokenAddresses[i])
	}

	if err := iter.provider.Client.CallCtx(ctx, calls...); err != nil {
		return nil, err
	}

	iter.currentIndex = endIndex
	return tokenAddresses, nil
}
