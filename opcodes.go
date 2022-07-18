package main

type Opcode string

const (
	CLS   Opcode = "Clear Screen"
	JP    Opcode = "Jump"
	LD    Opcode = "Set Vx = kk"
	ADD   Opcode = "Add Vx += kk"
	LD_I  Opcode = "Set register I = nnn"
	DRW   Opcode = "Draw"
	ERROR Opcode = "Error! Means the opcode is invalid"
	SET   Opcode = "Set Vx to value of Vy"
	OR    Opcode = "Binary OR"
	AND   Opcode = "Binary AND"
	XOR   Opcode = "Logical XOR"
	ADD_V Opcode = "Add Vx += Vy"
	SUB_V Opcode = "Sub"
	SHFT  Opcode = "Shift"
	JP_O  Opcode = "Jump with offset"
	RND   Opcode = "Random"
	SRD   Opcode = "Subroutines"
	SKP   Opcode = "Skip"
	SKPIF Opcode = "Skip if key"
	TIMER Opcode = "Timer"
	ADD_I Opcode = "Add Vx to I"
	KEY   Opcode = "Get key"
	FC    Opcode = "Font code"
	BCDC  Opcode = "Binary-coded decimal conversion"
	STORE Opcode = "Store mem"
	LOAD  Opcode = "Load mem"
)
