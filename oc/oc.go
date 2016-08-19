/*
Build with:
C:\> go install oc
Run:
C:\Data\Personal\go\bin> oc Test.m
*/

package main

import (
    "os"
    "fmt"
	"OSP"
)

func main() {
    
    if !(len(os.Args) > 1) {
        fmt.Printf("Usage: oc <filename.m>\n")
        return
    }
    OSP.Compile(os.Args[1])
	
}
