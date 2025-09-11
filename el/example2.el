(let
    Bool U_0                undef
    True Bool               undef
    False Bool              undef

    Nat U_0                 undef
    0 Nat                   undef
    succ {Nat -> Nat}       undef

    1 Nat (succ 0)
    2 Nat (succ 1)
    3 Nat (succ 2)

    3
)