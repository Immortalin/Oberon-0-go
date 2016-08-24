/*
OSP.go: The Oberon-0 Parser
25.07.2016 TSS
VAR parameters replaced by pointer to parameter. Procedure call becomes
eg OSS.Get(&sym), and usage of sym inside Get is always as *sym
*/

package OSP

import (
	"OSG"
	"OSS"
	"bytes"
	"fmt"
)

const (
	wordSize = 4
)

var (
	sym                int
	topScope, universe OSG.Object
	guard              OSG.Object
	Dump               bool
)

func printObj(x OSG.Object) {
	for x != guard {
		fmt.Println(x)
		x = x.Next
	}
}

func newObj(obj *OSG.Object, class int) {
	var (
		n, x OSG.Object
	)

	x = topScope
	guard.Name = OSS.Id

	for !bytes.Equal(x.Next.Name, OSS.Id) {
		x = x.Next
	}
	if x.Next == guard {
		n = new(OSG.ObjDesc)
		// "n.name = OSS.Id"
		n.Name = make([]byte, 0, 16)
		n.Name = append(n.Name, OSS.Id...)
		n.Class = class
		n.Next = guard
		x.Next = n
		*obj = n
	} else {
		*obj = x.Next
		OSS.Mark("newObj: name already defined")
		fmt.Println(string(OSS.Id))
	}
}

func find(obj *OSG.Object) {
	var (
		s, x OSG.Object
	)

	s = topScope
	guard.Name = OSS.Id

	for {
		x = s.Next

		for !bytes.Equal(x.Name, OSS.Id) {
			x = x.Next
		}
		if x != guard {
			*obj = x
			break
		}
		if s == universe {
			*obj = x
			OSS.Mark("find: undefined name")
			break
		}
		s = s.Dsc
	}
}

func findField(obj *OSG.Object, list OSG.Object) {
	guard.Name = OSS.Id
	for !bytes.Equal(list.Name, OSS.Id) {
		list = list.Next
	}
	*obj = list
}

func isParam(obj OSG.Object) bool {
	return (obj.Class == OSG.Par) || (obj.Class == OSG.Var) && (obj.Val > 0)
}

func OpenScope() {
	var s OSG.Object

	s = new(OSG.ObjDesc)
	s.Class = OSG.Head
	s.Dsc = topScope
	s.Next = guard
	s.Name = make([]byte, 0, 16)
	topScope = s
}

func CloseScope() {
	topScope = topScope.Dsc
}

func selector(x *OSG.Item) {
	var (
		y   OSG.Item
		obj OSG.Object
	)

	for (sym == OSS.Lbrak) || (sym == OSS.Period) {
		if sym == OSS.Lbrak {
			OSS.Get(&sym)
			expression(&y)
			if x.Tp.Form == OSG.Array {
				OSG.Index(x, &y)
			} else {
				OSS.Mark("selector: bracketed expressions only allowed on arrays")
			}
			if sym == OSS.Rbrak {
				OSS.Get(&sym)
			} else {
				OSS.Mark("selector: ] expected")
			}
		} else {
			OSS.Get(&sym)
			if sym == OSS.Ident {
				if x.Tp.Form == OSG.Record {
					findField(&obj, x.Tp.Fields)
					OSS.Get(&sym)
					if obj != guard {
						OSG.Field(x, obj)
					} else {
						OSS.Mark("selector: undefined record field")
					}
				}
			} else {
				OSS.Mark("selector: identifier expected after .")
			}
		}
	}
}

func factor(x *OSG.Item) {
	var obj OSG.Object

	if sym < OSS.Lparen {
		OSS.Mark("factor: identifier expected")
		for {
			OSS.Get(&sym)
			if sym >= OSS.Lparen {
				break
			}
		}
	}

	if sym == OSS.Ident {
		find(&obj)
		OSS.Get(&sym)
		OSG.MakeItem(x, obj)
		selector(x)
	} else if sym == OSS.Number {
		OSG.MakeConstItem(x, OSG.IntType, OSS.Val)
		OSS.Get(&sym)
	} else if sym == OSS.Lparen {
		OSS.Get(&sym)
		expression(x)
		if sym == OSS.Rparen {
			OSS.Get(&sym)
		} else {
			OSS.Mark("factor: ) expected")
		}
	} else if sym == OSS.Not {
		OSS.Get(&sym)
		factor(x)
		OSG.Op1(OSS.Not, x)
	} else {
		OSS.Mark("factor: factor expected")
		OSG.MakeItem(x, guard)
	}
}

