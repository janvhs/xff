package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
)

type FileBuffer struct {
	Path  string
	Value []byte
	File  *os.File
}

func main() {
	appName := filepath.Base(os.Args[0])

	flag := pflag.NewFlagSet(appName, pflag.ExitOnError)

	// Define flags
	// ---------------------------------------------------

	var modeDecimal bool
	flag.BoolVarP(&modeDecimal, "decimal", "d", false, "Show bytes as decimal")

	// Set flags
	// ---------------------------------------------------

	flag.Parse(os.Args)

	// ---------------------------------------------------

	args := flag.Args()

	if len(args) <= 1 {
		fmt.Printf("usage %s: [options] file\n", appName)
		os.Exit(1)
	}

	absFilePath, err := filepath.Abs(args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	f, err := os.Open(absFilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	defer f.Close()

	fileBytes, err := io.ReadAll(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	currentBuffer := &FileBuffer{
		Path:  absFilePath,
		Value: fileBytes,
		File:  f,
	}

	fmt.Printf("%d\n", currentBuffer.Value)
}
