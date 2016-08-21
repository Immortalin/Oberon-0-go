(* Array Test Module *)
MODULE Array;
    VAR a: ARRAY 2 OF INTEGER;
BEGIN
    a[0] := 7;
    a[1] := 3;
    Write(a[0]);
    Write(a[1])
END Array.