(* Types test module *)
MODULE Types;
    CONST 
        size = 5;
    
    TYPE 
        String = ARRAY size OF INTEGER;
    
    VAR 
        a: String;
    
BEGIN
    a[0] := 30;
    a[1] := 31;
    a[2] := 32;
    Write(a[0]); Write(a[1]); Write(a[2])
END Types.