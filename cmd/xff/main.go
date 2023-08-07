// TODO: Setup golangci to lint
// TODO: Wrap errors
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
)

// TODO: Add custom exit codes depending on error
var (
	ErrHelpRequested error = errors.New("xff: help requested")
	ErrUsage               = errors.New("xff: wrong usage")
)

type FileBuffer struct {
	Path  string
	Value []byte
	File  *os.File
}

// Only pass exit codes between 0 and 125
func printUsage(w io.Writer, exitCode uint, appName string, flag *pflag.FlagSet) {
	options := flag.FlagUsages()
	fmt.Fprintf(
		w,
		`usage:  %s [options] file

options:
%s`,
		appName,
		options,
	)
	os.Exit(int(exitCode))
}

// Do not defer functions in this proc, because they will not be executed, due to os.Exit()
func main() {
	
	appName := filepath.Base(os.Args[0])

	flag := pflag.NewFlagSet(appName, pflag.ContinueOnError)

	// Make Usage a no-op, because I handle this on my own
	flag.Usage = func() {}

	if err := mainE(appName, flag); err != nil {
		if errors.Is(err, ErrHelpRequested) {
			printUsage(os.Stdout, 0, appName, flag)
		} else if errors.Is(err, ErrUsage) {
			printUsage(os.Stderr, 1, appName, flag)
		} else {
			fmt.Fprintln(os.Stderr, err)
		}
	} else {
		os.Exit(0)
	}
}

type flagOptions struct {
	modeDecimal bool
}

func parseOptions(flag *pflag.FlagSet) (flagOptions, error) {

	fOpt := flagOptions{}

	// Define flags
	// ---------------------------------------------------

	flag.BoolVarP(&fOpt.modeDecimal, "decimal", "d", false, "Show bytes as decimal")

	// Set flags
	// ---------------------------------------------------

	if err := flag.Parse(os.Args); err != nil {
		if errors.Is(err, pflag.ErrHelp) {
			return fOpt, ErrHelpRequested
		} else {
			return fOpt, err
		}
	}

	return fOpt, nil
}

// Do not call os.Exit in this function, otherwise deferred functions will not execute!
func mainE(appName string, flag *pflag.FlagSet) error {

	options, err := parseOptions(flag)
	if err != nil {
		return err
	}

	args := flag.Args()

	if len(args) <= 1 {
		return ErrUsage
	}

	// ---------------------------------------------------

	if options.modeDecimal {
		fmt.Println("Decimal mode requested")
	}

	absFilePath, err := filepath.Abs(args[1])
	if err != nil {
		return err
	}

	// FIXME: Maybe create it, when it does not exist
	f, err := os.OpenFile(absFilePath, os.O_RDWR, 0)
	if err != nil {
		return err
	}

	fileBytes, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	currentBuffer := &FileBuffer{
		Path:  absFilePath,
		Value: fileBytes,
		File:  f,
	}

	defer currentBuffer.File.Close()

	fmt.Printf("%d\n", currentBuffer.Value)

	return nil
}
