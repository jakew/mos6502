package mos6502

import (
	"fmt"
	"testing"
)

/*
Test that we can print out the opration as a string.
*/
func TestString(t *testing.T) {
	op := Operation{Code: 0x69}
	expectString(t, "69 00 --", op.String())

	op.Byte1 = 0xBE
	op.Byte2 = 0xEF

	expectString(t, "69 BE --", op.String())
}

func TestFull(t *testing.T) {
	type fullOp struct {
		full  Address
		Byte1 byte
		Byte2 byte
	}

	var tests = map[string]fullOp{
		"basic": {0x1234, 0x12, 0x34},
	}

	for k, tt := range tests {
		t.Run(k, func(t *testing.T) {
			op := Operation{Code: 0x61}
			op.Byte1, op.Byte2 = tt.Byte1, tt.Byte2
			expectUint16(t, uint16(tt.full), uint16(op.Full()))
		})
	}
}

func TestAddressTypeCycles(t *testing.T) {
	type addressTypeCyclesTest struct {
		addressType AddressType
		cycles      int8
		paged       bool
		branch      bool
	}

	var tests = map[string]addressTypeCyclesTest{
		"accumulator": {a, 2, false, false},
		"absolute":    {abs, 4, false, false},
		"absolute,X":  {absX, 4, false, false},
		"absolute,Y":  {absY, 4, true, false},
		"immediate":   {imm, 2, false, false},
		"implied":     {impl, 2, false, false},
		"indirect":    {ind, 5, false, false},
		"indirect,X":  {xInd, 6, false, false},
		"indirect,Y":  {indY, 5, true, false},
		"relative":    {rel, 2, false, true},
		"zeropage":    {zpg, 3, false, false},
		"zeropage,X":  {zpgX, 4, false, false},
		"zeropage,Y":  {zpgY, 4, false, false},
	}

	for k, tt := range tests {
		t.Run(k, func(t *testing.T) {
			c, p, b := tt.addressType.Cycles()
			expectInt8(t, tt.cycles, c)
			expectBool(t, tt.paged, p)
			expectBool(t, tt.branch, b)
		})
	}
}

