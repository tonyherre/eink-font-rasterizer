package main

import (
	"bufio"
	"fmt"
	"log"
	"github.com/tarm/serial"
)

func writeToArduino(bytes []byte) {
	fmt.Println("Opening port")
	c := &serial.Config{Name: "COM4", Baud: 115200}
	s, err := serial.OpenPort(c)
	if err != nil {
			log.Fatal(err)
	}

	scanner := bufio.NewScanner(s)
	buf := make([]byte, 1)
	n, err := s.Read(buf)
	if err != nil {
			log.Fatal(err)
	}
	if n != 1 || buf[0] != 'g' {
		log.Fatalf("Unexpected char %s", string(buf))
	}

	for i := 0; i < len(bytes); i += 100 {
		end := i + 100
		if end > len(bytes) {
			end = len(bytes)
		}
		n, err := s.Write(bytes[i:end])
		if err != nil {
				log.Fatal(err)
		}
        buf := make([]byte, 1)
		n, err = s.Read(buf)
		if err != nil {
				log.Fatal(err)
		}
		if n != 1 || buf[0] != 'l' {
			log.Fatalf("Unexpected char %s", string(buf))
		}
	}

	for scanner.Scan() {
		fmt.Println("Serial: " + scanner.Text())
	}
}