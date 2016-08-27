/*
RISC.go: The CPU emulator
08.08.2016 TSS
*/

package RISC

import (
	"fmt"
)

// in bytes
const (
	MemSize = 1024
)

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

var (
	ir      uint32
	n, z    bool
	R       [16]int
	M       [MemSize / 4]int
	ProgOrg int
)

func regDump() {
	var i int

	for i = 0; i < 16; i += 4 {
		fmt.Printf("R[%#.2d]: %#.8x  ", i, R[i])
		fmt.Printf("R[%#.2d]: %#.8x  ", i+1, R[i+1])
		fmt.Printf("R[%#.2d]: %#.8x  ", i+2, R[i+2])
		fmt.Printf("R[%#.2d]: %#.8x\n", i+3, R[i+3])
	}
	fmt.Printf("\n")
}

func memDump() {
	var i int

	for i = 0; i < (MemSize / 4); i += 1 {
		fmt.Printf("%#.3x %#.8x\n", i*4, uint32(M[i]))
	}

}

func Execute(start int) {
	var (
		opc          uint32
		a, b, c, nxt int
	)

	R[14] = 0
	R[15] = start + ProgOrg

Loop:
	for {
		nxt = R[15] + 4
		ir = uint32(M[R[15]/4])
		opc = (ir / 0x4000000) % 0x40
		a = int((ir / 0x400000) % 0x10)
		b = int((ir / 0x40000) % 0x10)
		c = int(ir % 0x40000)
		//fmt.Printf("Executing opcode: %#.2d \n", uint8(opc))
		if opc < MOVI {
			// F0 instruction: c = register
			c = R[uint8(ir)%0x10]
		} else if opc < BEQ {
			// F1 instruction: c = 18-bit signed constant
			// F2 instruction: c = 18-bit signed displacement
			c = int(ir % 0x40000)
			if c >= 0x20000 {
				c -= 0x40000
			}
		} else {
			// F3 instruction: c = 26-bit signed displacement
			c = int(ir % 0x4000000)
			if c >= 0x2000000 {
				c -= 0x4000000
			}
		}
		switch opc {
		case MOV, MOVI:
			if b < 0 {
				b = -b
				R[a] = c >> uint(b)
			} else {
				R[a] = c << uint(b)
			}
		case MVN, MVNI:
			if b < 0 {
				b = -b
				R[a] = -(c >> uint(b))
			} else {
				R[a] = -(c << uint(b))
			}
		case ADD, ADDI:
			R[a] = R[b] + c
		case SUB, SUBI:
			R[a] = R[b] - c
		case MUL, MULI:
			R[a] = R[b] * c
		case Div, DIVI:
			R[a] = R[b] / c
		case Mod, MODI:
			R[a] = R[b] % c
		case CMP, CMPI:
			z = (R[b] == c)
			n = (R[b] < c)
		case CHKI:
			if R[a] < 0 || R[a] >= c {
				R[a] = 0
			}
		case LDW:
			R[a] = M[(R[b]+c)/4]
            // Set flags on load (##BUGFIX)
            z = (R[a] == 0)
            n = (R[a] < 0)
		case POP:
			R[a] = M[(R[b])/4]
			R[b] += c
		case STW:
			M[(R[b]+c)/4] = R[a]
		case PSH:
			R[b] -= c
			M[(R[b])/4] = R[a]
		case RD:
			fmt.Scanf("%d", &R[a])
		case WRD:
			fmt.Printf("%d ", R[c])
		case WRH:
			fmt.Printf("%#x ", R[c])
		case WRL:
			fmt.Printf("\n")
		case BEQ:
			if z {
				nxt = R[15] + c*4
			}
		case BNE:
			if !z {
				nxt = R[15] + c*4
			}
		case BLT:
			if n {
				nxt = R[15] + c*4
			}
		case BGE:
			if !n {
				nxt = R[15] + c*4
			}
		case BLE:
			if z || n {
				nxt = R[15] + c*4
			}
		case BGT:
			if !z && !n {
				nxt = R[15] + c*4
			}
		case BR:
			nxt = R[15] + c*4
		case BSR:
			nxt = R[15] + c*4
			R[14] = R[15] + 4
		case RET:
			nxt = R[c&0x0F]
			if nxt == 0 {
				break Loop
			}
		}
		//regDump()
		R[15] = nxt
	}
	
}

func Load(code [256]int, len int) {
	var i int

	// "zero out" memory
	i = 0
	for i < MemSize/4 {
		M[i] = 0xFFFFFFFF
		i += 1
	}
	// copy and relocate memory image
	i = 0
	for i < len {
		M[i+(ProgOrg/4)] = code[i]
		i += 1
	}
}
