(let
    Nil (inh Any_2)
    nil (inh Nil)
    Bool (inh Any_2)
    True (inh Bool)
    False (inh Bool)

    Nat (inh Any_2)
    n0 (inh Nat)
    succ (inh {Nat -> Nat})

    n1 (succ n0)
    n2 (succ n1)
    n3 (succ n2)
    n4 (succ n3)

    x {n1 ⊕ n2 ⊕ n3}
    x {n1 ⊗ n2 ⊗ n3 ⊗ n4}

    is_two (lambda x Nat (match x
        (exact n2)   True
                    False
    ))

    must_two (lambda x Nat (match x
        (exact n2)   n2
                     nil
    ))



    # (is_two n2)
    # is_two                  # is_two was successfully resolved type as Nat -> Bool
    must_two                  # is_two was successfully resolved type as Nat -> Bool
)






