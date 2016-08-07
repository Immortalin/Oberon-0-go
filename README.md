# Oberon-0-go
Port of Niklaus Wirth's Oberon-0 teaching language compiler to Go.

The code is as faithful a reproduction of the original source as possible. The resulting code is not idiomatic Go, but rather Oberon written in Go.

# Differences
The differences in the two languages were handled as follows:
  1. VAR parameters. In most cases the VAR parameter is used to pass an "empty" variable that will be "filled" by the procedure call. This was replaced with a pointer in the func declaration, and a dereference in the func body.
  2. Sets. A set is used to track the allocation of CPU registers: if it's in the set, it's allocated, otherwise it's free. The same functionality was provided using an array of bools.
  3. REPEAT-UNTIL. Go doesn't support a looping construct with a test at the end. A generic for loop was used, with an if test inside.
  4. Mutable strings. In Go strings are immutable. Instead, byte slices were used, which are mutable.
  5. Module bodies. In Oberon this is the code that falls between the BEGIN and END keywords after the declarations. In Go there is a function called init() in each module that does the same job.
  6. Nested procedure definitions. Go does not allow nested definitions. These were instead declared outside the function, and any variables referenced in an enclosing scope passed in by address.
  7. Character ranges. The large CASE statement in the lexer uses the ranges "A".."Z", "a".."z". Go supports expressions as cases in a switch statement, so the equivalent became case (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z')
