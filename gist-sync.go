//usr/bin/env go run $0 -- "$@"; exit $?

//
// Use this script to maintain a set of github gists as files.
//
// Usage:
//  - install 'gist' program (https://github.com/defunkt/gist)
//  - make sure you have Go (https://golang.org) installed
//  - make a folder somewhere where you want to store your gists.
//  - place gist-update.go (https://gist.github.com/ivanzoid/611177bbd3f5cb0604810f07080757b3#file-gist-update-go) to this folder
//  - place your gists as files to this folder
//  - when you create a new or update existing gist file, just run: './gist-update.go gist.foo'
//  - to sync all your gists, you may run './gist-update.go *.*'
//

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	GISTIDS_DIR = ".gistids"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %v file1 [file2 [...]]]\n", filepath.Base(os.Args[0]))
	os.Exit(2)
}

func dlog(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintf(os.Stderr, "\n")
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		usage()
	}

	for i := 0; i < flag.NArg(); i++ {
		filename := flag.Arg(i)
		processFile(filename)
	}
}

func processFile(filename string) {
	gistIdFilename := fmt.Sprintf("%v/%v.id", GISTIDS_DIR, filename)

	if fileExists(gistIdFilename) {
		gistId, err := readFirstLineFromFile(gistIdFilename)
		if err != nil {
			dlog("Can't read %v: %v", gistIdFilename, err)
			return
		}

		out, err := runProgram1("gist", "-u", gistId, filename)
		if err != nil {
			dlog("Error: %v", err)
			dlog("Out: %v", out)
		}
	} else {
		out, err := runProgram1("gist", filename)
		if err != nil {
			dlog("Error: %v", err)
			dlog("Out: %v", out)
			return
		}

		gistUrl := out

		gistId, err := gistIdFromUrl(gistUrl)
		if err != nil {
			dlog("%v", err)
			return
		}

		makeDirectoryIfNotExists(GISTIDS_DIR)

		writeStringToFile(gistIdFilename, gistId)
	}
}

func gistIdFromUrl(gistUrl string) (string, error) {
	comps := strings.Split(gistUrl, "/")
	if len(comps) == 0 {
		return "", errors.New("Can't parse url")
	}
	return comps[len(comps)-1], nil
}

func runProgram1(program string, args ...string) (string, error) {
	outStrings, err := runProgram(program, args...)
	if err != nil {
		return "", err
	}
	if len(outStrings) == 0 {
		return "", nil
	}

	return outStrings[0], nil
}

func runProgram(program string, args ...string) ([]string, error) {

	dlog("Running %v", cmdString(program, args))

	cmd := exec.Command(program, args...)

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	outString := string(out)
	if len(outString) == 0 {
		return nil, nil
	}

	outStrings := strings.Split(outString, "\n")
	return outStrings, nil
}

func cmdString(program string, args []string) string {
	comps := make([]string, 0)

	comps = append(comps, program)

	for _, arg := range args {
		if strings.Contains(arg, " ") {
			comps = append(comps, fmt.Sprintf("'%v'", arg))
		} else {
			comps = append(comps, arg)
		}
	}

	result := strings.Join(comps, " ")
	return result
}

func fileExists(reqFilePath string) bool {
	_, err := os.Stat(reqFilePath)
	exists := false
	if err == nil {
		exists = true
	} else {
		exists = !os.IsNotExist(err)
	}
	return exists
}

func readFirstLineFromFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	scanner.Scan()
	return scanner.Text(), nil
}

func writeStringToFile(path string, s string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	_, err = file.WriteString(s)
	file.Sync()
	file.Close()

	return err
}

func makeDirectoryIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.Mkdir(path, os.ModeDir|0755)
	}
	return nil
}
