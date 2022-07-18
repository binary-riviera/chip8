package main

func firstDigit(opcode uint16) uint16 {
	return opcode >> 12
}

func first2Digits(opcode uint16) uint16 {
	return opcode >> 8
}

func lastDigit(opcode uint16) uint16 {
	return opcode & 0x000f
}

func last3Digits(opcode uint16) uint16 {
	return opcode & 0x0fff
}

func last2Digits(opcode uint16) uint16 {
	return opcode & 0x00ff
}

func thirdDigit(opcode uint16) uint16 {
	return (opcode & 0x00f0) >> 4
}

func secondDigit(opcode uint16) uint16 {
	// this works by getting the second character then shifting it right twice
	return (opcode & 0x0f00) >> 8
}
