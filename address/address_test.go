package address

import (
	"testing"
)

func TestEVMAddressFromHex(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid with 0x prefix",
			input:   "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0",
			wantErr: false,
		},
		{
			name:    "valid without 0x prefix",
			input:   "742d35Cc6634C0532925a3b844Bc9e7595f0bEb0",
			wantErr: false,
		},
		{
			name:    "too short",
			input:   "0x742d35Cc",
			wantErr: true,
		},
		{
			name:    "invalid hex",
			input:   "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEbZ",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := EVMAddressFromHex(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("EVMAddressFromHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && addr.IsZero() {
				t.Error("expected non-zero address")
			}
		})
	}
}

func TestBech32RoundTrip(t *testing.T) {
	testAddr := EVMAddress{
		0x74, 0x2d, 0x35, 0xCc, 0x66, 0x34, 0xC0, 0x53,
		0x29, 0x25, 0xa3, 0xb8, 0x44, 0xBc, 0x9e, 0x75,
		0x95, 0xf0, 0xbE, 0xb0,
	}

	// Encode
	encoded, err := EncodeARCVAddress(testAddr, "arcv")
	if err != nil {
		t.Fatalf("EncodeARCVAddress() failed: %v", err)
	}

	if !startsWith(encoded, "arcv1") {
		t.Errorf("encoded address doesn't start with arcv1: %s", encoded)
	}

	// Decode
	decoded, err := DecodeARCVAddress(encoded, "arcv")
	if err != nil {
		t.Fatalf("DecodeARCVAddress() failed: %v", err)
	}

	// Compare
	if decoded != testAddr {
		t.Errorf("roundtrip failed: got %x, want %x", decoded, testAddr)
	}
}

func TestParseAddress(t *testing.T) {
	testAddr := EVMAddress{
		0x74, 0x2d, 0x35, 0xCc, 0x66, 0x34, 0xC0, 0x53,
		0x29, 0x25, 0xa3, 0xb8, 0x44, 0xBc, 0x9e, 0x75,
		0x95, 0xf0, 0xbE, 0xb0,
	}

	hexForm := testAddr.Hex()
	arcvForm, err := EncodeARCVAddress(testAddr, "arcv")
	if err != nil {
		t.Fatalf("failed to encode address: %v", err)
	}

	tests := []struct {
		name    string
		input   string
		want    EVMAddress
		wantErr bool
	}{
		{
			name:    "hex format",
			input:   hexForm,
			want:    testAddr,
			wantErr: false,
		},
		{
			name:    "arcv format",
			input:   arcvForm,
			want:    testAddr,
			wantErr: false,
		},
		{
			name:    "invalid format",
			input:   "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAddress(tt.input, "arcv")
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZeroAddress(t *testing.T) {
	zero := ZeroAddress()
	if !zero.IsZero() {
		t.Error("ZeroAddress() should return zero address")
	}

	// Test with actual zero address
	var zeroBytes EVMAddress
	if !zeroBytes.IsZero() {
		t.Error("empty EVMAddress should be zero")
	}

	// Test with non-zero
	nonZero := EVMAddress{0x01}
	if nonZero.IsZero() {
		t.Error("non-zero address reported as zero")
	}
}

func TestEdgeCases(t *testing.T) {
	// All zeros
	allZeros := ZeroAddress()
	encoded, err := EncodeARCVAddress(allZeros, "arcv")
	if err != nil {
		t.Fatalf("failed to encode zero address: %v", err)
	}
	decoded, err := DecodeARCVAddress(encoded, "arcv")
	if err != nil {
		t.Fatalf("failed to decode zero address: %v", err)
	}
	if decoded != allZeros {
		t.Error("zero address roundtrip failed")
	}

	// All 0xFF
	allOnes := EVMAddress{
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
		0xFF, 0xFF, 0xFF, 0xFF,
	}
	encoded, err = EncodeARCVAddress(allOnes, "arcv")
	if err != nil {
		t.Fatalf("failed to encode all-ones address: %v", err)
	}
	decoded, err = DecodeARCVAddress(encoded, "arcv")
	if err != nil {
		t.Fatalf("failed to decode all-ones address: %v", err)
	}
	if decoded != allOnes {
		t.Error("all-ones address roundtrip failed")
	}
}

func TestWrongHRP(t *testing.T) {
	testAddr := EVMAddress{0x01, 0x02, 0x03}
	
	// Encode with "arcv"
	encoded, err := EncodeARCVAddress(testAddr, "arcv")
	if err != nil {
		t.Fatalf("encode failed: %v", err)
	}

	// Try to decode with wrong HRP
	_, err = DecodeARCVAddress(encoded, "cosmos")
	if err == nil {
		t.Error("expected error when decoding with wrong HRP")
	}
}

func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

