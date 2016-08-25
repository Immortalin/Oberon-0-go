(* BOOLEANs Test Module *)
MODULE Booleans;
    VAR b: BOOLEAN;
BEGIN
    b := FALSE;
    IF b THEN
        Write(1)
    ELSE
        Write(0)
    END;
    
    b := TRUE;
    IF b THEN
        Write(1)
    ELSE
        Write(0)
    END;
    
    b := (1 < 0);
    IF b THEN
        Write(1)
    ELSE
        Write(0)
    END;
    
    WriteLn
    
END Booleans.