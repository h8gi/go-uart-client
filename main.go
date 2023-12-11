package main

import (
	"fmt"
	"log"

	"github.com/alecthomas/kong"
	"github.com/chzyer/readline"
	"github.com/tarm/serial"
)

const (
	ColorRed   = "\033[0;31m"
	ColorGreen = "\033[1;32m"
	ColorBlue  = "\033[1;34m"

	ColorOff = "\033[m"
)

var CLI struct {
	Path     string `arg name:"path" help: "USB device path." type:"path"`
	BaudRate int    `required help:"Baud rate" short:"b"`
	Repl     bool   `help:"launch repl mode" default:"false" short:"r"`
}

func repl(s *serial.Port) error {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          fmt.Sprintf("%s>%s ", ColorGreen, ColorOff),
		InterruptPrompt: fmt.Sprintf("\n%sinterrupt%s", ColorRed, ColorOff),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()

	buf := make([]byte, 128)

	for {
		line, err := rl.Readline()
		if err != nil {
			return err
		}

		if len(line) == 0 {
			continue
		}

		n, err := s.Write([]byte(line)[:1])
		if err != nil {
			return err
		}
		fmt.Printf("  write %d byte\n", n)

		n, err = s.Read(buf)
		if err != nil {
			return err
		}
		fmt.Printf("%q\n", buf[:n])
	}
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

	if CLI.Repl {
		err := repl(s)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		buf := make([]byte, 128)

		n, err := s.Write([]byte("Hello, World!"))
		fmt.Printf("(write %d byte)\n", n)
		if err != nil {
			log.Fatal(err)
		}

		n, err = s.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", buf[:n])
	}
}
