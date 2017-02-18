/*
Build with:
C:\> go install Oberon-0-go/oc
Run:
C:\Data\Personal\go\bin> oc Test.m
*/

package main

import (
	"Oberon-0-go/OSP"
	"bufio"
	"flag"
	"fmt"
	"os"
)

const Ver = "Oberon-0 compiler v0.5-alpha"

func init() {
	flag.BoolVar(&OSP.Dump, "d", false, "Dump listing output to terminal")
}

func main() {

	fmt.Printf("%s\n\n", Ver)

	// Handle args
	flag.Parse()
	// Exit on error
	if !(len(flag.Args()) > 0) {
		fmt.Printf("Usage: oc <flags> sourcefile.m\n")
		flag.PrintDefaults()
		return
	}

	// Compile source file
	filename := flag.Arg(0)
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		OSP.Compile(bufio.NewReader(file))
	}
}
