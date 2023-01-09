package main

import (
	"log"
	"os"

	"github.com/jannawro/jumpr/internal"
)

func main() {
	switch len(os.Args) {
	case 1:
		internal.TeaPrompt()
	case 2:
		internal.InlineLogin(os.Args[1])
	default:
		log.Fatal("Unsupported number of arguments. Use either 1 or none.")
	}
}
