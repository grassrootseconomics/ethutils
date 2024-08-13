package ethutils

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lmittmann/w3"
)

func TestProvider_SimulateRevertedTx(t *testing.T) {
	t.Skip("requires archive node")

	p := NewProvider("https://forno.celo.org", CeloMainnet)

	type args struct {
		txHash      common.Hash
		blockNumber *big.Int
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Reverted reason",
			args: args{
				txHash:      w3.H("0x30cc64b374656e517b161bd0d24d5bb401e9b018b30d864f437fa14ac879c1b0"),
				blockNumber: big.NewInt(24925269),
			},
			want:    "ERC20: transfer amount exceeds balance",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.SimulateRevertedTx(context.Background(), tt.args.txHash, tt.args.blockNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.SimulateRevertedTx() error = %v, wantErr %v, got = %s", err, tt.wantErr, got)
				return
			}
			if got != tt.want {
				t.Errorf("Provider.SimulateRevertedTx() = %v, want %v", got, tt.want)
			}
			t.Log(got)
		})
	}
}
