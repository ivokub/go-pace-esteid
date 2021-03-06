package gopace

import (
	"bytes"
	"testing"
)

// canReal is the CAN of an actual smart card used for tests
var canReal []byte = []byte("050746")

func TestPad(t *testing.T) {
	cases := []struct {
		input    []byte
		expected []byte
	}{
		{
			input:    []byte{0x50, 0x00},
			expected: []byte{0x50, 0x00, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
	}
	sc := new(SecureCard)
	for _, tc := range cases {
		padded := sc.PadData(tc.input)
		if !bytes.Equal(padded, tc.expected) {
			t.Errorf("Expected: %x, got %x", tc.expected, padded)
		}
	}
}

func TestEncData(t *testing.T) {
	kenc := []byte{0x18, 0x45, 0x21, 0x62, 0xcc, 0x45, 0x46, 0x15, 0xf6, 0x88, 0x1d, 0xb6, 0x9a, 0xa1, 0xb3, 0x33, 0x5e, 0x87, 0x43, 0xd7, 0x87, 0x19, 0x85, 0xa3, 0x1c, 0xc7, 0xdb, 0x80, 0x4b, 0xc9, 0xfd, 0xf3}
	data := []byte{0x50, 0x00}
	expected := []byte{0x90, 0x68, 0xdb, 0x9e, 0x71, 0x67, 0x66, 0x29, 0xb3, 0xfa, 0xa7, 0xb1, 0x26, 0x32, 0xc7, 0x30}
	sc := &SecureCard{kenc: kenc, ssc: 1}
	enced, err := sc.EncData(data)
	if err != nil {
		t.Errorf("Can not encrypt: %v", err)
	}
	if !bytes.Equal(enced, expected) {
		t.Errorf("Expected: %x, got %x", expected, enced)
	}
}

func TestPrepare(t *testing.T) {
	kenc := []byte{0x18, 0x45, 0x21, 0x62, 0xcc, 0x45, 0x46, 0x15, 0xf6, 0x88, 0x1d, 0xb6, 0x9a, 0xa1, 0xb3, 0x33, 0x5e, 0x87, 0x43, 0xd7, 0x87, 0x19, 0x85, 0xa3, 0x1c, 0xc7, 0xdb, 0x80, 0x4b, 0xc9, 0xfd, 0xf3}
	kmac := []byte{0xc6, 0x8b, 0xc4, 0xe8, 0x5e, 0x0e, 0x8f, 0x16, 0x86, 0x70, 0xc9, 0x56, 0x56, 0x3e, 0x6c, 0x9b, 0x3d, 0xe3, 0x8a, 0xf8, 0x22, 0x89, 0x4f, 0x35, 0x47, 0xbc, 0x3c, 0x0d, 0x6f, 0x04, 0x7f, 0x0b}

	cases := []struct {
		expected []byte
		ssc      uint64
		header   []byte
		data     []byte
		le       []byte
	}{
		{
			[]byte{0x0c, 0xa4, 0x01, 0x0c, 0x1d, 0x87, 0x11, 0x01, 0x90, 0x68, 0xdb, 0x9e, 0x71, 0x67, 0x66, 0x29, 0xb3, 0xfa, 0xa7, 0xb1, 0x26, 0x32, 0xc7, 0x30, 0x8e, 0x08, 0x80, 0xd7, 0xb8, 0x50, 0x09, 0xf5, 0xe6, 0xc0, 0x00},
			1,
			[]byte{0x00, 0xa4, 0x01, 0x0c},
			[]byte{0x50, 0x00},
			nil,
		},
		{
			[]byte{0x0c, 0xb0, 0x00, 0x00, 0x0d, 0x97, 0x01, 0x00, 0x8e, 0x08, 0xcc, 0x31, 0xbd, 0x80, 0xe0, 0xf5, 0x3b, 0x2a, 0x00},
			5,
			[]byte{0x00, 0xB0, 0x00, 0x00},
			nil,
			[]byte{0x00},
		},
	}
	for _, tc := range cases {
		sc := &SecureCard{kenc: kenc, kmac: kmac, ssc: tc.ssc}
		apdu, err := sc.Prepare(tc.header, tc.data, tc.le)
		if err != nil {
			t.Fatalf("Prepare: %v", err)
		}
		if !bytes.Equal(apdu, tc.expected) {
			t.Errorf("Expected: %x got %x", tc.expected, apdu)
		}
	}
}
