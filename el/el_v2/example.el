(let Bool
    {Bool: U_0}
    {True: Bool}
    {False: Bool}

    {Nat: U_0}
    {0: Nat}
    {succ: {Nat -> Nat}}

    {1 : Nat} {1 := (succ 0)}
    {2 : Nat} {2 := (succ 1)}
    {3 : Nat} {3 := (succ 2)}

    {is_two: {Nat -> Bool}}
    {is_two := {x => (match x
        (succ 1) True
        False
    )}}

    (is_two 3)
)