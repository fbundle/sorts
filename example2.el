(let
    Bool Any_2                undef
    True Bool               undef
    False Bool              undef

    Nat Any_2                 undef
    n0 Nat                  undef
    succ {Nat -> Nat}       undef

    n1 Nat (succ n0)
    n2 Nat (succ n1)
    n3 Nat (succ n2)
    n4 Nat (succ n3)

    x Any_1 {n1 ⊕ n2 ⊕ n3}
    x Any_1 {n1 ⊗ n2 ⊗ n3 ⊗ n4}

    is_two {Nat -> Bool} {x => (match x
        (exact (succ n1))   True
                            False
    )}

    add {Nat -> Nat -> Nat} {x => y => (match y
        (succ z)    (succ ((add x) z))
                    x
    )}

    # (is_two n2)
    (add n2 n3)             # output (succ (succ (succ (succ (succ n0)))))
)






