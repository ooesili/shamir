package main

import (
	"bufio"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/hashicorp/vault/shamir"
)

var config struct {
	in        string
	out       string
	threshold int
	parts     int
}

func init() {
	flag.StringVar(&config.in, "in", "-", "input file")
	flag.StringVar(&config.out, "out", "-", "output file")
	flag.IntVar(&config.threshold, "threshold", 2, "threshold for split operation")
	flag.IntVar(&config.parts, "parts", 3, "number of shares for split operation")

	flag.Usage = func() {
		fmt.Printf("Usage: %s [<options>] (split|combine)\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if err := mainErr(flag.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func mainErr(args []string) error {
	if len(args) == 0 {
		return errors.New("no operation given")
	}
	if len(args) > 1 {
		return errors.New("too many arguments")
	}

	var operation operation
	switch args[0] {
	case "combine":
		operation = combine
	case "split":
		operation = split
	default:
		return fmt.Errorf("invalid operation: %s", args[0])
	}

	var in io.Reader
	if config.in == "-" {
		in = os.Stdin
	} else {
		file, err := os.Open(config.in)
		if err != nil {
			return fmt.Errorf("opening input file: %s", err)
		}
		defer func() { _ = file.Close() }()
		in = file
	}

	var out io.Writer
	if config.out == "-" {
		out = os.Stdout
	} else {
		file, err := os.Create(config.out)
		if err != nil {
			return fmt.Errorf("opening out file: %s", err)
		}
		defer func() { _ = file.Close() }()
		out = file
	}

	return operation(in, out)
}

type operation func(in io.Reader, out io.Writer) error

func combine(in io.Reader, out io.Writer) error {
	var shares [][]byte
	scanner := bufio.NewScanner(in)
	encoding := base64.StdEncoding

	for lineno := 1; scanner.Scan(); lineno++ {
		line := scanner.Bytes()
		share := make([]byte, encoding.DecodedLen(len(line)))

		n, err := encoding.Decode(share, line)
		if err != nil {
			return fmt.Errorf("base64 decoding input: line %d: %s", lineno, err)
		}

		shares = append(shares, share[:n])
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("reading key shares from input: %s", err)
	}

	result, err := shamir.Combine(shares)
	if err != nil {
		return fmt.Errorf("combining key shares: %s", err)
	}

	if _, err := out.Write(result); err != nil {
		return fmt.Errorf("writing secret to output: %s", err)
	}

	return nil
}

func split(in io.Reader, out io.Writer) error {
	secret, err := ioutil.ReadAll(in)
	if err != nil {
		return fmt.Errorf("reading secret from input: %s", err)
	}

	shares, err := shamir.Split(secret, config.parts, config.threshold)
	if err != nil {
		return fmt.Errorf("splitting secret: %s", err)
	}

	encoding := base64.StdEncoding

	for _, share := range shares {
		lineLen := encoding.EncodedLen(len(share)) + 1
		line := make([]byte, lineLen)

		encoding.Encode(line, share)
		line[len(line)-1] = '\n'

		if _, err := out.Write(line); err != nil {
			return fmt.Errorf("writing shares to output: %s", err)
		}
	}

	return nil
}