func term(x *OSG.Item) {
	var (
		y  OSG.Item
		op int
	)

	factor(x)

	for (sym >= OSS.Times) && (sym <= OSS.And) {
		op = sym
		OSS.Get(&sym)
		if op == OSS.And {
			OSG.Op1(op, x)
		}
		factor(&y)
		OSG.Op2(op, x, &y)
	}
}

func simpleExpression(x *OSG.Item) {
	var (
		y  OSG.Item
		op int
	)

	if sym == OSS.Plus {
		// leading +
		OSS.Get(&sym)
		term(x)
	} else if sym == OSS.Minus {
		// leading -
		OSS.Get(&sym)
		term(x)
		OSG.Op1(OSS.Minus, x)
	} else {
		term(x)
	}

	for (sym >= OSS.Plus) && (sym <= OSS.Or) {
		op = sym
		OSS.Get(&sym)
		if op == OSS.Or {
			OSG.Op1(op, x)
		}
		term(&y)
		OSG.Op2(op, x, &y)
	}
}

func expression(x *OSG.Item) {
	var (
		y  OSG.Item
		op int
	)

	simpleExpression(x)
	if (sym >= OSS.Eql) && (sym <= OSS.Gtr) {
		op = sym
		OSS.Get(&sym)
		simpleExpression(&y)
		OSG.Relation(op, x, &y)
	}
}

func parameter(fp *OSG.Object) {
	var x OSG.Item

	expression(&x)
	if isParam(*fp) {
		OSG.Parameter(&x, (*fp).Tp, (*fp).Class)
		*fp = (*fp).Next
	} else {
		OSS.Mark("parameter: too many parameters")
	}
}

func param(x *OSG.Item) {

	if sym == OSS.Lparen {
		OSS.Get(&sym)
	} else {
		OSS.Mark("StatSequence/param: ( expected")
	}

	expression(x)

	if sym == OSS.Rparen {
		OSS.Get(&sym)
	} else {
		OSS.Mark("StatSequence/param: ) expected")
	}
}

func statSequence() {
	var (
		par, obj OSG.Object
		x, y     OSG.Item
		L        int
	)

	for {
		obj = guard
		if sym < OSS.Ident {
			OSS.Mark("StatSequence: statement expected")
			for {
				OSS.Get(&sym)
				if sym >= OSS.Ident {
					break
				}
			}
		}
		if sym == OSS.Ident {
			find(&obj)
			OSS.Get(&sym)
			OSG.MakeItem(&x, obj)
			selector(&x)
			if sym == OSS.Becomes {
				OSS.Get(&sym)
				expression(&y)
				OSG.Store(&x, &y)
			} else if sym == OSS.Eql {
				OSS.Mark("StatSequence: = found instead of :=")
				OSS.Get(&sym)
				expression(&y)
			} else if x.Mode == OSG.Proc {
				par = obj.Dsc
				if sym == OSS.Lparen {
					OSS.Get(&sym)
					if sym == OSS.Rparen {
						OSS.Get(&sym)
					} else {
						for {
							parameter(&par)
							if sym == OSS.Comma {
								OSS.Get(&sym)
							} else if sym == OSS.Rparen {
								OSS.Get(&sym)
								break
							} else if sym >= OSS.Semicolon {
								break
							} else {
								OSS.Mark("StatSequence: ) or , expected in procedure call")
							}
						}
					}
				}
				if obj.Val < 0 {
					OSS.Mark("StatSequence: forward call of procedure")
				} else if !isParam(par) {
					OSG.Call(&x)
				} else {
					OSS.Mark("StatSequence: too few parameters in procedure call")
				}
			} else if x.Mode == OSG.SProc {
				if obj.Val <= 3 {
					param(&y)
				}
				OSG.IOCall(&x, &y)
			} else if obj.Class == OSG.Typ {
				OSS.Mark("StatSequence: illegal assignment")
			} else {
				OSS.Mark("StatSequence: statement expected")
			}
		} else if sym == OSS.If {
			OSS.Get(&sym)
			expression(&x)
			OSG.CJump(&x)
			if sym == OSS.Then {
				OSS.Get(&sym)
			} else {
				OSS.Mark("StatSequence: THEN expected in if-statement")
			}
			statSequence()
			L = 0
			for sym == OSS.Elsif {
				OSS.Get(&sym)
				OSG.FJump(&L)
				OSG.FixLink(x.A)
				expression(&x)
				OSG.CJump(&x)
				if sym == OSS.Then {
					OSS.Get(&sym)
				} else {
					OSS.Mark("StatSequence: THEN expected in ELSIF branch")
				}
				statSequence()
			}
			if sym == OSS.Else {
				OSS.Get(&sym)
				OSG.FJump(&L)
				OSG.FixLink(x.A)
				statSequence()
			} else {
				OSG.FixLink(x.A)
			}
			OSG.FixLink(L)
			if sym == OSS.End {
				OSS.Get(&sym)
			} else {
				OSS.Mark("StatSequence: END expected in if-statement")
			}
		} else if sym == OSS.While {
			OSS.Get(&sym)
			L = OSG.Pc
			expression(&x)
			OSG.CJump(&x)
			if sym == OSS.Do {
				OSS.Get(&sym)
			} else {
				OSS.Mark("StatSequence: DO expected in while-statement")
			}
			statSequence()
			OSG.BJump(L)
			OSG.FixLink(x.A)
			if sym == OSS.End {
				OSS.Get(&sym)
			} else {
				OSS.Mark("StatSequence: END expected in while-statement")
			}
		}
		if sym == OSS.Semicolon {
			OSS.Get(&sym)
		} else if (sym >= OSS.Semicolon) && (sym < OSS.If) || (sym >= OSS.Array) {
			break
		} else {
			OSS.Mark("StatSequence: semicolon expected at end of statement")
		}
	}
}

