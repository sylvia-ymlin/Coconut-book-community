package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sylvia-ymlin/Coconut-book-community/initiate"
)

var (
	Version   = "dev"
	Commit    = "none"
	BuildTime = "unknown"

	showVersion = flag.Bool("version", false, "Show version information")
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("BookCommunity %s\n", Version)
		fmt.Printf("Commit: %s\n", Commit)
		fmt.Printf("Built: %s\n", BuildTime)
		os.Exit(0)
	}

	initiate.Run()
}
