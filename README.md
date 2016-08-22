# Oberon-0-go
Port of Niklaus Wirth's Oberon-0 teaching language compiler to Go.

The code is as faithful a reproduction of the original source as possible. The resulting code is not idiomatic Go, but rather Oberon written in Go.

# Differences
The differences in the two languages were handled as follows:
  1. **VAR parameters**. In most cases the VAR parameter is used to pass an "empty" variable that will be "filled" by the procedure call. This was replaced with a pointer in the func declaration, and a dereference in the func body.
  2. **Sets**. A set is used to track the allocation of CPU registers: if it's in the set, it's allocated, otherwise it's free. The same functionality was provided using an array of bools.
  3. **REPEAT-UNTIL**. Go doesn't support a looping construct with a test at the end. A generic for loop was used, with an if test inside.
  4. **Mutable strings**. In Go strings are immutable. Instead, byte slices were used, which are mutable.
  5. **Module bodies**. In Oberon this is the code that falls between the BEGIN and END keywords after the declarations, and is executed when the module is loaded. In Go there is a function called init() in each module that does the same job.
  6. **Nested procedure definitions**. Go does not allow nested definitions. These were instead declared outside the function, and any variables referenced in an enclosing scope passed in by address.
  7. **Character ranges**. The large CASE statement in the lexer uses the ranges "A".."Z", "a".."z". Go supports expressions as cases in a switch statement, so the equivalent became case (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z')
  8. **INC()/DEC()**. Go has convenient and compact equivalents with += and -=

# Conveniences
A number of features in Go are much more convenient to use than Oberon, and the temptation proved too strong to keep the original source "100% pure":
  1. **Initialization syntax**. Go allows data structures to be initialized using a literal notation. This obviates the need for doing this at run time, for example in the keyword table in the scanner (OSS.go)
  2. **Logical instructions**. When constructing and deconstructing 32-bit RISC CPU instructions, Oberon uses arithmetic operations. The C-derived bit shifting (>>, <<) and masking (&, |) notation in Go is used instead, its familiarity making it less error-prone to write & debug.

# Inconveniences
Oberon and Go are not the same languages. But it is easy to fall into that trap of assuming they are, when trying to write one in the other.
   1. **Slices vs Arrays** Have to be careful here - Go's syntax is built with a bias towards slices. In the Oberon usage, there is no equivalent concept, and no direct application of slices. Extra care to _force_ the use of arrays and not slices.
   2. **Immutable Strings** Extra hoops with byte arrays, and the confusion with byte slices.
