package main

func firstDigit(opcode uint16) uint16 {
	return opcode >> 12
}

func first2Digits(opcode uint16) uint16 {
	return opcode >> 8
}

func last3Digits(opcode uint16) uint16 {
	return opcode & 0x0fff
}

func last2Digits(opcode uint16) uint16 {
	return opcode & 0x00ff
}

func secondDigit(opcode uint16) uint16 {
	// this works by getting the second character then shifting it right twice
	return (opcode & 0x0f00) >> 8
}
