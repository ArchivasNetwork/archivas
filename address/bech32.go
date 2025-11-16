package address

// Minimal Bech32 implementation
// Based on BIP 173: https://github.com/bitcoin/bips/blob/master/bip-0173.mediawiki

import (
	"fmt"
	"strings"
)

const charset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

var gen = []int{0x3b6a57b2, 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3}

// bech32Polymod computes the Bech32 checksum polymod
func bech32Polymod(values []int) int {
	chk := 1
	for _, v := range values {
		top := chk >> 25
		chk = (chk&0x1ffffff)<<5 ^ v
		for i := 0; i < 5; i++ {
			if (top>>uint(i))&1 == 1 {
				chk ^= gen[i]
			}
		}
	}
	return chk
}

// bech32HrpExpand expands the human-readable part
func bech32HrpExpand(hrp string) []int {
	ret := make([]int, 0, len(hrp)*2+1)
	for _, c := range hrp {
		ret = append(ret, int(c)>>5)
	}
	ret = append(ret, 0)
	for _, c := range hrp {
		ret = append(ret, int(c)&31)
	}
	return ret
}

// bech32VerifyChecksum verifies the checksum
func bech32VerifyChecksum(hrp string, data []int) bool {
	return bech32Polymod(append(bech32HrpExpand(hrp), data...)) == 1
}

// bech32CreateChecksum creates a checksum
func bech32CreateChecksum(hrp string, data []int) []int {
	values := append(bech32HrpExpand(hrp), data...)
	values = append(values, []int{0, 0, 0, 0, 0, 0}...)
	polymod := bech32Polymod(values) ^ 1
	ret := make([]int, 6)
	for i := 0; i < 6; i++ {
		ret[i] = (polymod >> uint(5*(5-i))) & 31
	}
	return ret
}

// bech32Encode encodes hrp and data into a bech32 string
func bech32Encode(hrp string, data []int) (string, error) {
	if len(hrp) < 1 || len(hrp) > 83 {
		return "", fmt.Errorf("invalid hrp length")
	}
	for _, c := range hrp {
		if c < 33 || c > 126 {
			return "", fmt.Errorf("invalid hrp character")
		}
	}
	
	combined := append(data, bech32CreateChecksum(hrp, data)...)
	ret := hrp + "1"
	for _, d := range combined {
		if d < 0 || d >= len(charset) {
			return "", fmt.Errorf("invalid data value")
		}
		ret += string(charset[d])
	}
	return ret, nil
}

// bech32Decode decodes a bech32 string
func bech32Decode(bech string) (string, []int, error) {
	if len(bech) > 90 {
		return "", nil, fmt.Errorf("bech32 string too long")
	}
	
	lower := strings.ToLower(bech)
	upper := strings.ToUpper(bech)
	if bech != lower && bech != upper {
		return "", nil, fmt.Errorf("mixed case")
	}
	bech = lower
	
	pos := strings.LastIndex(bech, "1")
	if pos < 1 || pos+7 > len(bech) || len(bech)-1-pos > 90 {
		return "", nil, fmt.Errorf("invalid separator position")
	}
	
	hrp := bech[:pos]
	data := make([]int, 0, len(bech)-pos-1)
	for _, c := range bech[pos+1:] {
		d := strings.IndexRune(charset, c)
		if d == -1 {
			return "", nil, fmt.Errorf("invalid character")
		}
		data = append(data, d)
	}
	
	if !bech32VerifyChecksum(hrp, data) {
		return "", nil, fmt.Errorf("invalid checksum")
	}
	
	return hrp, data[:len(data)-6], nil
}

// convertBits converts data between bit groups
func convertBits(data []byte, fromBits, toBits uint, pad bool) ([]int, error) {
	acc := 0
	bits := uint(0)
	ret := make([]int, 0, (len(data)*int(fromBits)+int(toBits)-1)/int(toBits))
	maxv := (1 << toBits) - 1
	
	for _, value := range data {
		acc = (acc << fromBits) | int(value)
		bits += fromBits
		for bits >= toBits {
			bits -= toBits
			ret = append(ret, (acc>>bits)&maxv)
		}
	}
	
	if pad {
		if bits > 0 {
			ret = append(ret, (acc<<(toBits-bits))&maxv)
		}
	} else if bits >= fromBits || ((acc<<(toBits-bits))&maxv) != 0 {
		return nil, fmt.Errorf("invalid padding")
	}
	
	return ret, nil
}

// ConvertBits8To5 converts 8-bit data to 5-bit groups (for encoding)
func ConvertBits8To5(data []byte) ([]int, error) {
	return convertBits(data, 8, 5, true)
}

// ConvertBits5To8 converts 5-bit data to 8-bit bytes (for decoding)
func ConvertBits5To8(data []int) ([]byte, error) {
	bits, err := convertBits5To8Bytes(data)
	if err != nil {
		return nil, err
	}
	return bits, nil
}

func convertBits5To8Bytes(data []int) ([]byte, error) {
	acc := 0
	bits := uint(0)
	ret := make([]byte, 0, (len(data)*5)/8)
	maxv := 255
	
	for _, value := range data {
		acc = (acc << 5) | value
		bits += 5
		for bits >= 8 {
			bits -= 8
			ret = append(ret, byte((acc>>bits)&maxv))
		}
	}
	
	if bits >= 5 || ((acc<<(8-bits))&maxv) != 0 {
		return nil, fmt.Errorf("invalid padding")
	}
	
	return ret, nil
}

