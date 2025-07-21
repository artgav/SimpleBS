package main

import (
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		os.Args = append([]string{os.Args[0]}, os.Args[2:]...)
		runGRPCServer()
	} else {
		runCLI()
	}
}
