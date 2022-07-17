// main entry point for program
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

func readROM(path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Error while opening file", err)
	}
	//defer file.Close()

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
	verbosePtr := flag.Bool("verbose", false, "a bool")

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

	/*surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0)

	rect := sdl.Rect{0, 0, 200, 200}
	surface.FillRect(&rect, 0xffff0000)
	window.UpdateSurface()
	*/
	running := true
	for running {
		//c.emulateCycle()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}
}

func setupSDL() *sdl.Window {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	//defer sdl.Quit()

	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 64, 32, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	//defer window.Destroy()

	return window
}
