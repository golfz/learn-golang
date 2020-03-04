package main

import (
	"log"
	"os"
)

func main() {
	file, err := os.Open("sloth.txt")
	if err != nil {
		log.Fatalf("error open sloth.txt: %v", err)
	}

	// close file before read file
	file.Close()

	bytesRead := make([]byte, 33)

	_, err = file.Read(bytesRead)
	if err != nil {
		log.Fatalf("error reading from sloth.txt: %#v", err)
	}

}
