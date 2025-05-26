package ethutils

import (
	"context"
	"testing"

	"github.com/lmittmann/w3"
)

func TestIndexIter_Next(t *testing.T) {
	t.Skip("This test requires a live uncapped Celo node to run against. It is not suitable for CI.")
	p := NewProvider("https://forno.celo.org", CeloMainnet)

	ctx := context.Background()
	iter, err := p.NewBatchIterator(ctx, w3.A("0xe2CEf4000d6003958c891D251328850f84654eb9"))
	if err != nil {
		t.Fatal(err)
	}

	for {
		batch, err := iter.Next(ctx)
		if err != nil {
			t.Error(err)
		}
		if batch == nil {
			break
		}
		t.Logf("index batch size %d", len(batch))
	}
}
