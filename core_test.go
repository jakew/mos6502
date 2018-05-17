package mos6502

import (
	"testing"
)

type addressTestCase struct {
	address   Address
	core      Core
	operation Operation
}

func TestAddress(t *testing.T) {
	bus := Bus{
		data: map[Address]byte{
			0x00B4: 0xEE,
			0x00B5: 0x12,
			0x00BA: 0x12,
			0x00BB: 0xEE,
			0x1F11: 0x34,
			0x1F12: 0x12,
		},
	}

	var tests = map[string]addressTestCase{
		/* Name,  Addres, Op,   Byt1, Byt2, PC,     AC,   X,    Y,  */
		"abs":   {address: 0x1234, operation: Operation{Code: 0x6D, Byte1: 0x12, Byte2: 0x34}, core: Core{PC: 0x0000, AC: 0x00, X: 0x00, Y: 0x00}},
		"abs,x": {address: 0x12FF, operation: Operation{Code: 0x7D, Byte1: 0x12, Byte2: 0x34}, core: Core{PC: 0x0000, AC: 0x00, X: 0xCB, Y: 0x00}},
		"abs,y": {address: 0x1301, operation: Operation{Code: 0x79, Byte1: 0x12, Byte2: 0x34}, core: Core{PC: 0x0000, AC: 0x00, X: 0x00, Y: 0xCD}},
		"ind":   {address: 0x1234, operation: Operation{Code: 0x6C, Byte1: 0x1F, Byte2: 0x11}, core: Core{PC: 0x0000, AC: 0x00, X: 0x00, Y: 0x00}},
		"x,ind": {address: 0xEE12, operation: Operation{Code: 0x61, Byte1: 0xB4, Byte2: 0x00}, core: Core{PC: 0x0000, AC: 0x00, X: 0x06, Y: 0x00}},
		"ind,y": {address: 0x12F4, operation: Operation{Code: 0x71, Byte1: 0xB4, Byte2: 0x00}, core: Core{PC: 0x0000, AC: 0x00, X: 0x00, Y: 0x06}},
		"rel+":  {address: 0x9A18, operation: Operation{Code: 0x90, Byte1: 0x7F, Byte2: 0x00}, core: Core{PC: 0x9999, AC: 0x00, X: 0x00, Y: 0x00}},
		"rel-":  {address: 0x9919, operation: Operation{Code: 0x90, Byte1: 0x80, Byte2: 0x00}, core: Core{PC: 0x9999, AC: 0x00, X: 0x00, Y: 0x00}},
		"zpg":   {address: 0x00AA, operation: Operation{Code: 0x65, Byte1: 0xAA, Byte2: 0x00}, core: Core{PC: 0x0000, AC: 0x00, X: 0x00, Y: 0x00}},
		"zpg,x": {address: 0x00BB, operation: Operation{Code: 0x75, Byte1: 0xB0, Byte2: 0x00}, core: Core{PC: 0x0000, AC: 0x00, X: 0x0B, Y: 0x00}},
		"zpg,y": {address: 0x00CC, operation: Operation{Code: 0xB6, Byte1: 0xC0, Byte2: 0x00}, core: Core{PC: 0x0000, AC: 0x00, X: 0x00, Y: 0x0C}},
	}

	for k, tt := range tests {
		t.Run(k, func(t *testing.T) {
			tt.core.Bus = bus
			address := tt.core.Address(tt.operation)
			expectAddress(t, tt.address, address)
		})
	}
}

func TestValue(t *testing.T) {

	var tests = map[string]struct {
		value byte
		op    Operation
		core  Core
		bus   Bus
	}{
		"absolute": {
			value: 0x12,
			op:    Operation{Code: 0x6D, Byte1: 0x12, Byte2: 0x34},
			core:  Core{},
			bus:   Bus{data: map[Address]byte{0x1234: 0x12}},
		},
		"absolute X": {
			value: 0x22,
			op:    Operation{Code: 0x7D, Byte1: 0x12, Byte2: 0x34},
			core:  Core{X: 0xFF},
			bus:   Bus{data: map[Address]byte{0x1333: 0x22}},
		},
		"absolute Y": {
			value: 0x22,
			op:    Operation{Code: 0x79, Byte1: 0x12, Byte2: 0x34},
			core:  Core{Y: 0xFE},
			bus:   Bus{data: map[Address]byte{0x1332: 0x22}},
		},
		"immediate": {
			value: 0x12,
			op:    Operation{Code: 0x69, Byte1: 0x12},
			core:  Core{},
		},
		"indirect X": {
			value: 0x32,
			op:    Operation{Code: 0x61, Byte1: 0xB4},
			core:  Core{X: 0x06},
			bus: Bus{data: map[Address]byte{
				0x00BA: 0x12,
				0x00BB: 0xEE,
				0xEE12: 0x32,
			}},
		},
		"indirect Y": {
			value: 0xDD,
			op:    Operation{Code: 0x71, Byte1: 0xB4},
			core:  Core{Y: 0x06},
			bus: Bus{data: map[Address]byte{
				0x00B4: 0xEE,
				0x00B5: 0x12,
				0x12F4: 0xDD,
			}},
		},
		"zeropage": {
			value: 0xBB,
			op:    Operation{Code: 0x65, Byte1: 0x23},
			core:  Core{},
			bus:   Bus{data: map[Address]byte{0x0023: 0xBB}},
		},
		"zeropage X": {
			value: 0xBB,
			op:    Operation{Code: 0x75, Byte1: 0xBB},
			core:  Core{X: 0x10},
			bus:   Bus{data: map[Address]byte{0x00CB: 0xBB}},
		},
		"zeropage Y": {
			value: 0xEE,
			op:    Operation{Code: 0xB6, Byte1: 0x22},
			core:  Core{Y: 0x11},
			bus:   Bus{data: map[Address]byte{0x0033: 0xEE}},
		},
	}

	for k, tt := range tests {
		t.Run(k, func(t *testing.T) {
			tt.core.Bus = tt.bus
			value := tt.core.Value(tt.op)
			expectByte(t, tt.value, value)
		})
	}
}

type opTest struct {
	op       Operation
	start    Core
	expected Core
}

func TestOPS(t *testing.T) {
	var tests = map[string]opTest{
		"0x69 basic": {
			op:       Operation{Code: 0x69, Byte1: 0x36},
			start:    Core{AC: 0xC8, Carry: true},
			expected: Core{AC: 0xFF, Negative: true},
		},
		"0x69 add to 128": {
			op:       Operation{Code: 0x69, Byte1: 0x40},
			start:    Core{AC: 0x40},
			expected: Core{AC: 0x80, Negative: true, Overflow: true},
		},
		"0x69 add to 0": {
			op:       Operation{Code: 0x69, Byte1: 0x37},
			start:    Core{AC: 0xC8, Carry: true},
			expected: Core{AC: 0x00, Carry: true, Zero: true},
		},
		"0x69 add to 1": {
			op:       Operation{Code: 0x69, Byte1: 0x38},
			start:    Core{AC: 0xC8, Carry: true},
			expected: Core{AC: 0x01, Carry: true},
		},
		"0x69 add -129": {
			op:       Operation{Code: 0x69, Byte1: 0xFF},
			start:    Core{AC: 0x80, Negative: true},
			expected: Core{AC: 0x7F, Carry: true, Overflow: true},
		},
	}

	for k, tt := range tests {
		t.Run(k, func(t *testing.T) {
			tt.start.Execute(tt.op)
			expectCore(t, &tt.expected, &tt.start)
		})
	}
}
