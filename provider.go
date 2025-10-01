package ethutils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lmittmann/w3"
)

type (
	Option func(c *Provider)

	Provider struct {
		Client                        *w3.Client
		Signer                        types.Signer
		BalanceScannerContractAddress string
	}
)

func WithSigner(signer types.Signer) Option {
	return func(p *Provider) {
		p.Signer = signer
	}
}

func WithClient(w3Client *w3.Client) Option {
	return func(p *Provider) {
		p.Client = w3Client
	}
}

func WithBalanceScannerAddress(address string) Option {
	return func(p *Provider) {
		p.BalanceScannerContractAddress = address
	}
}

func NewProvider(url string, chainID int64, opts ...Option) *Provider {
	defaultProvider := &Provider{
		Client: w3.MustDial(url),
		Signer: types.LatestSignerForChainID(big.NewInt(chainID)),
	}

	for _, opt := range opts {
		opt(defaultProvider)
	}

	return defaultProvider
}
