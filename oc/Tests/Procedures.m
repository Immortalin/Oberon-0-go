(* Procedures Test Module *)
MODULE Procedures;
    VAR v: INTEGER;
    
    PROCEDURE p1;
    BEGIN
        v := 1
    END p1;
    
    PROCEDURE p2(x: INTEGER);
    BEGIN
        Write(x)
    END p2;
    
    PROCEDURE p3(x: INTEGER; y: INTEGER);
    BEGIN
        Write(x+y)
    END p3;
    
    PROCEDURE p4(VAR x: INTEGER);
    BEGIN
        x := x + 1
    END p4;
    
    PROCEDURE p5;
        VAR k: INTEGER;
    BEGIN
        k := 0;
        Write(k);
        p4(k); Write(k);
        p4(k); Write(k)
    END p5;
 
BEGIN
    p1;
    Write(v);
    p2(104);
    p3(722, 3);
    
    v := 0;
    (* Should print 1 *)
    p4(v); Write(v);
    (* Should print 2 *)
    p4(v); Write(v);
    
    p5
END Procedures.