package ethutils

import (
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

const (
	divviMagicPrefix               = "6decb85d"
	knownEmptyAddressArrayEncoding = "00000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000000"
)

var referalTagFormat1Byte = []byte{0x01}

func (p *Provider) GetReferalTag(user common.Address) []byte {
	encodedBytes := concatBytes(encodeAddress(user), encodeAddress(p.DivviConsumerAddress), encodeAddressArrayEmpty())
	payloadLen := []byte{byte(len(encodedBytes) >> 8), byte(len(encodedBytes))}
	magicPrefix, _ := hex.DecodeString(divviMagicPrefix)
	payload := concatBytes(magicPrefix, referalTagFormat1Byte, payloadLen, encodedBytes)
	return payload
}

func encodeAddress(address common.Address) []byte {
	return common.LeftPadBytes(address.Bytes(), 32)
}

func encodeAddressArray(addresses []common.Address) []byte {
	arrayDataOffset := big.NewInt(96)
	arrayDataOffsetBytes := common.LeftPadBytes(arrayDataOffset.Bytes(), 32)

	arrayLength := big.NewInt(int64(len(addresses)))
	arrayLengthBytes := common.LeftPadBytes(arrayLength.Bytes(), 32)

	var addressesBytes []byte
	for _, addr := range addresses {
		addressesBytes = append(addressesBytes, encodeAddress(addr)...)
	}

	result := make([]byte, 0, len(arrayDataOffsetBytes)+len(arrayLengthBytes)+len(addressesBytes))
	result = append(result, arrayDataOffsetBytes...)
	result = append(result, arrayLengthBytes...)
	result = append(result, addressesBytes...)

	return result
}

func encodeAddressArrayEmpty() []byte {
	emptyBytes, _ := hex.DecodeString(knownEmptyAddressArrayEncoding)
	return emptyBytes
}

func concatBytes(slices ...[]byte) []byte {
	var result []byte
	for _, slice := range slices {
		result = append(result, slice...)
	}

	return result
}
