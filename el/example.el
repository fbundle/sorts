(let
    Bool U_1                undef
    True Bool               undef
    False Bool              undef

    Nat U_1                 undef
    0 Nat                   undef
    succ {Nat -> Nat}       undef

    1 Nat (succ 0)
    2 Nat (succ 1)
    3 Nat (succ 2)
    4 Nat (succ 3)

    x Any_0 {1 ⊕ 2 ⊕ 3}
    x Any_0 {1 ⊗ 2 ⊗ 3 ⊗ 4}

    is_two {Nat -> Bool} {x => (match x
        (succ 1) True
        False
    )}

    add {Nat -> Nat -> Nat} {x => {y => (match y
        (succ @z) (succ ((add x) @z))
        x
    )}}

    # (is_two 2)
    ((add 2) 3)             # output (succ (succ (succ (succ (succ 0)))))
)





