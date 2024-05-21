package ethutils

import (
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"
)

// HexToAddress is similar to [common.HexToAddress] but converts to a checksumed address first.
// It is necessary to perform an address validation first before using HexToAddress e.g. with go-playground/validator.
func HexToAddress(hexAddress string) common.Address {
	return common.HexToAddress(ChecksumAddress(hexAddress))
}

// ChecksumAddress converts a mixed case hex address string to a checksumed one.
func ChecksumAddress(address string) string {
	address = strings.ToLower(address)
	address = strings.Replace(address, "0x", "", 1)

	sha := sha3.NewLegacyKeccak256()
	sha.Write([]byte(address))
	hash := sha.Sum(nil)
	hashstr := hex.EncodeToString(hash)
	result := []string{"0x"}
	for i, v := range address {
		res, _ := strconv.ParseInt(string(hashstr[i]), 16, 64)
		if res > 7 {
			result = append(result, strings.ToUpper(string(v)))
			continue
		}
		result = append(result, string(v))
	}

	return strings.Join(result, "")
}
