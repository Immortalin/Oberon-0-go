/*
OSS.go: The Oberon-0 Scanner
24.07.2016 TSS
Texts.Read(R, ch) -> ch, err = r.ReadByte()
*/

package OSS

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	IdLen     = 16
	Null      = 0
	Times     = 1
	Div       = 3
	Mod       = 4
	And       = 5
	Plus      = 6
	Minus     = 7
	Or        = 8
	Eql       = 9
	Neq       = 10
	Lss       = 11
	Geq       = 12
	Leq       = 13
	Gtr       = 14
	Period    = 18
	Comma     = 19
	Colon     = 20
	Rparen    = 22
	Rbrak     = 23
	Of        = 25
	Then      = 26
	Do        = 27
	Lparen    = 29
	Lbrak     = 30
	Not       = 32
	Becomes   = 33
	Number    = 34
	Ident     = 37
	Semicolon = 38
	End       = 40
	Else      = 41
	Elsif     = 42
	If        = 44
	While     = 46
	Array     = 54
	Record    = 55
	Const     = 57
	Type      = 58
	Var       = 59
	Procedure = 60
	Begin     = 61
	Module    = 63
	Eof       = 64
)

type ident [IdLen]byte

var (
	Val    int
	Id     []byte
	Error  bool
	ch     byte
	err    error
	r      *bufio.Reader
	keyTab = []struct {
		sym int
		id  string
	}{{Null, "BY"},
		{Do, "DO"},
		{If, "IF"},
		{Null, "IN"},
		{Null, "IS"},
		{Of, "OF"},
		{Or, "OR"},
		{Null, "TO"},
		{End, "END"},
		{Null, "FOR"},
		{Mod, "MOD"},
		{Null, "NIL"},
		{Var, "VAR"},
		{Null, "CASE"},
		{Else, "ELSE"},
		{Null, "EXIT"},
		{Then, "THEN"},
		{Type, "TYPE"},
		{Null, "WITH"},
		{Array, "ARRAY"},
		{Begin, "BEGIN"},
		{Const, "CONST"},
		{Elsif, "ELSIF"},
		{Null, "IMPORT"},
		{Null, "UNTIL"},
		{While, "WHILE"},
		{Record, "RECORD"},
		{Null, "REPEAT"},
		{Null, "RETURN"},
		{Null, "POINTER"},
		{Procedure, "PROCEDURE"},
		{Div, "DIV"},
		{Null, "LOOP"},
		{Module, "MODULE"}}
)

func Mark(msg string) {
	fmt.Println(msg)
	Error = true
}

func identifier() int {
	Id = Id[:0]
	i := 0
	for {
		if i < 16 {
			Id = append(Id, ch)
			i += 1
		}
		ch, _ = r.ReadByte()
		if (ch < '0') ||
			(ch > '9' && ch < 'A') ||
			(ch > 'Z' && ch < 'a') ||
			(ch > 'z') {
			break
		}
	}
	k := 0
	for k < len(keyTab) {
		if keyTab[k].id == string(Id[:]) {
			break
		}
		k += 1
	}
	if k < len(keyTab) {
		return keyTab[k].sym
	}
	return Ident
}

func number() (n int) {
	n = 0
	for {
		if n <= 1048576 {
			n = 10*n + int(ch-'0')
		} else {
			Mark("Number too large")
			n = 0
		}
		ch, err = r.ReadByte()
		if (ch < '0') || (ch > '9') {
			break
		}
	}
	return
}

func comment() {
	ch, err = r.ReadByte()
	for {
		for {
			for ch == '(' {
				ch, err = r.ReadByte()
				if ch == '*' {
					comment()
				}
			}
			if ch == '*' {
				ch, err = r.ReadByte()
				break
			}
			if err == io.EOF {
				break
			}
			ch, err = r.ReadByte()
		}
		if ch == ')' {
			ch, err = r.ReadByte()
			break
		}
		if err == io.EOF {
			Mark("Comment not terminated")
			break
		}
	}
}

func Get(sym *int) {
	for err != io.EOF && (ch <= ' ') {
		ch, err = r.ReadByte()
	}
	if err == io.EOF {
		*sym = Eof
	} else {
		switch {
		case ch == '&':
			ch, err = r.ReadByte()
			*sym = And
		case ch == '*':
			ch, err = r.ReadByte()
			*sym = Times
		case ch == '+':
			ch, err = r.ReadByte()
			*sym = Plus
		case ch == '-':
			ch, err = r.ReadByte()
			*sym = Minus
		case ch == '=':
			ch, err = r.ReadByte()
			*sym = Eql
		case ch == '#':
			ch, err = r.ReadByte()
			*sym = Neq
		case ch == '<':
			ch, err = r.ReadByte()
			if ch == '=' {
				ch, err = r.ReadByte()
				*sym = Leq
			} else {
				*sym = Lss
			}
		case ch == '>':
			ch, err = r.ReadByte()
			if ch == '=' {
				ch, err = r.ReadByte()
				*sym = Geq
			} else {
				*sym = Gtr
			}
		case ch == ';':
			ch, err = r.ReadByte()
			*sym = Semicolon
		case ch == ',':
			ch, err = r.ReadByte()
			*sym = Comma
		case ch == ':':
			ch, err = r.ReadByte()
			if ch == '=' {
				ch, err = r.ReadByte()
				*sym = Becomes
			} else {
				*sym = Colon
			}
		case ch == '.':
			ch, err = r.ReadByte()
			*sym = Period
		case ch == '(':
			ch, err = r.ReadByte()
			if ch == '*' {
				comment()
				ch, err = r.ReadByte()
			} else {
				*sym = Lparen
			}
		case ch == ')':
			ch, err = r.ReadByte()
			*sym = Rparen
		case ch == '[':
			ch, err = r.ReadByte()
			*sym = Lbrak
		case ch == ']':
			ch, err = r.ReadByte()
			*sym = Rbrak
		case ch >= '0' && ch <= '9':
			Val = number()
			*sym = Number
		case (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z'):
			*sym = identifier()
		case ch == '~':
			ch, err = r.ReadByte()
			*sym = Not
		default:
			ch, err = r.ReadByte()
			*sym = Null
		}
	}
}

/* Scanner init */
func Init(fname string) {
	file, err := os.Open(fname)
	if err != nil {
		fmt.Println(err)
	} else {
		r = bufio.NewReader(file)
		ch, _ = r.ReadByte()
		Error = false
	}
}

/* Go version of Module "body" */
func init() {
	Error = true
	Id = make([]byte, 16)
}
