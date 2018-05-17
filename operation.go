package mos6502

import "fmt"


/*
Operation is the instruction to the Core converted into a struct.
*/
type Operation struct {
	Code  byte
	Byte1 byte
	Byte2 byte
}

/*
Returns the Operation as a series of bytes in a string readble format.
*/
func (o Operation) String() string {
	switch o.Size() {
	case 3:
		return fmt.Sprintf("%02X %02X %02X", o.Code, o.Byte1, o.Byte2)
	case 2:
		return fmt.Sprintf("%02X %02X --", o.Code, o.Byte1)
	default:
		return fmt.Sprintf("%02X -- --", o.Code)
	}
}

/*
Clone the structure into a new copy of the operation.
*/
func (o Operation) Clone() Operation {
	return Operation{
		Code:  o.Code,
		Byte1: o.Byte1,
		Byte2: o.Byte2}
}

/*
Full returns the full two-byte content of the operation.
*/
func (o Operation) Full() Address {
	return AddressFromBytes(o.Byte1, o.Byte2)
}

/*
Addressing returns the addressing type of the operation.
*/
func (o Operation) Addressing() AddressType {
	switch {
	case o.Code == 0x6C:
		return ind
	case o.Code == 0xBE:
		return absY
	case o.Code == 0x96, o.Code == 0xB6:
		return zpgY
	case (o.Code & 0x9F) == 0xA:
		return a
	case (o.Code & 0x1C) == 0x0C,
		o.Code == 0x20:
		return abs
	case (o.Code & 0x1C) == 0x1C:
		return absX
	case (o.Code & 0x1F) == 0x19:
		return absY
	case (o.Code & 0x1F) == 0x09,
		(o.Code & 0x9D) == 0x80:
		return imm
	case (o.Code & 0x0F) == 0x08,
		(o.Code & 0x8F) == 0x8A,
		(o.Code & 0x9F) == 0x00:
		return impl
	case (o.Code & 0x1F) == 0x01:
		return xInd
	case (o.Code & 0x1F) == 0x11:
		return indY
	case (o.Code & 0x1F) == 0x10:
		return rel
	case (o.Code & 0x1C) == 0x04:
		return zpg
	case (o.Code & 0x1C) == 0x14:
		return zpgX
	default:
		panic("Invalid address used. No address type found.")
	}
}

/*
cyclesByPattern returns the base cycle count (c), if the cycle count is affected by pages (p) and if it is affected by
branches (b).
*/
func (o Operation) cyclesByPattern(c int8, p bool, b bool) (int8, bool, bool) {
	if (o.Code&0xC0) != 0x80 && (o.Code&0x07) == 0x06 {
		if (o.Code & 0x1C) == 0x1C {
			return c + 3, p, b
		}
		if o.Code == 0x91 {
			panic("This is right.")
		}
		return c + 2, p, b
	} else if (o.Code & 0x1C) == 0x1C {
		return c, true, b
	}

	return c, p, b
}

/*
cycleOverrides returns the exact cycle count (c), if the cycle count is affected by pages (p) and if it is affected by
branches (b) for operations that don't fall under the general rules.
*/
func (o Operation) cycleOverrides(c int8, p bool, b bool) (int8, bool, bool) {
	switch o.Code {
	case 0x08, 0x48, 0x9D, 0x91:
		return c + 1, false, false
	case 0x20, 0x28, 0x68:
		return c + 2, false, false
	case 0x40, 0x60:
		return c + 4, false, false
	case 0x4C:
		return 3, false, false
	case 0x99:
		return 5, false, false
	case 0xCE:
		return 3, false, false
	case 0x00:
		return c + 5, false, false
	default:
		return c, p, b
	}
}

// TODO: Refactor this.
func (o Operation) Cycles() (int8, bool, bool) {
	return o.cycleOverrides(o.cyclesByPattern(o.Addressing().Cycles()))
}

// TODO: Refactor this.
func (o Operation) Size() int8 {
	switch o.Code {
	case 0xB6, 0x61, 0x65, 0x69, 0x75, 0x90:
		return 2
	case 0x6C, 0x6D, 0x71, 0x79, 0x7D:
		return 3
	default:
		return 1
	}
}
