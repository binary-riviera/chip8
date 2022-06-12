// main entry point for program
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
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

	fmt.Println("Starting chip-8 emulator")
	c := chip8{}
	c.initialise(*verbosePtr)
	fmt.Println("Initialised emulator...")
	bytes := readROM("roms/chip8-roms/programs/IBM Logo.ch8")
	fmt.Println("Loaded ROM file")
	c.loadGame(bytes)

	// emulation loop
	for {
		// emulate one cycle
		c.emulateCycle()
	}
}
