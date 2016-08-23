/*
OSG.go: The Oberon-0 Code Generator
25.07.2016 TSS
*/

package OSG

import (
	"OSS"
	"RISC"
	"fmt"
)

const (
	maxCode = 128
	maxRel  = 200
	nofCom  = 16
)

// Class/mode
const (
	Head  = 0
	Var   = 1
	Par   = 2
	Const = 3
	Fld   = 4
	Typ   = 5
	Proc  = 6
	SProc = 7
	Reg   = 10
	Cond  = 11
)

// Form
const (
	Boolean = 0
	Integer = 1
	Array   = 2
	Record  = 3
)

// Mnemonics
// Not sure why these are not imported from RISC
const (
	MOV  = 0
	MVN  = 1
	ADD  = 2
	SUB  = 3
	MUL  = 4
	Div  = 5
	Mod  = 6
	CMP  = 7
	MOVI = 16
	MVNI = 17
	ADDI = 18
	SUBI = 19
	MULI = 20
	DIVI = 21
	MODI = 22
	CMPI = 23
	CHKI = 24
	LDW  = 32
	LDB  = 33
	POP  = 34
	STW  = 36
	STB  = 37
	PSH  = 38
	RD   = 40
	WRD  = 41
	WRH  = 42
	WRL  = 43
	BEQ  = 48
	BNE  = 49
	BLT  = 50
	BGE  = 51
	BLE  = 52
	BGT  = 53
	BR   = 56
	BSR  = 57
	RET  = 58
)

// Reserved registers
const (
	FP  = 12
	SP  = 13
	LNK = 14
	PC  = 15
)

type Item struct {
	Mode, Lev  int
	Tp         Type
	A, b, c, r int
}

type Object *ObjDesc
type ObjDesc struct {
	Class, Lev int
	Next, Dsc  Object
	Tp         Type
	Name       []byte
	Val        int
}

type Type *TypeDesc
type TypeDesc struct {
	Form      int
	Fields    Object
	Base      Type
	Size, Len int
}

var (
	BoolType, IntType Type
	Curlev            int
	Pc, entry         int
	regs              [32]bool
	code              [maxCode]int
)

var mnemo = []string{
	"MOV", "MVN", "ADD", "SUB", "MUL", "DIV", "MOD", "CMP", "OP8", "OP9", "OP10", "OP11", "OP12", "OP13", "OP14", "OP15",
	"MOVI", "MVNI", "ADDI", "SUBI", "MULI", "DIVI", "MODI", "CMPI", "CHKI", "OP25", "OP26", "OP27", "OP28", "OP29", "OP30", "OP31",
	"LDW", "LDB", "POP", "OP35",
	"STW", "STB", "PSH", "OP39",
	"RD", "WRD", "WRH", "WRL", "OP44", "OP45", "OP46", "OP47",
	"BEQ", "BNE", "BLT", "BGE", "BLE", "BGT", "OP54", "OP55", "BR", "BSR", "RET"}

func getReg(r *int) {

	*r = 0
	for (*r < FP) && regs[*r] {
		*r += 1
	}
	regs[*r] = true
}

func put(op, a, b, c int) {
	if op >= 32 {
		op -= 64
	}
	//code[Pc] = ASH(ASH(ASH(op,4)+a,4)+b, 18) + (c MOD 40000H)
	code[Pc] = (((op<<4|a)<<4 | b) << 18) | (c & 0x3FFFF)
	Pc += 1
}

func putBR(op, disp int) {
	//code[Pc] = ASH(op-40H, 26) + (disp MOD 4000000H)
	code[Pc] = (op-0x40)<<26 | (disp & 0x3FFFFFF)
	Pc += 1
}

func testRange(x int) {
	if (x >= 0x20000) || (x < -0x20000) {
		OSS.Mark("TestRange: value too large")
	}
}

