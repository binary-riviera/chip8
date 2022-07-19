package main

import (
	"fmt"
	"strconv"

	"github.com/veandco/go-sdl2/sdl"
)

type chip8 struct {
	memory      [4096]byte                         // Chip8 has 4k of memory
	V           [16]byte                           // CPU Registers V0-VE
	display     [WINDOW_WIDTH * WINDOW_HEIGHT]byte // graphics
	ins         uint16                             // the current opcode/instruction (2 bytes long)
	I           uint16                             // index register
	pc          uint16                             // program counter
	delay_timer byte                               //
	sound_timer byte                               //
	stack       [16]uint16                         //
	sp          uint16                             //
	key         [16]byte                           // current keypad state
	verbose     bool                               // whether to run the emulator in verbose mode
	last_ins    uint16                             // the instruction from the last cycle
	sdl_window  *sdl.Window                        // the SDL window to write to
}

func (c *chip8) initialise(window *sdl.Window, verbose bool) {
	c.verbose = verbose
	c.pc = 0x200 // program counter starts at 0x200
	c.ins = 0
	c.I = 0
	c.sp = 0
	c.delay_timer = 0
	c.sound_timer = 0
	c.sdl_window = window
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
	c.loadFont()
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

func (c *chip8) loadFont() {
	// starts at 0x050, ends at 0x09F
	for i := 0; i < 80; i++ {
		c.memory[i+0x050] = FONT[i]
	}
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
	for i := 0; i < WINDOW_WIDTH*WINDOW_HEIGHT; i++ {
		row := i / WINDOW_WIDTH
		col := i % WINDOW_WIDTH
		rect := sdl.Rect{X: int32(col) * WINDOW_SCALE, Y: int32(row) * WINDOW_SCALE, W: 1 * WINDOW_SCALE, H: 1 * WINDOW_SCALE}
		colour := 0xffffffff
		if c.display[i] == 1 {
			colour = 0x00000000
		}
		surface.FillRect(&rect, uint32(colour))
	}
	c.sdl_window.UpdateSurface()
}

func (c *chip8) draw(xV uint16, yV uint16, n uint16) {
	// first we need to get the X and Y coordinates from VX and VY
	// We need to account for wrapping by modding the width of the screen
	x := c.V[xV] % WINDOW_WIDTH
	old_x := x
	y := c.V[yV] % WINDOW_HEIGHT
	if c.verbose {
		fmt.Println("x: " + strconv.Itoa(int(x)) + ", y: " + strconv.Itoa(int(y)))
	}
	// next we need to set VF to 0 for later
	c.V[0xF] = 0
	// now the hard part...
	for i := 0; i < int(n); i++ {
		x = old_x
		// represents a new row
		sprite := c.memory[int(c.I)+i]
		// we need to iterate through every bit in the sprite
		for mask := byte(0x80); mask != 0; mask >>= 1 {
			if (sprite & mask) != 0 { // it's 1
				// so we need to get the xy coord from c.display to compare
				// it should be something like Y * width + X + j
				idx := int(y)*WINDOW_WIDTH + int(x) // j
				if c.display[idx] == 1 {
					c.display[idx] = 0
					c.V[0xF] = 1
				} else {
					c.display[idx] = 1
				}
			} else { // it's 0

			}
			x++
		}
		y++
	}
	c.showDisplay()
}

func (c *chip8) execute(opcode Opcode, ins uint16) {
	if opcode == CLS {
		c.clearDisplay()
		c.showDisplay()
	} else if opcode == DRW {
		c.draw(secondDigit(ins), thirdDigit(ins), lastDigit(ins))
	} else if opcode == JP {
		c.pc = last3Digits(ins) // set PC to NNN
	} else if opcode == LD_I {
		c.I = last3Digits(ins) // set I to NNN
	} else if opcode == LD {
		c.V[secondDigit(ins)] = byte(last2Digits(ins)) // set V[x] to NN
	} else if opcode == ADD {
		c.V[secondDigit(ins)] += byte(last2Digits(ins)) // Add Vx += kk
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
	} else if firstDigit(ins) == 0x2 || ins == 0x00EE {
		return SRD
	} else if firstDigit(ins) == 0x3 || firstDigit(ins) == 0x4 || firstDigit(ins) == 0x5 || firstDigit(ins) == 0x9 {
		return SKP
	} else if firstDigit(ins) == 0x6 {
		return LD
	} else if firstDigit(ins) == 0x7 {
		return ADD
	} else if firstDigit(ins) == 0x8 {
		if lastDigit(ins) == 0 {
			return SET
		} else if lastDigit(ins) == 0x1 {
			return OR
		} else if lastDigit(ins) == 0x2 {
			return AND
		} else if lastDigit(ins) == 0x3 {
			return XOR
		} else if lastDigit(ins) == 0x4 {
			return ADD_V
		} else if lastDigit(ins) == 0x5 || lastDigit(ins) == 0x7 {
			return SUB_V
		} else if lastDigit(ins) == 0x6 || lastDigit(ins) == 0xE {
			return SHFT
		} else {
			return ERROR
		}
	} else if firstDigit(ins) == 0xA {
		return LD_I
	} else if firstDigit(ins) == 0xB {
		return JP_O
	} else if firstDigit(ins) == 0xC {
		return RND
	} else if firstDigit(ins) == 0xD {
		return DRW
	} else if firstDigit(ins) == 0xE {
		return SKPIF
	} else if firstDigit(ins) == 0xF && last2Digits(ins) == 0x1E {
		return ADD_I
	} else if firstDigit(ins) == 0xF && last2Digits(ins) == 0x0A {
		return KEY
	} else if firstDigit(ins) == 0xF && last2Digits(ins) == 0x29 {
		return FC
	} else if firstDigit(ins) == 0xF && last2Digits(ins) == 0x33 {
		return BCDC
	} else if firstDigit(ins) == 0xF && last2Digits(ins) == 0x55 {
		return STORE
	} else if firstDigit(ins) == 0xF && last2Digits(ins) == 0x65 {
		return LOAD
	} else if firstDigit(ins) == 0xF {
		return TIMER
	} else {
		return ERROR
	}
}
