package main

import (
	"fmt"
	"log"

	"github.com/alecthomas/kong"
	"github.com/chzyer/readline"
	"github.com/tarm/serial"
)

const (
	COLOR_RED   = "\033[0;31m"
	COLOR_GREEN = "\033[1;32m"
	COLOR_BLUE  = "\033[1;34m"

	COLOR_OFF = "\033[m"
)

var CLI struct {
	BaudRate int    `required help:"Baud rate" short:"b"`
	Path     string `arg name:"path" help: "USB device path." type:"path"`
}

func main() {

	kong.Parse(&CLI)
	fmt.Println(CLI.BaudRate, CLI.Path)

	c := &serial.Config{
		Name: CLI.Path,
		Baud: CLI.BaudRate,
	}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          fmt.Sprintf("%s>%s ", COLOR_GREEN, COLOR_OFF),
		InterruptPrompt: fmt.Sprintf("\n%sinterrupt%s", COLOR_RED, COLOR_OFF),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()

	buf := make([]byte, 128)

	for {
		line, err := rl.Readline()
		if err != nil {
			log.Fatal(err)
		}

		if len(line) == 0 {
			continue
		}

		s.Write([]byte(line)[:1])

		n, err := s.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%q\n", buf[:n])
	}
}
