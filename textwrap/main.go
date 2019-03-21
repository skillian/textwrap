package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/skillian/textwrap"
)

func main() {
	var columns int
	flag.IntVar(&columns, "c", 80, "Columns to wrap output")

	defaultUsage := flag.Usage

	flag.Usage = func() {
		defaultUsage()
		fmt.Println(`
positional arguments:
  source
                The source file to read from (default: "-" for stdin)

  target
                The target file to write to (default: "-" for stdout)`)
	}

	flag.Parse()

	args := flag.Args()

	sourceName, targetName := "", ""

	switch len(args) {
	case 2:
		targetName = args[1]
		fallthrough
	case 1:
		sourceName = args[0]
	default:
		log.Fatal("expected source and target filenames")
	}

	var err error

	var source io.ReadCloser
	if sourceName == "" || sourceName == "-" {
		source = ioutil.NopCloser(os.Stdin)
	} else {
		source, err = os.Open(sourceName)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer logClose(source)

	var target io.WriteCloser
	if targetName == "" || targetName == "-" {
		target = &nopWriteCloser{os.Stdout}
	} else {
		target, err = os.Create(sourceName)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer logClose(target)

	source, err = textwrap.NewReader(source, columns)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = io.Copy(target, source); err != nil {
		log.Fatal(err)
	}
}

func logClose(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Fatal(err)
	}
}

type nopWriteCloser struct{ io.Writer }

func (wc *nopWriteCloser) Close() error { return nil }
