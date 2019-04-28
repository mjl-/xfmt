package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"strings"
)

var (
	width = flag.Int("width", 80, "max width of a line, not including non-text prefix")
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

	out := bufio.NewWriter(os.Stdout)

	write := func(s string) {
		_, err := out.Write([]byte(s))
		if err != nil {
			log.Fatalf("write: %s", err)
		}
	}

	p := &parser{in: bufio.NewReader(os.Stdin)}

	for {
		line, end := p.gatherLine()
		if line == "" && end == "" {
			break
		}
		prefix, rem := parseLine(line)
		for _, s := range splitLine(rem) {
			write(prefix)
			write(s)
			write(end)
		}
	}
	err := out.Flush()
	if err != nil {
		log.Fatalf("write: %s\n", err)
	}
}

type parser struct {
	in         *bufio.Reader
	curLine    string
	curLineend string
}

func (p *parser) peekLine() (string, string) {
	if p.curLine != "" || p.curLineend != "" {
		return p.curLine, p.curLineend
	}

	line, err := p.in.ReadString('\n')
	if err != nil && err != io.EOF {
		log.Fatalf("read: %s\n", err)
	}
	if line == "" {
		return "", ""
	}
	if strings.HasSuffix(line, "\r\n") {
		p.curLine, p.curLineend = line[:len(line)-2], "\r\n"
	} else if strings.HasSuffix(line, "\n") {
		p.curLine, p.curLineend = line[:len(line)-1], "\n"
	} else {
		p.curLine, p.curLineend = line, ""
	}
	return p.curLine, p.curLineend
}

func (p *parser) consumeLine() {
	if p.curLine == "" && p.curLineend == "" {
		panic("bad")
	}
	p.curLine = ""
	p.curLineend = ""
}

func (p *parser) gatherLine() (string, string) {
	var curLine, curLineend string
	var curPrefix string

	for {
		line, end := p.peekLine()
		if line == "" && end == "" {
			break
		}
		if curLine == "" {
			curLineend = end
		}
		prefix, rem := parseLine(line)
		if curLine != "" && (curPrefix != prefix || rem == "" || startsWithNonText(rem)) {
			break
		}
		curPrefix = prefix
		if curLine != "" {
			curLine += " "
		}
		curLine += rem
		p.consumeLine()
	}

	return curPrefix + curLine, curLineend
}

func startsWithNonText(s string) bool {
	c := s[0]
	return c < 0x80 && !(c >= 'A' && c < 'Z') && !(c >= 'a' && c <= 'z')
}

func parseLine(s string) (string, string) {
	orig := s
	s = strings.TrimLeft(orig, " \t")
	prefix := orig[:len(orig)-len(s)]
	if strings.HasPrefix(s, "//") {
		prefix += "//"
		s = s[2:]
	} else if strings.HasPrefix(s, "#") {
		prefix += "#"
		s = s[1:]
	}
	ns := strings.TrimLeft(s, " \t")
	prefix += s[:len(s)-len(ns)]
	s = ns
	return prefix, s
}

func splitLine(s string) []string {
	if len(s) <= *width {
		return []string{s}
	}

	line := ""
	r := []string{}
	for _, w := range strings.Split(s, " ") {
		if line != "" && len(line)+1+len(w) > *width {
			r = append(r, line)
			line = w
			continue
		}
		if line != "" {
			line += " "
		}
		line += w
	}
	if line != "" {
		r = append(r, line)
	}
	return r
}
