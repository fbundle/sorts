(let Bool
    {Bool: U_0}
    {True: Bool}
    {False: Bool}

    {Nat: U_0}
    {0: Nat}
    {succ: {Nat -> Nat}}

    (def {1: Nat} (succ 0))
    (def {2: Nat} (succ 1))
    (def {3: Nat} (succ 2))

    {is_two: {Nat -> Bool} : {x => (match x
        (succ 1) True
        False
    )}}

    (is_two 3)
)