// main entry point for program
package main

import "fmt"

func main() {
	fmt.Println("Starting chip-8 emulator")
	c := chip8{}
	c.initialise()
	fmt.Println("Initialised emulator...")

}
