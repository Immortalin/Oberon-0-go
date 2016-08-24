(* Array Test Module *)
(* Expanded to handle const and var indices *)
MODULE Array;
    VAR 
        a: ARRAY 5 OF INTEGER;
        i: INTEGER;
BEGIN
    (* Constant indices *)
    a[0] := 7;
    a[1] := 3;
    Write(a[0]);
    Write(a[1]);
    
    (* Variable indices *)
    (* This doesnt crash but prints 4294967295 *)
    i := 1;
    a[i] := i;    
    Write(a[i]);
    
    (* This doesnt crash but prints 4294967295 *)
    i := 2;
    a[i] := i;    
    Write(a[i]);
    
    (* This crashes the compiler *)
    i := 0;
    a[i] := i;
    Write(a[i]);
    
    WriteLn
    
END Array.