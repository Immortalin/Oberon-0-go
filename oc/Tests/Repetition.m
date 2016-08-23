(* While-statement Test Module *)
MODULE Repetition;
    VAR v: INTEGER;
    
BEGIN
    v := 0;
    WHILE v < 10 DO
        Write(v);
        v := v + 1
    END
END Repetition.