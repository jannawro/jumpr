package main

import (
	"log"
	"os"

	"github.com/jannawro/jumpr/internal/cli"
)

func main() {
	switch len(os.Args) {
	case 1:
		cli.TeaPrompt()
	case 2:
		cli.InlineLogin(os.Args[1])
	default:
		log.Fatal("Unsupported number of arguments. Use either 1 or none.")
	}
}
