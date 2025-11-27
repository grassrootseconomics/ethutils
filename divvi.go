package ethutils

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
)

type SubmitBody struct {
	TxHash  string `json:"txHash"`
	ChainID int64  `json:"chainId"`
}

const (
	divviMagicPrefix               = "6decb85d"
	knownEmptyAddressArrayEncoding = "00000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000000000000000000"

	divviAPIEndpoint = "https://api.divvi.xyz/submitReferral"
	userAgent        = "ethutils/v1.x.x"
	contentType      = "application/json"
)

var (
	referalTagFormat1Byte = []byte{0x01}

	ErrDivviClientNotInitialized = errors.New("divvi client not initialized")
)

func (p *Provider) GetReferalTag(user common.Address) []byte {
	encodedBytes := ConcatBytes(encodeAddress(user), encodeAddress(p.DivviConsumerAddress), encodeAddressArrayEmpty())
	payloadLen := []byte{byte(len(encodedBytes) >> 8), byte(len(encodedBytes))}
	magicPrefix, _ := hex.DecodeString(divviMagicPrefix)
	payload := ConcatBytes(magicPrefix, referalTagFormat1Byte, payloadLen, encodedBytes)
	return payload
}

func (p *Provider) SubmitReferral(ctx context.Context, txHash common.Hash) error {
	if p.DivviClient == nil {
		return ErrDivviClientNotInitialized
	}
	body := SubmitBody{
		TxHash:  txHash.Hex(),
		ChainID: CeloMainnet,
	}

	b, err := json.Marshal(&body)
	if err != nil {
		return err
	}

	resp, err := p.requestWithCtx(ctx, http.MethodPost, divviAPIEndpoint, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	if err := parseResponse(resp); err != nil {
		return err
	}

	return nil
}

func (p *Provider) requestWithCtx(ctx context.Context, method string, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	return p.do(req)
}

func (p *Provider) do(req *http.Request) (*http.Response, error) {
	builtRequest, err := setHeaders(req)
	if err != nil {
		return nil, err
	}

	return p.DivviClient.Do(builtRequest)
}

func setHeaders(req *http.Request) (*http.Request, error) {
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", contentType)
	req.Header.Set("Content-Type", contentType)

	return req, nil
}

func parseResponse(resp *http.Response) error {
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("divvi api error: status=%s", resp.Status)
		}
		return fmt.Errorf("divvi api error: status=%s body=%s", resp.Status, string(b))
	}

	return nil
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

func ConcatBytes(slices ...[]byte) []byte {
	var result []byte
	for _, slice := range slices {
		result = append(result, slice...)
	}

	return result
}