func identList(class int, first *OSG.Object) {
	var obj OSG.Object

	if sym == OSS.Ident {
		newObj(first, class)
		OSS.Get(&sym)
		for sym == OSS.Comma {
			OSS.Get(&sym)
			if sym == OSS.Ident {
				newObj(&obj, class)
				OSS.Get(&sym)
			} else {
				OSS.Mark("identList: Identifier expected")
			}
		}
		if sym == OSS.Colon {
			OSS.Get(&sym)
		} else {
			OSS.Mark("identList: : expected")
		}
	}
}

func Type(Tp *OSG.Type) {
	var (
		obj, first OSG.Object
		x          OSG.Item
		tp         OSG.Type
	)

	*Tp = OSG.IntType
	if (sym != OSS.Ident) && (sym < OSS.Array) {
		OSS.Mark("Type declaration: type declaration expected")
		for {
			OSS.Get(&sym)
			if (sym == OSS.Ident) || (sym >= OSS.Array) {
				break
			}
		}
	}
	if sym == OSS.Ident {
		find(&obj)
		OSS.Get(&sym)
		if obj.Class == OSG.Typ {
			*Tp = obj.Tp
		} else {
			OSS.Mark("Type declaration: unrecognized type")
		}
	} else if sym == OSS.Array {
		OSS.Get(&sym)
		expression(&x)
		if (x.Mode != OSG.Const) || (x.A < 0) {
			OSS.Mark("Type declaration: bad array size")
		}
		if sym == OSS.Of {
			OSS.Get(&sym)
		} else {
			OSS.Mark("Type declaration: OF expected")
		}
		Type(&tp)
		*Tp = new(OSG.TypeDesc)
		(*Tp).Form = OSG.Array
		(*Tp).Base = tp
		(*Tp).Len = x.A
		(*Tp).Size = (*Tp).Len * tp.Size
	} else if sym == OSS.Record {
		OSS.Get(&sym)
		(*Tp) = new(OSG.TypeDesc)
		(*Tp).Form = OSG.Record
		(*Tp).Size = 0
		OpenScope()
		for {
			if sym == OSS.Ident {
				identList(OSG.Fld, &first)
				Type(&tp)
				obj = first
				for obj != guard {
					obj.Tp = tp
					obj.Val = (*Tp).Size
					(*Tp).Size += obj.Tp.Size
					obj = obj.Next
				}
			}
			if sym == OSS.Semicolon {
				OSS.Get(&sym)
			} else if sym == OSS.Ident {
				OSS.Mark("Type declaration: ; expected")
			} else {
				break
			}
		}
		(*Tp).Fields = topScope.Next
		CloseScope()
		if sym == OSS.End {
			OSS.Get(&sym)
		} else {
			OSS.Mark("Type declaration: END expected")
		}
	} else {
		OSS.Mark("Type declaration: identifier, ARRAY or RECORD expected")
	}
}

