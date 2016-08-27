(* Types test module *)
MODULE Record;
    CONST 
        size = 5;
    
    TYPE 
        Item = RECORD
            A: INTEGER;
            b: BOOLEAN;
        END;
    
    VAR 
        a: Item;
    
BEGIN
    a.A := 1;
    a.b := TRUE;
    Write(a.A);
    IF a.b THEN
        Write(1)
    ELSE
        Write(0)
    END
END Record.