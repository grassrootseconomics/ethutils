package ethutils

import (
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/lmittmann/w3"
)

// Throwaway test keys.
var (
	privateKey  = "c74ecb579f9822e196b1866fef65950f5f9b8ed128ca92260b0de3c4dca8d436"
	publicKey   = "0x3AA8028a5FD03a0D35C32347e746842689b30987"
	deadAddress = w3.A("0x000000000000000000000000000000000000dEaD")
)

func TestProvider_SignContractExecutionTx(t *testing.T) {
	p := NewProvider("https://alfajores-forno.celo-testnet.org", CeloAlfajores)

	privateKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		t.Fatalf("Failed to parse private key %v", err)
	}

	sampleFunc := w3.MustNewFunc("transfer(address to, uint256 amount)", "bool")
	input, err := sampleFunc.EncodeArgs(deadAddress, w3.I("1"))
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		privateKey *ecdsa.PrivateKey
		txData     ContractExecutionTxOpts
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		wantType uint8
	}{
		{
			name: "Sign ERC20 transfer",
			args: args{
				privateKey: privateKey,
				txData: ContractExecutionTxOpts{
					// cUSD
					ContractAddress: w3.A("0x874069Fa1Eb16D44d622F2e0Ca25eeA172369bC1"),
					InputData:       input,
					GasLimit:        250000,
					GasFeeCap:       SafeGasFeeCap,
					GasTipCap:       SafeGasTipCap,
					Nonce:           0,
				},
			},
			wantErr:  false,
			wantType: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := p.SignContractExecutionTx(tt.args.privateKey, tt.args.txData)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.SignContractExecutionTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tx.Type() != tt.wantType {
				t.Errorf("Provider.SignContractExecutionTx() want type = %d, got %d", tt.wantType, tx.Type())
			}
		})
	}
}

func TestProvider_SignGasTransferTx(t *testing.T) {
	p := NewProvider("https://alfajores-forno.celo-testnet.org", CeloAlfajores)

	privateKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		t.Fatalf("Failed to parse private key %v", err)
	}

	type args struct {
		privateKey *ecdsa.PrivateKey
		txData     GasTransferTxOpts
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		wantType uint8
	}{
		{
			name: "Sign gas transfer",
			args: args{
				privateKey: privateKey,
				txData: GasTransferTxOpts{
					To:        deadAddress,
					Value:     big.NewInt(1),
					GasFeeCap: SafeGasFeeCap,
					GasTipCap: SafeGasTipCap,
					Nonce:     0,
				},
			},
			wantErr:  false,
			wantType: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := p.SignGasTransferTx(tt.args.privateKey, tt.args.txData)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.SignGasTransferTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tx.Type() != tt.wantType {
				t.Errorf("Provider.SignGasTransferTx() want type = %d, got %d", tt.wantType, tx.Type())
			}
		})
	}
}

func TestProvider_SignContractPublishTx(t *testing.T) {
	p := NewProvider("https://alfajores-forno.celo-testnet.org", CeloAlfajores)

	privateKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		t.Fatalf("Failed to parse private key %v", err)
	}

	contractBytecode, err := PrepareContractBytecodeData(
		"6080604052348015600f57600080fd5b5060405160ee38038060ee8339818101604052810190602d91906069565b505060a2565b600080fd5b6000819050919050565b6049816038565b8114605357600080fd5b50565b6000815190506063816042565b92915050565b60008060408385031215607d57607c6033565b5b60006089858286016056565b92505060206098858286016056565b9150509250929050565b603f8060af6000396000f3fe6080604052600080fdfea26469706673582212207b3c6ed6af876cf2dec1553939d6f993c60f20d2cd459bc2f3e6aed8f87fa31a64736f6c63430008180033",
		`[{"inputs": [{"internalType": "uint256","name": "a","type": "uint256"	},{	"internalType": "uint256","name": "b","type": "uint256"}],"stateMutability": "nonpayable","type": "constructor"}]`,
		[]any{
			w3.Big1,
			w3.Big2,
		},
	)
	if err != nil {
		t.Fatalf("Failed to prepare contract bytecode data %v", err)
	}

	type args struct {
		privateKey *ecdsa.PrivateKey
		txData     ContractPublishTxOpts
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		wantType uint8
	}{
		{
			name: "Sign contract publish",
			args: args{
				privateKey: privateKey,
				txData: ContractPublishTxOpts{
					ContractByteCode: contractBytecode,
					GasFeeCap:        SafeGasFeeCap,
					GasTipCap:        SafeGasTipCap,
					GasLimit:         1_000_000,
					Nonce:            0,
				},
			},
			wantErr:  false,
			wantType: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := p.SignContractPublishTx(tt.args.privateKey, tt.args.txData)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.SignContractPublishTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tx.Type() != tt.wantType {
				t.Errorf("Provider.SignContractPublishTx() want type = %d, got %d", tt.wantType, tx.Type())
			}
		})
	}
}
