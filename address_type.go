package mos6502

/*
AddressType used to describe how an operation access a value.
*/
type AddressType uint8

/*
The types of addressing used by the operations.
*/
const (
	a AddressType = iota + 1
	abs
	absX
	absY
	impl
	imm
	ind
	xInd
	indY
	rel
	zpg
	zpgX
	zpgY
)

/*
Cycles returns the default number of cycles for an address type.
*/
func (t AddressType) Cycles() (int8, bool, bool) {
	switch t {
	// case A, Immediate, Implied:
	case a, imm, impl:
		return 2, false, false
	case rel:
		return 2, false, true
	case zpg:
		return 3, false, false
	case ind:
		return 5, false, false
	case indY:
		return 5, true, false
	case xInd:
		return 6, false, false
	case absY:
		return 4, true, false
	default:
		return 4, false, false
	}
}
