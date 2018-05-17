package mos6502

/*
Address is a dual byte format used for memory locations.
*/
type Address uint16

/*
AddressFromBytes converts two bits into an address.
*/
func AddressFromBytes(high byte, low byte) Address {
	return (Address(high) << 8) + Address(low)
}

/*
WithOffset returns the address modified by the given offset. If over 0x7F, the number is a negative.
*/
func (a Address) WithOffset(o byte) Address {
	// If the left most digit is 1, get the 2's compliment and subtract.
	if o >= 0x80 {
		return a - Address((o^0xFF)+1)
	}

	return a + Address(o)
}
