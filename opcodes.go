package main

type Opcode string

const (
	CLS   Opcode = "Clear Screen"
	JP    Opcode = "Jump"
	LD_Vx Opcode = "Set Vx = kk"
	ADD   Opcode = "Add Vx += kk"
	LD_I  Opcode = "Set register I = nnn"
	DRW   Opcode = "Draw"
	ERROR Opcode = "Error! Means the opcode is invalid"
)