func declarations(varsize *int) {
	var (
		obj, first OSG.Object
		x          OSG.Item
		tp         OSG.Type
	)

	if (sym < OSS.Const) && (sym != OSS.End) {
		OSS.Mark("declarations: declaration expected")
		for {
			OSS.Get(&sym)
			if (sym >= OSS.Const) || (sym == OSS.End) {
				break
			}
		}
	}

	for {
		if sym == OSS.Const {
			OSS.Get(&sym)
			for sym == OSS.Ident {
				newObj(&obj, OSG.Const)
				OSS.Get(&sym)
				if sym == OSS.Eql {
					OSS.Get(&sym)
				} else {
					OSS.Mark("Constant declaration: = expected")
				}
				expression(&x)
				if x.Mode == OSG.Const {
					obj.Val = x.A
					obj.Tp = x.Tp
				} else {
					OSS.Mark("Constant declaration: expression is not constant")
				}
				if sym == OSS.Semicolon {
					OSS.Get(&sym)
				} else {
					OSS.Mark("Constant declaration: ; expected")
				}
			}
		}
		if sym == OSS.Type {
			OSS.Get(&sym)
			for sym == OSS.Ident {
				newObj(&obj, OSG.Typ)
				OSS.Get(&sym)
				if sym == OSS.Eql {
					OSS.Get(&sym)
				} else {
					OSS.Mark("Type declaration: = expected")
				}
				Type(&obj.Tp)
				if sym == OSS.Semicolon {
					OSS.Get(&sym)
				} else {
					OSS.Mark("Type declaration: ; expected")
				}
			}
		}
		if sym == OSS.Var {
			OSS.Get(&sym)
			for sym == OSS.Ident {
				identList(OSG.Var, &first)
				Type(&tp)
				obj = first
				for obj != guard {
					obj.Tp = tp
					obj.Lev = OSG.Curlev
					*varsize += obj.Tp.Size
					obj.Val = -*varsize
					obj = obj.Next
				}
				if sym == OSS.Semicolon {
					OSS.Get(&sym)
				} else {
					OSS.Mark("Variable declaration: ; expected")
				}
			}
		}
		if (sym >= OSS.Const) && (sym <= OSS.Var) {
			OSS.Mark("declarations: order must be CONST TYPE VAR")
		} else {
			break
		}
	}
}

func fpSection(parblksize *int) {
	var (
		obj, first OSG.Object
		tp         OSG.Type
		parsize    int
	)

	if sym == OSS.Var {
		OSS.Get(&sym)
		identList(OSG.Par, &first)
	} else {
		identList(OSG.Var, &first)
	}
	if sym == OSS.Ident {
		find(&obj)
		OSS.Get(&sym)
		if obj.Class == OSG.Typ {
			tp = obj.Tp
		} else {
			OSS.Mark("FPSection: unknown parameter type")
			tp = OSG.IntType
		}
	} else {
		OSS.Mark("FPSection: type identifier expected")
		tp = OSG.IntType
	}
	if first.Class == OSG.Var {
		parsize = tp.Size
		if tp.Form >= OSG.Array {
			OSS.Mark("FPSection: array or struct params not supported")
		}
	} else {
		parsize = wordSize
	}
	obj = first
	for obj != guard {
		obj.Tp = tp
		*parblksize += parsize
		obj = obj.Next
	}
}

func procedureDecl() {
	const marksize = 8
	var (
		proc, obj              OSG.Object
		procid                 []byte
		locblksize, parblksize int
	)

	OSS.Get(&sym)
	if sym == OSS.Ident {
		procid = make([]byte, 0, 16)
		procid = append(procid, OSS.Id...)
		newObj(&proc, OSG.Proc)
		OSS.Get(&sym)
		parblksize = marksize
		OSG.IncLevel(1)
		OpenScope()
		proc.Val = -1
		if sym == OSS.Lparen {
			OSS.Get(&sym)
			if sym == OSS.Rparen {
				OSS.Get(&sym)
			} else {
				fpSection(&parblksize)
				for sym == OSS.Semicolon {
					OSS.Get(&sym)
					fpSection(&parblksize)
				}
				if sym == OSS.Rparen {
					OSS.Get(&sym)
				} else {
					OSS.Mark("ProcedureDecl: ) expected")
				}
			}
		} else if OSG.Curlev == 1 {
			//OSG.EnterCmd(procid)
		}
		//fmt.Printf("ProcedureDecl: parblksize = %d\n", parblksize)
		obj = topScope.Next
		locblksize = parblksize
		for obj != guard {
			obj.Lev = OSG.Curlev
			// Bug in Appendix text in ELSE clause
			// Correct version on P.76
			if obj.Class == OSG.Par {
				locblksize -= wordSize
			} else {
				locblksize -= obj.Tp.Size
			}
			obj.Val = locblksize
			obj = obj.Next
		}
		proc.Dsc = topScope.Next
		if sym == OSS.Semicolon {
			OSS.Get(&sym)
		} else {
			OSS.Mark("ProcedureDecl: ; expected")
		}
		locblksize = 0
		declarations(&locblksize)
		for sym == OSS.Procedure {
			procedureDecl()
			if sym == OSS.Semicolon {
				OSS.Get(&sym)
			} else {
				OSS.Mark("ProcedureDecl: ; expected after nested procedure declaration")
			}
		}
		proc.Val = OSG.Pc
		OSG.Enter(locblksize)
		if sym == OSS.Begin {
			OSS.Get(&sym)
			statSequence()
		}
		if sym == OSS.End {
			OSS.Get(&sym)
		} else {
			OSS.Mark("ProcedureDecl: END expected")
		}
		if sym == OSS.Ident {
			if !bytes.Equal(procid, OSS.Id) {
				OSS.Mark("ProcedureDecl: names must match")
			}
			OSS.Get(&sym)
		}
		OSG.Return(parblksize - marksize)
		CloseScope()
		OSG.IncLevel(-1)
	}
}

