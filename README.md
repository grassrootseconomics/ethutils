# ethutils
[![Go test](https://github.com/grassrootseconomics/ethutils/actions/workflows/test.yaml/badge.svg)](https://github.com/grassrootseconomics/ethutils/actions/workflows/test.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/grassrootseconomics/ethutils.svg)](https://pkg.go.dev/github.com/grassrootseconomics/ethutils)


High level Go wrapper for Ethereum, primarily targetting [Celo](https://celo.org).

Features:

* Publishing smart contracts
* Gas transfers
* Smart contract interactions
* RPC error messages
* Address conversion/parsing utilities
* CICRegistry loader
* ABI parser (event and method signatures)
* Dump reverted tx reason
* CELO/ERC20 Balances scanner

## Installation

```bash
$ go get github.com/grassrootseconomics/ethutils
```

## License

[LGPL-3.0](COPYING.LESSER)