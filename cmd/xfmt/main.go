// Xfmt reads text from stdin, reformats it and writes it stdout, wrapping lines
// and recognizing code comments.
package main

import (
	"flag"
	"log"
	"os"

	"github.com/mjl-/xfmt"
)

var (
	width = flag.Int("width", 80, "max width of a line, not including non-text prefix")
	dash  = flag.Bool("dash", false, "merge lines that start with a dash (\"-\")")
	star  = flag.Bool("star", false, "merge lines that start with an asterisk (\"*\")")
)

func main() {
	log.SetFlags(0)
	flag.Usage = func() {
		log.Println("usage: xfmt [flags]")
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()
	if len(args) != 0 {
		flag.Usage()
		os.Exit(2)
	}

	config := xfmt.Config{
		MaxWidth:      *width,
		BreakPrefixes: []string{},
	}
	if !*dash {
		config.BreakPrefixes = append(config.BreakPrefixes, "- ")
	}
	if !*star {
		config.BreakPrefixes = append(config.BreakPrefixes, "* ")
	}

	err := xfmt.Format(os.Stdout, os.Stdin, config)
	if err != nil {
		log.Fatalf("format: %s", err)
	}
}
