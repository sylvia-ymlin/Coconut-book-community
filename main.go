package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/sylvia-ymlin/Coconut-book-community/docs"
	"github.com/sylvia-ymlin/Coconut-book-community/initiate"
)

// @title BookCommunity API
// @version 1.0
// @description High-performance book community backend API built with Go
// @description
// @description Tech Stack:
// @description - Backend: Go 1.20 + Gin + GORM
// @description - Database: PostgreSQL 15
// @description - Cache: Redis 7.0 + In-Memory LRU
// @description - Message Queue: RabbitMQ 3.12
//
// @contact.name API Support
// @contact.url https://github.com/sylvia-ymlin/Coconut-book-community
// @contact.email support@bookcommunity.io
//
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
//
// @host localhost:8080
// @BasePath /douyin
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

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
