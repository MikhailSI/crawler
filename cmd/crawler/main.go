package main

import (
	"log"
	"os"
	"semtest/pkg/parser"
	"strconv"
)

const (
	RPSDefault = 10
)

func main() {
	url := os.Getenv("URL")
	RPS, _ := strconv.Atoi(os.Getenv("RPS"))

	if url == "" {
		log.Fatalln("env URL not specified")
	}

	if RPS == 0 {
		RPS = RPSDefault
	}

	parser.Run(url, RPS)
}