func Module() {
	var (
		modid   []byte
		varsize int
	)

	if sym == OSS.Module {
		OSS.Get(&sym)
		OSG.Open()
		OpenScope()
		varsize = 0
		if sym == OSS.Ident {
			modid = make([]byte, 0, 16)
			modid = append(modid, OSS.Id...)
			OSS.Get(&sym)
		} else {
			OSS.Mark("Module: identifier expected after MODULE")
		}
		if sym == OSS.Semicolon {
			OSS.Get(&sym)
		} else {
			OSS.Mark("Module: ; expected after identifier")
		}
		declarations(&varsize)
		for sym == OSS.Procedure {
			procedureDecl()
			if sym == OSS.Semicolon {
				OSS.Get(&sym)
			} else {
				OSS.Mark("Module: ; expected to end procedure declaration")
			}
		}
		OSG.Header(varsize)
		if sym == OSS.Begin {
			OSS.Get(&sym)
			statSequence()
		}
		if sym == OSS.End {
			OSS.Get(&sym)
		} else {
			OSS.Mark("Module: END expected")
		}
		if sym == OSS.Ident {
			if !(bytes.Compare(OSS.Id, modid) == 0) {
				OSS.Mark("Module: module names do not match")
				fmt.Println(string(OSS.Id))
				fmt.Println(string(modid))
			}
			OSS.Get(&sym)
		} else {
			OSS.Mark("Module: identifier expected after END")
		}
		if sym != OSS.Period {
			OSS.Mark("Module: . expected after identifier")
		}
		CloseScope()
		if OSS.Error == false {
			fmt.Printf("Compile successful\n")
			OSG.Close()
			if Dump {
				OSG.Decode()
			}
			OSG.Load()
		}
	} else {
		OSS.Mark("Module: MODULE expected")
	}
}

func Compile(fname string) {
	OSS.Init(fname)
	OSS.Get(&sym)
	Module()
}

func enter(cl int, n int, name string, tp OSG.Type) {
	var obj OSG.Object

	obj = new(OSG.ObjDesc)
	obj.Class = cl
	obj.Val = n
	obj.Name = make([]byte, 0, 16)
	obj.Name = append(obj.Name, name...)
	obj.Tp = tp
	obj.Dsc = nil
	obj.Next = topScope.Next
	topScope.Next = obj
}

// Module body
func init() {
	fmt.Printf("Oberon-0 Compiler\n")

	guard = new(OSG.ObjDesc)
	guard.Class = OSG.Var
	guard.Tp = OSG.IntType
	guard.Val = 0

	topScope = nil
	OpenScope()

	enter(OSG.Typ, 1, "BOOLEAN", OSG.BoolType)
	enter(OSG.Typ, 2, "INTEGER", OSG.IntType)
	enter(OSG.Const, 1, "TRUE", OSG.BoolType)
	enter(OSG.Const, 0, "FALSE", OSG.BoolType)
	enter(OSG.SProc, 1, "Read", nil)
	enter(OSG.SProc, 2, "Write", nil)
	enter(OSG.SProc, 3, "WriteHex", nil)
	enter(OSG.SProc, 4, "WriteLn", nil)

	universe = topScope
}
