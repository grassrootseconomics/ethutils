package ethutils

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lmittmann/w3"
)

func Test_encodeAddress(t *testing.T) {
	address := "0x5523058cdFfe5F3c1EaDADD5015E55C6E00fb439"
	expected := "0000000000000000000000005523058cdffe5f3c1eadadd5015e55c6e00fb439"
	actual := encodeAddress(w3.A(address))
	if hex.EncodeToString(actual) != expected {
		t.Errorf("encodeAddress(%q) = %q; want %q", address, actual, expected)
	}
}

func Test_encodeAddressArray(t *testing.T) {
	addressArray := []common.Address{}
	expected := "00000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000000"
	actual := encodeAddressArray(addressArray)
	if hex.EncodeToString(actual) != expected {
		t.Errorf("encodeAddress(%q) = %q; want %q", addressArray, actual, expected)
	}
}

func Test_encodeAddressArrayEmpty(t *testing.T) {
	addressArray := []common.Address{}
	expected := encodeAddressArray(addressArray)
	actual := encodeAddressArrayEmpty()
	if !bytes.Equal(expected, actual) {
		t.Errorf("byte slices not equal\nexpected: %x\nactual:   %x", expected, actual)
	}
}

func Test_getReferalTag(t *testing.T) {
	p := &Provider{
		DivviConsumerAddress: w3.A("0x5523058cdFfe5F3c1EaDADD5015E55C6E00fb439"),
	}
	user := "0x079A79890B33b4574Ac7084683Cc38A5A10DA0Dd"
	expected := "6decb85d010080000000000000000000000000079a79890b33b4574ac7084683cc38a5a10da0dd0000000000000000000000005523058cdffe5f3c1eadadd5015e55c6e00fb43900000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000000"
	actual := p.GetReferalTag(w3.A(user))
	t.Logf("payload=%x", actual)
	if hex.EncodeToString(actual) != expected {
		t.Errorf("getReferalTag(%q, %q) = %q; want %q", user, p.DivviConsumerAddress, actual, expected)
	}
}
