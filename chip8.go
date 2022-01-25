package main

type chip8 struct {
	memory      [4096]byte    // Chip8 has 4k of memory
	V           [16]byte      // CPU Registers V0-VE
	display     [64 * 32]byte // graphics
	opcode      uint16        // the current opcode (2 bytes long)
	I           uint16        // index register
	pc          uint16        // program counter
	delay_timer byte          //
	sound_timer byte          //
	stack       [16]uint16    //
	sp          uint16        //
	key         [16]byte      // current keypad state
}

func (c *chip8) initialise() {
	// TODO: is there a way to set default values for the struct and then just reset it?
	c.pc = 0x200 // program counter starts at 0x200
	c.opcode = 0
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

func (c *chip8) loadGame(bytes []byte) {
	size := len(bytes)
	println(size)
	for i := 0; i < size; i++ {
		c.memory[i+0x200] = bytes[i]
	}
}

func (c *chip8) emulateCycle() {

}

func (c *chip8) setKeys() {

}
