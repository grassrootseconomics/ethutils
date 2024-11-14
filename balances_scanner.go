package ethutils

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lmittmann/w3"
	"github.com/lmittmann/w3/module/eth"
)

type (
	TokenBalances map[common.Address]*big.Int

	tokensBalanceResp struct {
		Success bool
		Data    []byte
	}
)

const (
	BalanceScannerContractAddress = "0xF62107c53a5b18646E823a21ed531ED934B1CE9E"
)

func (p *Provider) TokensBalance(ctx context.Context, owner common.Address, tokenAddresses []common.Address) (TokenBalances, error) {
	var (
		resp []tokensBalanceResp

		tokensBalanceFunc = w3.MustNewFunc("tokensBalance(address,address[])", "(bool success, bytes data)[]")
	)

	err := p.Client.CallCtx(
		ctx,
		eth.CallFunc(
			w3.A(BalanceScannerContractAddress),
			tokensBalanceFunc,
			owner,
			tokenAddresses,
		).Returns(&resp),
	)
	if err != nil {
		return nil, err
	}

	tokenBalances := make(TokenBalances)
	for i, v := range resp {
		tokenBalances[tokenAddresses[i]] = new(big.Int).SetBytes(v.Data)
	}

	return tokenBalances, nil
}
