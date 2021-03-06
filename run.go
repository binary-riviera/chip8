// main entry point for program
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

func readROM(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Error while opening file", err)
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		log.Fatal("Error while getting file stats", err)
	}
	// calculate the bytes size of the file
	var size int64 = info.Size()
	bytes := make([]byte, size)
	// read into buffer
	buffer := bufio.NewReader(file)
	_, err = buffer.Read(bytes)
	return bytes
}

func main() {
	verbosePtr := flag.Bool("verbose", false, "verbose mode")

	flag.Parse()
	fmt.Println("Loading SDL...")
	window := setupSDL()
	fmt.Println("Starting chip-8 emulator")
	c := chip8{}
	c.initialise(window, *verbosePtr)
	fmt.Println("Initialised emulator...")
	bytes := readROM("roms/chip8-roms/programs/IBM Logo.ch8")
	fmt.Println("Loaded ROM file")
	c.loadGame(bytes)

	running := true
	for running {
		c.emulateCycle()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
		millisecondsToSleep := 1000 / CPU_FREQ
		time.Sleep(time.Duration(millisecondsToSleep) * time.Millisecond)
	}

	window.Destroy()
	sdl.Quit()
}

func setupSDL() *sdl.Window {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(WINDOW_WIDTH*WINDOW_SCALE), int32(WINDOW_HEIGHT*WINDOW_SCALE), sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	return window
}
