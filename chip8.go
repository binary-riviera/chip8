package main

import "fmt"

type chip8 struct {
	memory      [4096]byte    // Chip8 has 4k of memory
	V           [16]byte      // CPU Registers V0-VE
	display     [64 * 32]byte // graphics
	ins         uint16        // the current opcode/instruction (2 bytes long)
	I           uint16        // index register
	pc          uint16        // program counter
	delay_timer byte          //
	sound_timer byte          //
	stack       [16]uint16    //
	sp          uint16        //
	key         [16]byte      // current keypad state
	verbose     bool          // whether to run the emulator in verbose mode
}

func (c *chip8) initialise(verbose bool) {
	c.verbose = verbose
	c.pc = 0x200 // program counter starts at 0x200
	c.ins = 0
	c.I = 0
	c.sp = 0
	c.delay_timer = 0
	c.sound_timer = 0
	// clear display
	for i := 0; i < len(c.display); i++ {
		c.display[i] = 0
	}
	// clear stack
	for i := 0; i < len(c.stack); i++ {
		c.stack[i] = 0
	}
	// clear registers
	for i := 0; i < len(c.V); i++ {
		c.V[i] = 0
	}
	// clear memory
	for i := 0; i < len(c.memory); i++ {
		c.memory[i] = 0
	}
	// TODO: implement font loading
}

func (c *chip8) emulateCycle() {
	// fetch opcode
	c.ins = uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1])
	// decode opcode
	if c.verbose {
		fmt.Printf("Instruction: 0x%x\n", c.ins)
	}
	opcode := decodeOpcode(c.ins)
	if opcode == ERROR {
		panic("Error, invalid instruction!")
	}
	if c.verbose {
		fmt.Println(opcode)
	}

	// execute opcode
	c.execute(opcode, c.ins)

	// update timers

	if opcode != JP { // If jumping PC shouldn't be increased
		c.pc += 2 // since an opcode is 2 bytes long
	}
}

func (c *chip8) loadGame(bytes []byte) {
	size := len(bytes)
	if c.verbose {
		fmt.Println(size)
	}
	for i := 0; i < size; i++ {
		c.memory[i+0x200] = bytes[i]
	}
}

func (c *chip8) setKeys() {

}

func (c *chip8) clearDisplay() {
	for i := range c.display {
		c.display[i] = 0
	}
}

func (c *chip8) execute(opcode Opcode, ins uint16) {

	if opcode == DRW {
		c.clearDisplay()
	} else if opcode == JP {
		// set PC to NNN
		c.pc = last3Digits(ins)
	}

}

func decodeOpcode(ins uint16) Opcode {
	// should return some enum representing the code
	if ins == 0x00E0 {
		return CLS
	} else if firstDigit(ins) == 0x1 {
		return JP
	} else if firstDigit(ins) == 0x6 {
		return LD_Vx
	} else if firstDigit(ins) == 0x7 {
		return ADD
	} else if firstDigit(ins) == 0xA {
		return LD_I
	} else if firstDigit(ins) == 0xD {
		return DRW
	} else {
		return ERROR
	}
}

func firstDigit(opcode uint16) uint16 {
	return opcode >> 12
}

func first2Digits(opcode uint16) uint16 {
	return opcode >> 8
}

func last3Digits(opcode uint16) uint16 {
	return opcode & 0x0fff
}