func TestOperationCycles(t *testing.T) {
	type operationCyclesTest struct {
		operation Operation
		cycles    int8
		paged     bool
		branch    bool
	}

	var tests = map[string]operationCyclesTest{
		"BRK":       {Operation{Code: 0x00}, 7, false, false},
		"ORA X,ind": {Operation{Code: 0x01}, 6, false, false},
		"ORA zpg":   {Operation{Code: 0x05}, 3, false, false},
		"ASL zpg":   {Operation{Code: 0x06}, 5, false, false},
		"PHP impl":  {Operation{Code: 0x08}, 3, false, false},
		"ORA #":     {Operation{Code: 0x09}, 2, false, false},
		"ASL A":     {Operation{Code: 0x0A}, 2, false, false},
		"ORA abs":   {Operation{Code: 0x0D}, 4, false, false},
		"ASL abs":   {Operation{Code: 0x0E}, 6, false, false},
		"BPL rel":   {Operation{Code: 0x10}, 2, false, true},
		"ORA ind,Y": {Operation{Code: 0x11}, 5, true, false},
		"ORA zpg,X": {Operation{Code: 0x15}, 4, false, false},
		"ASL zpg,X": {Operation{Code: 0x16}, 6, false, false},
		"CLC impl":  {Operation{Code: 0x18}, 2, false, false},
		"ORA abs,Y": {Operation{Code: 0x19}, 4, true, false},
		"ORA abs,X": {Operation{Code: 0x1D}, 4, true, false},
		"ASL abs,X": {Operation{Code: 0x1E}, 7, false, false},
		"JSR abs":   {Operation{Code: 0x20}, 6, false, false},
		"AND X,ind": {Operation{Code: 0x21}, 6, false, false},
		"BIT zpg":   {Operation{Code: 0x24}, 3, false, false},
		"AND zpg":   {Operation{Code: 0x25}, 3, false, false},
		"ROL zpg":   {Operation{Code: 0x26}, 5, false, false},
		"PLP impl":  {Operation{Code: 0x28}, 4, false, false},
		"AND #":     {Operation{Code: 0x29}, 2, false, false},
		"ROL A":     {Operation{Code: 0x2A}, 2, false, false},
		"BIT abs":   {Operation{Code: 0x2C}, 4, false, false},
		"AND abs":   {Operation{Code: 0x2D}, 4, false, false},
		"ROL abs":   {Operation{Code: 0x2E}, 6, false, false},
		"BMI rel":   {Operation{Code: 0x30}, 2, false, true},
		"AND ind,Y": {Operation{Code: 0x31}, 5, true, false},
		"AND zpg,X": {Operation{Code: 0x35}, 4, false, false},
		"ROL zpg,X": {Operation{Code: 0x36}, 6, false, false},
		"SEC impl":  {Operation{Code: 0x38}, 2, false, false},
		"AND abs,Y": {Operation{Code: 0x39}, 4, true, false},
		"AND abs,X": {Operation{Code: 0x3D}, 4, true, false},
		"ROL abs,x": {Operation{Code: 0x3E}, 7, false, false},
		"RTI impl":  {Operation{Code: 0x40}, 6, false, false},
		"EOR X,ind": {Operation{Code: 0x41}, 6, false, false},
		"EOR zpg":   {Operation{Code: 0x45}, 3, false, false},
		"LSR zpg":   {Operation{Code: 0x46}, 5, false, false},
		"PHA impl":  {Operation{Code: 0x48}, 3, false, false},
		"EOR #":     {Operation{Code: 0x49}, 2, false, false},
		"LSR A":     {Operation{Code: 0x4A}, 2, false, false},
		"JMP abs":   {Operation{Code: 0x4C}, 3, false, false},
		"EOR abs":   {Operation{Code: 0x4D}, 4, false, false},
		"LSR abs":   {Operation{Code: 0x4E}, 6, false, false},
		"BVC rel":   {Operation{Code: 0x50}, 2, false, true},
		"EOR ind,Y": {Operation{Code: 0x51}, 5, true, false},
		"EOR zpg,X": {Operation{Code: 0x55}, 4, false, false},
		"LSR zpg,X": {Operation{Code: 0x56}, 6, false, false},
		"CLI impl":  {Operation{Code: 0x58}, 2, false, false},
		"EOR abs,Y": {Operation{Code: 0x59}, 4, true, false},
		"EOR abs,X": {Operation{Code: 0x5D}, 4, true, false},
		"LSR abs,X": {Operation{Code: 0x5E}, 7, false, false},
		"RTS impl":  {Operation{Code: 0x60}, 6, false, false},
		"ABC X,ind": {Operation{Code: 0x61}, 6, false, false},
		"ADC zpg":   {Operation{Code: 0x65}, 3, false, false},
		"ROR zpg":   {Operation{Code: 0x66}, 5, false, false},
		"PLA impl":  {Operation{Code: 0x68}, 4, false, false},
		"ADC #":     {Operation{Code: 0x69}, 2, false, false},
		"ROR A":     {Operation{Code: 0x6A}, 2, false, false},
		"JMP ind":   {Operation{Code: 0x6C}, 5, false, false},
		"ADC abs":   {Operation{Code: 0x6D}, 4, false, false},
		"ROR abs":   {Operation{Code: 0x6E}, 6, false, false},
		"BVS rel":   {Operation{Code: 0x70}, 2, false, true},
		"ADC ind,Y": {Operation{Code: 0x71}, 5, true, false},
		"ADC zpg,X": {Operation{Code: 0x75}, 4, false, false},
		"ROR zpg,X": {Operation{Code: 0x76}, 6, false, false},
		"SEI impl":  {Operation{Code: 0x78}, 2, false, false},
		"ADC abs,Y": {Operation{Code: 0x79}, 4, true, false},
		"ADC abs,X": {Operation{Code: 0x7D}, 4, true, false},
		"ROR abs,X": {Operation{Code: 0x7E}, 7, false, false},
		"STA X,ind": {Operation{Code: 0x81}, 6, false, false},
		"STY zpg":   {Operation{Code: 0x84}, 3, false, false},
		"STA zpg":   {Operation{Code: 0x85}, 3, false, false},
		"STX zpg":   {Operation{Code: 0x86}, 3, false, false},
		"DEY impl":  {Operation{Code: 0x88}, 2, false, false},
		"TXA impl":  {Operation{Code: 0x8A}, 2, false, false},
		"STY abs":   {Operation{Code: 0x8C}, 4, false, false},
		"STA abs":   {Operation{Code: 0x8D}, 4, false, false},
		"STX abs":   {Operation{Code: 0x8E}, 4, false, false},
		"BCC rel":   {Operation{Code: 0x90}, 2, false, true},
		"STA ind,Y": {Operation{Code: 0x91}, 6, false, false},
		"STY zpg,X": {Operation{Code: 0x94}, 4, false, false},
		"STA zpg,X": {Operation{Code: 0x95}, 4, false, false},
		"STX zpg,Y": {Operation{Code: 0x96}, 4, false, false},
		"TYA impl":  {Operation{Code: 0x98}, 2, false, false},
		"STA abs,Y": {Operation{Code: 0x99}, 5, false, false},
		"TXS impl":  {Operation{Code: 0x9A}, 2, false, false},
		"STA abs,X": {Operation{Code: 0x9D}, 5, false, false},
		"LDY #":     {Operation{Code: 0xA0}, 2, false, false},
		"LDA X,ind": {Operation{Code: 0xA1}, 6, false, false},
		"LDX #":     {Operation{Code: 0xA2}, 2, false, false},
		"LDY zpg":   {Operation{Code: 0xA4}, 3, false, false},
		"LDA zpg":   {Operation{Code: 0xA5}, 3, false, false},
		"LDX zpg":   {Operation{Code: 0xA6}, 3, false, false},
		"TAY impl":  {Operation{Code: 0xA8}, 2, false, false},
		"LDA #":     {Operation{Code: 0xA9}, 2, false, false},
		"TAX impl":  {Operation{Code: 0xAA}, 2, false, false},
		"LDY abs":   {Operation{Code: 0xAC}, 4, false, false},
		"LDA abs":   {Operation{Code: 0xAD}, 4, false, false},
		"LDX abs":   {Operation{Code: 0xAE}, 4, false, false},
		"BCS rel":   {Operation{Code: 0xB0}, 2, false, true},
		"LDA ind,Y": {Operation{Code: 0xB1}, 5, true, false},
		"LDY zpg,X": {Operation{Code: 0xB4}, 4, false, false},
		"LDA zpg,X": {Operation{Code: 0xB5}, 4, false, false},
		"LDX zpg,Y": {Operation{Code: 0xB6}, 4, false, false},
		"CLV impl":  {Operation{Code: 0xB8}, 2, false, false},
		"LDA abs,Y": {Operation{Code: 0xB9}, 4, true, false},
		"TSX impl":  {Operation{Code: 0xBA}, 2, false, false},
		"LDY abs,X": {Operation{Code: 0xBC}, 4, true, false},
		"LDA abs,X": {Operation{Code: 0xBD}, 4, true, false},
		"LDX abs,Y": {Operation{Code: 0xBE}, 4, true, false},
		"CPY #":     {Operation{Code: 0xC0}, 2, false, false},
		"CMP X,ind": {Operation{Code: 0xC1}, 6, false, false},
		"CPY zpg":   {Operation{Code: 0xC4}, 3, false, false},
		"CMP zpg":   {Operation{Code: 0xC5}, 3, false, false},
		"DEC zpg":   {Operation{Code: 0xC6}, 5, false, false},
		"INY impl":  {Operation{Code: 0xC8}, 2, false, false},
		"CMP #":     {Operation{Code: 0xC9}, 2, false, false},
		"DEX impl":  {Operation{Code: 0xCA}, 2, false, false},
		"CPY abs":   {Operation{Code: 0xCC}, 4, false, false},
		"CMP abs":   {Operation{Code: 0xCD}, 4, false, false},
		"DEC abs":   {Operation{Code: 0xCE}, 3, false, false},
		"BNE rel":   {Operation{Code: 0xD0}, 2, false, true},
		"CMP ind,Y": {Operation{Code: 0xD1}, 5, true, false},
		"CMP zpg,X": {Operation{Code: 0xD5}, 4, false, false},
		"DEC zpg,X": {Operation{Code: 0xD6}, 6, false, false},
		"CLD impl":  {Operation{Code: 0xD8}, 2, false, false},
		"CMP abs,Y": {Operation{Code: 0xD9}, 4, true, false},
		"CMP abs,X": {Operation{Code: 0xDD}, 4, true, false},
		"DEC abs,X": {Operation{Code: 0xDE}, 7, false, false},
		"CPX #":     {Operation{Code: 0xE0}, 2, false, false},
		"SPC X,ind": {Operation{Code: 0xE1}, 6, false, false},
		"CPX zpg":   {Operation{Code: 0xE4}, 3, false, false},
		"SPC zpg":   {Operation{Code: 0xE5}, 3, false, false},
		"INC zpg":   {Operation{Code: 0xE6}, 5, false, false},
		"INX impl":  {Operation{Code: 0xE8}, 2, false, false},
		"SBC #":     {Operation{Code: 0xE9}, 2, false, false},
		"NOP impl":  {Operation{Code: 0xEA}, 2, false, false},
		"CPX abs":   {Operation{Code: 0xEC}, 4, false, false},
		"SBC abs":   {Operation{Code: 0xED}, 4, false, false},
		"INC abs":   {Operation{Code: 0xEE}, 6, false, false},
		"BEQ rel":   {Operation{Code: 0xF0}, 2, false, true},
		"SBC ind,Y": {Operation{Code: 0xF1}, 5, true, false},
		"SBC zpg,X": {Operation{Code: 0xF5}, 4, false, false},
		"INC zpg,X": {Operation{Code: 0xF6}, 6, false, false},
		"SED impl":  {Operation{Code: 0xF8}, 2, false, false},
		"SBC abs,Y": {Operation{Code: 0xF9}, 4, true, false},
		"SBC abs,X": {Operation{Code: 0xFD}, 4, true, false},
		"INC abs,X": {Operation{Code: 0xFE}, 7, false, false},
	}

	for k, tt := range tests {
		title := fmt.Sprintf("%s (%02X)", k, tt.operation.Code)
		t.Run(title, func(t *testing.T) {
			c, p, b := tt.operation.Cycles()

			expectInt8(t, tt.cycles, c)
			expectBool(t, tt.paged, p)
			expectBool(t, tt.branch, b)
		})
	}
}
