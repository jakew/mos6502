package mos6502

import "testing"

func expectString(t *testing.T, expected string, actual string) {
	if expected != actual {
		t.Logf("Expected \"%s\" but got \"%s\".", expected, actual)
		t.Fail()
	}
}

/*
Given an expectation of a byte and the actual byte, report back a failure if they are different, logging both.
*/
func expectByte(t *testing.T, expected byte, actual byte) {
	if expected != actual {
		t.Logf("Expected \"%02X\" but got \"%02X\".", expected, actual)
		t.Fail()
	}
}

/*
Given an expectation of a int8 and the actual int8, report back a failure if they are different, logging both.
*/
func expectInt8(t *testing.T, expected int8, actual int8) {
	if expected != actual {
		t.Logf("Expected \"%03d\" but got \"%03d\".", expected, actual)
		t.Fail()
	}
}

/*
Given an expectation of a uint8 and the actual uint8, report back a failure if they are different, logging both.
*/
func expectUint8(t *testing.T, expected uint8, actual uint8) {
	if expected != actual {
		t.Logf("Expected \"%03d\" but got \"%03d\".", expected, actual)
		t.Fail()
	}
}

/*
Given an expectation of a Address and the actual Address, report back a failure if they are different, logging both.
*/
func expectAddress(t *testing.T, expected Address, actual Address) {
	if expected != actual {
		t.Logf("Expected \"%04X\" but got \"%04X\".", expected, actual)
		t.Fail()
	}
}

/*
Given an expectation of a uint16 and the actual uint16, report back a failure if they are different, logging both.
*/
func expectUint16(t *testing.T, expected uint16, actual uint16) {
	if expected != actual {
		t.Logf("Expected \"%04X\" but got \"%04X\".", expected, actual)
		t.Fail()
	}
}

/*
Given an expectation of a bool and the actual bool, report back a failure if they are different, logging both.
*/
func expectBool(t *testing.T, expected bool, actual bool) {
	if expected != actual {
		t.Logf("Expected \"%t\" but got \"%t\".", expected, actual)
		t.Fail()
	}
}
func expectCore(t *testing.T, expected *Core, actual *Core) {
	fail := false

	if expected.PC != actual.PC {
		t.Logf("Expected PC to be \"%04X\" but got \"%04X\".", expected.PC, actual.PC)
		fail = true
	}

	if expected.AC != actual.AC {
		t.Logf("Expected AC to be \"%02X\" but got \"%02X\".", expected.AC, actual.AC)
		fail = true
	}

	if expected.X != actual.X {
		t.Logf("Expected X to be \"%02X\" but got \"%02X\".", expected.X, actual.X)
		fail = true
	}

	if expected.Y != actual.Y {
		t.Logf("Expected Y to be \"%02X\" but got \"%02X\".", expected.Y, actual.Y)
		fail = true
	}

	if expected.SP != actual.SP {
		t.Logf("Expected SP to be \"%02X\" but got \"%02X\".", expected.SP, actual.SP)
		fail = true
	}

	if expected.Negative != actual.Negative {
		t.Logf("Expected Negative to be \"%t\" but got \"%t\".", expected.Negative, actual.Negative)
		fail = true
	}

	if expected.Overflow != actual.Overflow {
		t.Logf("Expected Overflow to be \"%t\" but got \"%t\".", expected.Overflow, actual.Overflow)
		fail = true
	}

	if expected.Decimal != actual.Decimal {
		t.Logf("Expected Decimal to be \"%t\" but got \"%t\".", expected.Decimal, actual.Decimal)
		fail = true
	}

	if expected.Interrupt != actual.Interrupt {
		t.Logf("Expected Interrupt to be \"%t\" but got \"%t\".", expected.Interrupt, actual.Interrupt)
		fail = true
	}

	if expected.Zero != actual.Zero {
		t.Logf("Expected Zero to be \"%t\" but got \"%t\".", expected.Zero, actual.Zero)
		fail = true
	}

	if expected.Carry != actual.Carry {
		t.Logf("Expected Carry to be \"%t\" but got \"%t\".", expected.Carry, actual.Carry)
		fail = true
	}

	if fail {
		t.Fail()
	}
}
