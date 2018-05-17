package mos6502

import "testing"

func TestAddressWithBytes(t *testing.T) {
	type addressWithBytesTests struct {
		address Address
		high    byte
		low     byte
	}

	var tests = map[string]addressWithBytesTests{
		"basic": {0x1234, 0x12, 0x34},
	}

	for k, tt := range tests {
		t.Run(k, func(t *testing.T) {
			a := AddressFromBytes(tt.high, tt.low)
			expectUint16(t, uint16(tt.address), uint16(a))
		})
	}
}

func TestWithOffset(t *testing.T) {
	type addressWithBytesTests struct {
		expected Address
		original Address
		offset   byte
	}

	var tests = map[string]addressWithBytesTests{
		"basic":    {0x1111, 0x1100, 0x11},
		"negative": {0x1111, 0x1122, 0xEF},
	}

	for k, tt := range tests {
		t.Run(k, func(t *testing.T) {
			a := tt.original.WithOffset(tt.offset)
			expectUint16(t, uint16(tt.expected), uint16(a))
		})
	}
}
