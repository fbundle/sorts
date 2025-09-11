#infix_left ⊕ ⊗
#infix_right => -> : , = :=

(let
    {Bool : U_0}
    {True : Bool}
    {False : Bool}

    {Nat : U_0}
    {0 : Nat}
    {succ : {Nat -> Nat}}

    {1 : Nat}
    {2 : Nat}
    {3 : Nat}
    {1 := (succ 0)}
    {2 := (succ 1)}
    {3 := (succ 2)}

    {is_two : {Nat -> Bool}}
    {is_two := {x => (match x
        (succ 1) True
        False
    )}}

    {add: {Nat -> Nat -> Nat}}
    {add := {x => y => (match x
        (succ z) (add z (succ y))   # if x is succ z then add z (succ y)
        y                           # else y
    )}}

    (is_two 3)
)