package mos6502

func (b *byte) hasBit(bit uint8) bool {
	if bit > 7 {
		return false
	}
	return (*b & (1 << bit)) != 0
}

func (b *byte) withBit(bit uint8) byte {
	return *b + (byte(1 << bit))
}

/*
Core of the mos6502 processor.
*/
type Core struct {

	/* Program Counter */
	PC Address

	/* Accumulator */
	AC byte

	/* X Register */
	X byte

	/* Y Register */
	Y byte

	/* Stack Pointer */
	SP byte

	// Status flags; bomined make up the SR register.
	Negative  bool // bit 7
	Overflow  bool // bit 6
	Break     bool // bit 4
	Decimal   bool // bit 3
	Interrupt bool // bit 2
	Zero      bool // bit 1
	Carry     bool // bit 0

	// The connection to get memory stored values from.
	Bus Bus

	opCycles uint8
	op       chan int
}

/*
Tick the processor once, causing operatons to be performed.
*/
func (c *Core) Tick() {
	if c.opCycles > 0 {
		c.opCycles--
		return
	}

	// Wait for previous operation to finish:
	<-c.op

	// READ in Op.
	// c.Bus.ReadOp()
	readCh := make(chan Operation)
	go c.ReadInstruction(readCh)

	// If chan exists, recieve from it.
	go c.readyExecution(c.op, readCh)
}

/*
ReadInstruction reads in the next op.
*/
func (c *Core) ReadInstruction(ch chan Operation) {
	ch <- Operation{Code: 0x69}
}

/*
Execute an operation, sending it to a channel.
*/
func (c *Core) readyExecution(ch chan int, readCh chan Operation) {
	op := <-readCh
	c.Execute(op)
	ch <- 1
}

func (c *Core) Execute(op Operation) {
	v := c.Value(op)
	switch {
	case (op.Code & 0xE3) == 0x61:
		c.ADC(v)
	}
}

/*
IndirectAddress locates the proper address on the memory bus given the addresses's address.
*/
func (c *Core) IndirectAddress(start Address) Address {
	return AddressFromBytes(c.Bus.Content(start+1), c.Bus.Content(start))
}

/*
Address returns the address of the data used by the operation.
*/
func (c *Core) Address(op Operation) Address {
	switch op.Addressing() {
	case abs:
		return op.Full()
	case absX:
		return op.Full() + Address(c.X)
	case absY:
		return op.Full() + Address(c.Y)
	case ind:
		return c.IndirectAddress(op.Full())
	case xInd:
		return c.IndirectAddress(Address(op.Byte1 + c.X))
	case indY:
		return c.IndirectAddress(Address(op.Byte1)) + Address(c.Y)
	case rel:
		return c.PC.WithOffset(op.Byte1)
	case zpg:
		return Address(op.Byte1)
	case zpgX:
		return Address(op.Byte1) + Address(c.X)
	case zpgY:
		return Address(op.Byte1) + Address(c.Y)
	}

	return 0
}

/*
Value returns the value to be used by the operation. This depends on the addressing type.
*/
func (c *Core) Value(op Operation) *byte {
	switch op.Addressing() {
	case a:
		return &c.AC
	case imm:
		return &(op.Byte1)
	case impl:
		return 0
	default:
		return c.Bus.Content(c.Address(op))
	}
}

func (c *Core) ADC(v *byte) {
	a = &v
	likeSignedN := a > 0x7F && c.AC > 0x7F
	likeSignedP := a < 0x80 && c.AC < 0x80
	carry := 0
	if c.Carry {
		carry = 1
	}
	c.setACandFlags(c.AC + a + carry)

	// Overflow: If two negatives are over FF or two positives over F3
	c.Overflow = (likeSignedN && c.AC < 0x80) || (likeSignedP && c.AC > 0x7F)
}

func (c *Core) branch(v *byte) bool {
	o := c.PC
	c.PC = c.PC + Address(*v)
	return (o & 0xFF00) != (c.PC & 0xFF00)
}

func (c *Core) setACAndFlags(ac *byte) {
	oldAC := c.AC
	c.AC = ac

	// Flags
	// Zero: True if the value is zero.
	c.Zero = c.AC == 0x00

	// Negative: True if the first bit is set.
	c.Negative = c.AC > 0x7F

	// Carry: True if the value is greater than a byte.
	c.Carry = c.AC < oldAC
}

func (c *Core) AND(v *byte) {
	c.setACAndFlags(c.AC & *v)
}

func (c *Core) ASL(v *byte) {
	c.Carry = *v > 0x79
	*v = *v << 1
}

func (c *Core) BCC(v *byte) (bool, bool) {
	if c.Carry {
		return false, false
	}
	return true, branch(v)
}

func (c *Core) BCS(v *byte) (bool, bool) {
	if !c.Carry {
		return false, false
	}
	return true, branch(v)
}

func (c *Core) BEQ(v *byte) (bool, bool) {
	if c.Zero {
		return false, false
	}
	return true, branch(v)
}

func (c *Core) BIT(v *byte) {
	r := c.AC & v

	c.Zero = r == 0x00
	c.Overflow = (r & 0x40) == 0x40
	c.Negative = (r & 0x80) == 0x80
}

func (c *Core) BMI(v *byte) (bool, bool) {
	if !c.Negative {
		return false, false
	}
	return true, branch(v)
}

func (c *Core) BNE(v *byte) (bool, bool) {
	if !c.Zero {
		return false, false
	}
	return true, branch(v)
}

func (c *Core) BPL(v *byte) (bool, bool) {
	if c.Negative {
		return false, false
	}
	return true, branch(v)
}

func (c *Core) BRK(v *byte) (bool, bool) {
	if c.Negative {
		return false, false
	}
	return true, branch(v)
}
