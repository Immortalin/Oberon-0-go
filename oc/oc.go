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
	"fmt"
	"os"
)

func main() {

	if !(len(os.Args) > 1) {
		fmt.Printf("Usage: oc <filename.m>\n")
		return
	}
	// ##TODO: Add option parsing and set this with -d flag
	OSP.Dump = true

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		OSP.Compile(bufio.NewReader(file))
	}
}
