package ethutils

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lmittmann/w3"
)

func TestProvider_TokensBalance(t *testing.T) {
	p := NewProvider("https://forno.celo.org", CeloMainnet)

	tokens := []common.Address{
		w3.A("0x02cc0715E844a45bA56Ad391D92DCd6537315177"),
		w3.A("0x2105a206B7bec31E2F90acF7385cc8F7F5f9D273"),
		w3.A("0x45d747172e77d55575c197CbA9451bC2CD8F4958"),
	}

	balances, err := p.TokensBalance(context.Background(), w3.A("0x0030cfF17fAf04a4Bb0657d47999099B3cbF9ccc"), tokens)
	if err != nil {
		t.Error(err)
	}
	t.Log(balances)
}
