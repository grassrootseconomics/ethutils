package ethutils

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lmittmann/w3"
)

type (
	ContractExecutionTxOpts struct {
		ContractAddress common.Address
		InputData       []byte
		GasFeeCap       *big.Int
		GasTipCap       *big.Int
		GasLimit        uint64
		Nonce           uint64
	}

	GasTransferTxOpts struct {
		To        common.Address
		Value     *big.Int
		GasFeeCap *big.Int
		GasTipCap *big.Int
		Nonce     uint64
	}

	ContractPublishTxOpts struct {
		ContractByteCode []byte
		GasFeeCap        *big.Int
		GasTipCap        *big.Int
		GasLimit         uint64
		Nonce            uint64
	}
)

func (p *Provider) SignContractExecutionTx(privateKey *ecdsa.PrivateKey, txData ContractExecutionTxOpts) (*types.Transaction, error) {
	tx, err := types.SignNewTx(privateKey, p.Signer, &types.DynamicFeeTx{
		To:        &txData.ContractAddress,
		Nonce:     txData.Nonce,
		Data:      txData.InputData,
		Gas:       txData.GasLimit,
		GasFeeCap: txData.GasFeeCap,
		GasTipCap: txData.GasTipCap,
	})
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (p *Provider) SignGasTransferTx(privateKey *ecdsa.PrivateKey, txData GasTransferTxOpts) (*types.Transaction, error) {
	tx, err := types.SignNewTx(privateKey, p.Signer, &types.DynamicFeeTx{
		Value:     txData.Value,
		To:        &txData.To,
		Nonce:     txData.Nonce,
		Gas:       21000,
		GasFeeCap: txData.GasFeeCap,
		GasTipCap: txData.GasTipCap,
	})
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (p *Provider) SignContractPublishTx(privateKey *ecdsa.PrivateKey, txData ContractPublishTxOpts) (*types.Transaction, error) {
	tx, err := types.SignNewTx(privateKey, p.Signer, &types.DynamicFeeTx{
		Data:      txData.ContractByteCode,
		Nonce:     txData.Nonce,
		Gas:       txData.GasLimit,
		GasFeeCap: txData.GasFeeCap,
		GasTipCap: txData.GasTipCap,
	})
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (p *Provider) SignContractExecutionTxWithcKES(privateKey *ecdsa.PrivateKey, txData ContractExecutionTxOpts) (*types.Transaction, error) {
	cKES := w3.A("0x456a3D042C0DbD3db53D5489e98dFb038553B0d0")

	tx, err := types.SignNewTx(privateKey, p.Signer, &types.CeloDynamicFeeTxV2{
		Data:        txData.InputData,
		Nonce:       txData.Nonce,
		Gas:         txData.GasLimit,
		GasFeeCap:   txData.GasFeeCap,
		GasTipCap:   txData.GasTipCap,
		FeeCurrency: &cKES,
		To:          &txData.ContractAddress,
	})
	if err != nil {
		return nil, err
	}

	return tx, nil
}
