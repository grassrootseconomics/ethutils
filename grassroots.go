package ethutils

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lmittmann/w3"
	"github.com/lmittmann/w3/module/eth"
	"github.com/lmittmann/w3/w3types"
)

type (
	GESmartContracts map[string]bool
	GERegistry       map[string]common.Address
)

const (
	AccountIndex    = "AccountIndex"
	CustodialProxy  = "CustodialRegistrationProxy"
	GasFaucet       = "GasFaucet"
	TrainingVoucher = "TrainingVoucher"
	TokenIndex      = "TokenIndex"
	PoolIndex       = "PoolIndex"
)

var (
	tokenRegistryGetter = w3.MustNewFunc("tokenRegistry()", "address")
	quoterGetter        = w3.MustNewFunc("quoter()", "address")
	entryCountFunc      = w3.MustNewFunc("entryCount()", "uint256")
	entrySig            = w3.MustNewFunc("entry(uint256 _idx)", "address")

	ZeroAddress           = w3.A("0x0000000000000000000000000000000000000000")
	SarafuNetworkRegistry = w3.A("0xd1FB944748aca327a1ba036B082993D9dd9Bfa0C")
	CustodialRegistry     = w3.A("0x0cc9f4fff962def35bb34a53691180b13e653030")
)

// RegistryMap auto-loads all well known custodial related smart contract addresses through the CICRegistry pointer contract.
func (p *Provider) RegistryMap(ctx context.Context, registryAddress common.Address) (GERegistry, error) {
	var (
		listFunc = w3.MustNewFunc("addressOf(bytes32)", "address")

		accountIndexAddress    common.Address
		custodialProxy         common.Address
		gasFaucetAddress       common.Address
		tokenIndexAddress      common.Address
		trainingVoucherAddress common.Address
		poolIndex              common.Address
	)

	err := p.Client.CallCtx(
		ctx,
		eth.CallFunc(registryAddress, listFunc, common.BytesToHash(common.RightPadBytes([]byte(AccountIndex), 32))).Returns(&accountIndexAddress),
		eth.CallFunc(registryAddress, listFunc, common.BytesToHash(common.RightPadBytes([]byte(CustodialProxy), 32))).Returns(&custodialProxy),
		eth.CallFunc(registryAddress, listFunc, common.BytesToHash(common.RightPadBytes([]byte(GasFaucet), 32))).Returns(&gasFaucetAddress),
		eth.CallFunc(registryAddress, listFunc, common.BytesToHash(common.RightPadBytes([]byte(TokenIndex), 32))).Returns(&tokenIndexAddress),
		eth.CallFunc(registryAddress, listFunc, common.BytesToHash(common.RightPadBytes([]byte(TrainingVoucher), 32))).Returns(&trainingVoucherAddress),
		eth.CallFunc(registryAddress, listFunc, common.BytesToHash(common.RightPadBytes([]byte(PoolIndex), 32))).Returns(&poolIndex),
	)
	if err != nil {
		return nil, err
	}

	registryMap := make(GERegistry)
	registryMap[CustodialProxy] = custodialProxy
	registryMap[AccountIndex] = accountIndexAddress
	registryMap[GasFaucet] = gasFaucetAddress
	registryMap[TokenIndex] = tokenIndexAddress
	registryMap[TrainingVoucher] = trainingVoucherAddress
	registryMap[PoolIndex] = poolIndex

	return registryMap, nil
}

func (p *Provider) GetGESmartContracts(ctx context.Context, registries []string) (GESmartContracts, error) {
	geSmartContracts := make(GESmartContracts)

	for _, registry := range registries {
		registryMap, err := p.RegistryMap(ctx, w3.A(registry))
		if err != nil {
			return nil, err
		}

		for _, address := range registryMap {
			geSmartContracts[address.Hex()] = false
		}

		if tokenIndex := registryMap[TokenIndex]; tokenIndex != ZeroAddress {
			tokens, err := p.getAllTokensFromTokenIndex(ctx, tokenIndex)
			if err != nil {
				return nil, err
			}

			geSmartContracts[tokenIndex.Hex()] = true
			for _, token := range tokens {
				geSmartContracts[token.Hex()] = false
			}
		}

		if poolIndex := registryMap[PoolIndex]; poolIndex != ZeroAddress {
			pools, err := p.getAllTokensFromTokenIndex(ctx, poolIndex)
			if err != nil {
				return nil, err
			}

			geSmartContracts[poolIndex.Hex()] = true
			for _, pool := range pools {
				geSmartContracts[pool.Hex()] = false

				var poolTokenRegistry, priceQuoter common.Address
				err := p.Client.CallCtx(
					ctx,
					eth.CallFunc(pool, tokenRegistryGetter).Returns(&poolTokenRegistry),
					eth.CallFunc(pool, quoterGetter).Returns(&priceQuoter),
				)
				if err != nil {
					return nil, err
				}
				geSmartContracts[priceQuoter.Hex()] = false

				poolTokens, err := p.getAllTokensFromTokenIndex(ctx, poolTokenRegistry)
				if err != nil {
					return nil, err
				}
				geSmartContracts[poolTokenRegistry.Hex()] = true
				for _, token := range poolTokens {
					geSmartContracts[token.Hex()] = false
				}
			}
		}
	}

	return geSmartContracts, nil
}

func (p *Provider) getAllTokensFromTokenIndex(ctx context.Context, tokenIndex common.Address) ([]common.Address, error) {
	var tokenIndexEntryCount big.Int

	if err := p.Client.CallCtx(
		ctx,
		eth.CallFunc(tokenIndex, entryCountFunc).Returns(&tokenIndexEntryCount),
	); err != nil {
		return nil, err
	}

	calls := make([]w3types.RPCCaller, tokenIndexEntryCount.Int64())
	tokenAddresses := make([]common.Address, tokenIndexEntryCount.Int64())

	for i := 0; i < int(tokenIndexEntryCount.Int64()); i++ {
		calls[i] = eth.CallFunc(tokenIndex, entrySig, new(big.Int).SetInt64(int64(i))).Returns(&tokenAddresses[i])
	}

	if err := p.Client.CallCtx(ctx, calls...); err != nil {
		return nil, err
	}

	return tokenAddresses, nil
}