func load(x *Item) {
	var r int

	if x.Mode == Var {
		if x.Lev == 0 {
			x.A -= Pc * 4
		}
		getReg(&r)
		put(LDW, r, x.r, x.A)
		regs[x.r] = false
		x.r = r
	} else if x.Mode == Const {
		testRange(x.A)
		getReg(&x.r)
		put(MOVI, x.r, 0, x.A)
	}
	x.Mode = Reg
}

func loadBool(x *Item) {
	if x.Tp.Form != Boolean {
		OSS.Mark("loadBool: boolean expected")
	}
	load(x)
	x.Mode = Cond
	x.A = 0
	x.b = 0
	x.c = 1
}

func putOp(cd int, x *Item, y *Item) {
	if x.Mode != Reg {
		load(x)
	}
	if y.Mode == Const {
		testRange(y.A)
		put(cd+16, x.r, x.r, y.A)
	} else {
		if y.Mode != Reg {
			load(y)
		}
		put(cd, x.r, x.r, y.r)
		regs[y.r] = false
	}
}

func negated(cond int) int {
	if (cond % 2) != 0 {
		return cond - 1
	}
	return cond + 1
}

func merged(L0 int, L1 int) int {
	var L2, L3 int

	if L0 != 0 {
		L2 = L0
		for {
			//L3 = code[L2] % 0x40000
			L3 = code[L2] & 0x3FFFF
			if L3 == 0 {
				break
			}
			L2 = L3
		}
		code[L2] = code[L2] - L3 + L1
		return L0
	}
	return L1
}

func fix(at int, with int) {
	//code[at] = code[at]%0x400000*0x400000 + (with % 0x400000)
	code[at] = (code[at] & 0xFFC00000) | (with & 0x3FFFFF)
}

func FixLink(L int) {
	var L1 int
	for L != 0 {
		//L1 = code[L] % 0x40000
		L1 = code[L] & 0x3FFFF
		fix(L, Pc-L)
		L = L1
	}
}

func IncLevel(n int) {
	Curlev += n
}

func MakeConstItem(x *Item, tp Type, val int) {
	x.Mode = Const
	x.Tp = tp
	x.A = val
}

func MakeItem(x *Item, y Object) {
	var r int

	x.Mode = y.Class
	x.Tp = y.Tp
	x.Lev = y.Lev
	x.A = y.Val
	x.b = 0

	if y.Lev == 0 {
		x.r = PC
	} else if y.Lev == Curlev {
		x.r = FP
	} else {
		x.r = 0
	}

	if y.Class == Par {
		getReg(&r)
		put(LDW, r, x.r, x.A)
		x.Mode = Var
		x.r = r
		x.A = 0
	}
}

// x := x.y
func Field(x *Item, y Object) {
	x.A += y.Val
	x.Tp = y.Tp
}

// x := x[y]
func Index(x *Item, y *Item) {
	if y.Tp != IntType {
		OSS.Mark("Index: array index must be integer")
	}
	if y.Mode == Const {
		if (y.A < 0) || (y.A >= x.Tp.Len) {
			OSS.Mark("Index: array index out of range")
		}
		x.A += y.A * x.Tp.Base.Size
	} else {
		if y.Mode != Reg {
			load(y)
		}
		put(CHKI, y.r, 0, x.Tp.Len)
		put(MULI, y.r, y.r, x.Tp.Base.Size)
		put(ADD, y.r, x.r, y.r)
		regs[x.r] = false
		x.r = y.r
	}
	x.Tp = x.Tp.Base
}

// x := op x
func Op1(op int, x *Item) {
	var t int

	if op == OSS.Minus {
		if x.Tp.Form != Integer {
			OSS.Mark("Op1: integer type expected")
		} else if x.Mode == Const {
			x.A = -x.A
		} else {
			if x.Mode == Var {
				load(x)
			}
			put(MVN, x.r, 0, x.r)
		}
	} else if op == OSS.Not {
		if x.Mode != Cond {
			loadBool(x)
		}
		x.c = negated(x.c)
		t = x.A
		x.A = x.b
		x.b = t
	} else if op == OSS.And {
		if x.Mode != Cond {
			loadBool(x)
		}
		putBR(BEQ+negated(x.c), x.A)
		regs[x.r] = false
		x.A = Pc - 1
		FixLink(x.b)
		x.b = 0
	} else if op == OSS.Or {
		if x.Mode != Cond {
			loadBool(x)
		}
		putBR(BEQ+x.c, x.b)
		regs[x.r] = false
		x.b = Pc - 1
		FixLink(x.A)
		x.A = 0
	}
}

