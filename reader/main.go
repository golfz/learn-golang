package main

import (
	"log"
	"os"
)

func main() {
	file, err := os.Open("sloth.txt")
	if err != nil {
		log.Fatalf("error opening sloth.txt: %v", err)
	}
	defer file.Close()

	bytesRead := make([]byte, 33)

	n, err := file.Read(bytesRead)
	if err != nil {
		log.Fatalf("error reading from sloth.txt: %v", err)
	}

	log.Printf("We read \"%s\" into bytesRead (%d bytes)", string(bytesRead), n)

	n, err = file.Read(bytesRead)
	if err != nil {
		log.Fatalf("error reading from sloth.txt: %v", err)
	}

	log.Printf("We read \"%s\" into bytesRead (%d bytes)", string(bytesRead), n)
}
