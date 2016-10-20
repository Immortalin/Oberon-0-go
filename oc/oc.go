/*
Build with:
C:\> go install oc
Run:
C:\Data\Personal\go\bin> oc Test.m
*/

package main

import (
	"OSP"
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
	OSP.Compile(os.Args[1])
}