func Op2(op int, x *Item, y *Item) {
	if (x.Tp.Form == Integer) && (y.Tp.Form == Integer) {
		if (x.Mode == Const) && (y.Mode == Const) {
			if op == OSS.Plus {
				x.A += y.A
			} else if op == OSS.Minus {
				x.A -= y.A
			} else if op == OSS.Times {
				x.A *= y.A
			} else if op == OSS.Div {
				x.A /= y.A
			} else if op == OSS.Mod {
				x.A %= y.A
			} else {
				OSS.Mark("Op2: bad operator (constant)")
			}
		} else {
			if op == OSS.Plus {
				putOp(ADD, x, y)
			} else if op == OSS.Minus {
				putOp(SUB, x, y)
			} else if op == OSS.Times {
				putOp(MUL, x, y)
			} else if op == OSS.Div {
				putOp(Div, x, y)
			} else if op == OSS.Mod {
				putOp(Mod, x, y)
			} else {
				OSS.Mark("Op2: bad operator (variable)")
			}
		}
	} else if (x.Tp.Form == Boolean) && (y.Tp.Form == Boolean) {
		if y.Mode != Cond {
			loadBool(y)
		}
		if op == OSS.Or {
			x.A = y.A
			x.b = merged(y.b, x.b)
			x.c = y.c
		} else if op == OSS.And {
			x.A = merged(y.A, x.A)
			x.b = y.b
			x.c = y.c
		}
	} else {
		OSS.Mark("Op2: Unrecognized type in dyadic expression")
	}
}

func Relation(op int, x *Item, y *Item) {
	if (x.Tp.Form != Integer) || (y.Tp.Form != Integer) {
		OSS.Mark("Relation: both arguments must be integers")
	} else {
		putOp(CMP, x, y)
		x.c = op - OSS.Eql
		regs[y.r] = false
	}
	x.Mode = Cond
	x.Tp = BoolType
	x.A = 0
	x.b = 0
}

func Store(x *Item, y *Item) {

	if (x.Tp.Form == Boolean || x.Tp.Form == Integer) && (x.Tp.Form == y.Tp.Form) {
		if y.Mode == Cond {
			put(BEQ+negated(y.c), y.r, 0, y.A)
			regs[y.r] = false
			y.A = Pc - 1
			FixLink(y.b)
			getReg(&y.r)
			put(MOVI, y.r, 0, 1)
			putBR(BR, 2)
			FixLink(y.A)
			put(MOVI, y.r, 0, 0)
		} else if y.Mode != Reg {
			load(y)
		}
		if x.Mode == Var {
			if x.Lev == 0 {
				x.A = x.A - Pc*4
			}
			put(STW, y.r, x.r, x.A)
		} else {
			OSS.Mark("Store: illegal assignment")
		}
		regs[x.r] = false
		regs[y.r] = false
	} else {
		OSS.Mark("Store: incompatible assignment")
	}
}

func Parameter(x *Item, ftyp Type, class int) {
	var r int

	if x.Tp == ftyp {
		if class == Par {
			if x.Mode == Var {
				if x.A != 0 {
					getReg(&r)
                    // Addition to original source: handle VAR parameters at module scope correctly
                    if x.Lev == 0 {
                        x.A -= Pc * 4
                    }
					put(ADDI, r, x.r, x.A)
				} else {
					r = x.r
				}
			} else {
				OSS.Mark("Parameter: VAR param expected")
			}
			put(PSH, r, SP, 4)
			regs[r] = false
		} else {
			if x.Mode != Reg {
				load(x)
			}
			put(PSH, x.r, SP, 4)
			regs[x.r] = false
		}
	} else {
		OSS.Mark("Parameter: unknown parameter type")
	}
}

