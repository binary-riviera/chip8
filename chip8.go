package main

import (
	"fmt"
	"strconv"

	"github.com/veandco/go-sdl2/sdl"
)

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
	last_ins    uint16        // the instruction from the last cycle
	sdl_window  *sdl.Window   // the SDL window to write to
	sdl_scale   int32         // the scale for the SDL window
}

func (c *chip8) initialise(window *sdl.Window, scale int, verbose bool) {
	c.verbose = verbose
	c.pc = 0x200 // program counter starts at 0x200
	c.ins = 0
	c.I = 0
	c.sp = 0
	c.delay_timer = 0
	c.sound_timer = 0
	c.sdl_window = window
	c.sdl_scale = int32(scale)
	// clear and show display
	c.clearDisplay()
	c.showDisplay()
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
	// FETCH
	c.ins = uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1])

	// DECODE
	if c.verbose && c.ins != c.last_ins {
		fmt.Printf("Instruction: 0x%x\n", c.ins)

	}
	opcode := decodeOpcode(c.ins)
	if opcode == ERROR {
		panic("Error, invalid instruction!")
	}
	if c.verbose && c.ins != c.last_ins {
		fmt.Println(opcode)
	}

	// EXECUTE
	c.execute(opcode, c.ins)

	// RESET
	// update timer
	if opcode != JP { // If jumping PC shouldn't be increased
		c.pc += 2 // since an opcode is 2 bytes long
	}
	c.last_ins = c.ins // record the last instruction
}

func (c *chip8) loadGame(bytes []byte) {
	size := len(bytes)
	if c.verbose {
		fmt.Println("Loaded " + strconv.Itoa(size) + " lines")
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

func (c *chip8) showDisplay() {
	surface, err := c.sdl_window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0)
	rect := sdl.Rect{0, 0, 32 * c.sdl_scale, 16 * c.sdl_scale}
	surface.FillRect(&rect, 0xffff0000)
	c.sdl_window.UpdateSurface()
}

func (c *chip8) execute(opcode Opcode, ins uint16) {

	if opcode == CLS {
		c.clearDisplay()
	} else if opcode == JP {
		c.pc = last3Digits(ins) // set PC to NNN
	} else if opcode == LD_I {
		c.I = last3Digits(ins) // set I to NNN
	} else if opcode == LD_Vx {
		c.V[secondDigit(ins)] = byte(last2Digits(ins)) // set V[x] to NN
	} else {
		panic("ERROR! opcode " + opcode + " not implemented")
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
