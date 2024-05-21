package ethutils

import (
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestHexToAddress(t *testing.T) {
	type args struct {
		hexAddress string
	}
	tests := []struct {
		args args
		want common.Address
	}{
		{
			args: args{
				hexAddress: "0xd1FB944748aca327a1ba036B082993D9dd9Bfa0C",
			},
			want: SarafuNetworkRegistry,
		},
	}
	for _, tt := range tests {
		t.Run("Hex2Address", func(t *testing.T) {
			if got := HexToAddress(tt.args.hexAddress); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HexToAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChecksumAddress(t *testing.T) {
	type args struct {
		hexAddress string
	}
	tests := []struct {
		args args
		want common.Address
	}{
		{
			args: args{
				// mixedCase
				hexAddress: "0xd1fb944748aca327a1ba036B082993D9dd9Bfa0C",
			},
			want: SarafuNetworkRegistry,
		},
	}
	for _, tt := range tests {
		t.Run("ChecksumAddress", func(t *testing.T) {
			if got := HexToAddress(ChecksumAddress(tt.args.hexAddress)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChecksumAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