func CJump(x *Item) {
	if x.Tp.Form == Boolean {
		if x.Mode != Cond {
			loadBool(x)
		}
		putBR(BEQ+negated(x.c), x.A)
		regs[x.r] = false
		FixLink(x.b)
		x.A = Pc - 1
	} else {
		OSS.Mark("CJump: Boolean expected")
		x.A = Pc
	}
}

func BJump(L int) {
	putBR(BR, L-Pc)
}

func FJump(L *int) {
	putBR(BR, *L)
	*L = Pc - 1
}

func Call(x *Item) {
	putBR(BSR, x.A-Pc)
}

func IOCall(x *Item, y *Item) {
	var z Item

	if x.A < 4 {
		if y.Tp.Form != Integer {
			OSS.Mark("IOCall: integer expected")
		}
	}
	if x.A == 1 {
		getReg(&z.r)
		z.Mode = Reg
		z.Tp = IntType
		put(RD, z.r, 0, 0)
		Store(y, &z)
	} else if x.A == 2 {
		load(y)
		put(WRD, 0, 0, y.r)
		regs[y.r] = false
	} else if x.A == 3 {
		load(y)
		put(WRH, 0, 0, y.r)
		regs[y.r] = false
	} else {
		put(WRL, 0, 0, 0)
	}
}

func Header(size int) {
	entry = Pc
	put(MOVI, SP, 0, RISC.MemSize-size)
	put(PSH, LNK, SP, 4)
}

func Enter(size int) {
	put(PSH, LNK, SP, 4)
	put(PSH, FP, SP, 4)
	put(MOV, FP, 0, SP)
	put(SUBI, SP, SP, size)
}

func Return(size int) {
	put(MOV, SP, 0, FP)
	put(POP, FP, SP, 4)
	put(POP, LNK, SP, size+4)
	putBR(RET, LNK)
}

func Open() {
	var i = 0

	Curlev = 0
	Pc = 0

	for i < 32 {
		regs[i] = false
		i += 1
	}

}

func Close() {
	put(POP, LNK, SP, 4)
	putBR(RET, LNK)
}

func Decode() {
	var (
		c, i  int
		w, op uint32
	)

	fmt.Printf("\nEntry address: %#.8x\n", entry*4)
	i = 0
	for i < Pc {
		w = uint32(code[i])
		op = (w >> 26) & 0x3F
		fmt.Printf("%#.8x %#.8x %-4s ", i*4, w, mnemo[op])
		if op < MOVI {
			// c = register
			fmt.Printf("R%.2d, R%.2d, R%.2d\n", uint32((w>>22)&0x0F), uint32((w>>18)&0x0F), uint32(w&0x0F))
		} else if op < BEQ {
			// c = 18-bit signed constant or displacement
			c = int(w & 0x3FFFF)
			if c >= 0x20000 {
				c -= 0x40000
			}
			fmt.Printf("R%.2d, R%.2d, %#+.5x\n", uint32((w>>22)&0x0F), uint32((w>>18)&0x0F), c)
		} else {
			c = int(w & 0x3FFFFFF)

			if op == RET {
				// c = link register
				fmt.Printf("R%.2d\n", c)
			} else {
				// c = 26-bit signed offset
				if c >= 0x2000000 {
					c -= 0x4000000
				}
				fmt.Printf("%#+.6x\n", c*4)
			}
		}
		i += 1
	}
    fmt.Printf("\n%d bytes\n", Pc)
}

func Dump(n int) {
	for i := 0; i < n; i++ {
		fmt.Printf("%#.8x\n", uint32(code[i]))
	}
}

// "Load Module": execute compiled code's module body
func Load() {
	RISC.Load(code, Pc)
	RISC.Execute(entry * 4)
}

// Module body
func init() {

	BoolType = new(TypeDesc)
	BoolType.Form = Boolean
	BoolType.Size = 4

	IntType = new(TypeDesc)
	IntType.Form = Integer
	IntType.Size = 4

}
