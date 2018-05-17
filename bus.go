package mos6502

/*
Bus struct which can be used to emulate a memory bus.
*/
type Bus struct {
	// Not sure yet.
	data map[Address]byte
}

/*
Content returns the value at a specific address on the bus.
*/
func (b *Bus) Content(a Address) *byte {
	return *(*b.data[a])
}

/*
SetContent sets the content at a specific address on the bus.
*/
func (b *Bus) SetContent(a Address, d byte) {
	b.data[a] = d
}
