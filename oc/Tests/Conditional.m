(* If-statement Test Module *)
MODULE Conditional;
    VAR v: INTEGER;
    
BEGIN
    (* Single IF *)
    v := 2;
    IF v > 1 THEN
        Write(0);
    END;
    
    (* IF-ELSE *)
    v := 1;
    IF v > 1 THEN
        Write(1)
    ELSE
        Write(2)
    END;
    
    (* IF-ELSIF-ELSE *)
    v := 2;
    IF v > 3 THEN
        Write(10)
    ELSIF v > 2 THEN
        Write(20)
    ELSIF v > 1 THEN
        Write(30)
    ELSE
        Write(40)
    END;
    
    (* IF-ELSIF only *)
    v := 5;
    IF v = 4 THEN
        Write(100)
    ELSIF v = 5 THEN
        Write(200)
    END
 
END Conditional.