package main

import (
	"log"
	"os"
	"strings"

	"github.com/mgrim/sakpro/cleaner"
)

func main() {
	arg := os.Args[1]

	file, err := os.Open(arg)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	html, err := cleaner.CleanHTML(file)
	if err != nil {
		log.Fatal(err)
	}

	target, err := os.Create(strings.Replace(arg, ".htm", "_clean.htm", 1))
	if err != nil {
		log.Fatal(err)
	}
	defer target.Close()

	target.WriteString(html)
}
