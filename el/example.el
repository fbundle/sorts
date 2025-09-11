#infix_left ⊕ ⊗
#infix_right => -> : , = :=

(let
    {Bool : U_0}
    {True : Bool}
    {False : Bool}
    {Nat : U_0}
    {0 : Nat}
    {succ : Nat -> Nat}
    {is_two : Nat -> Bool}

    {1 := (succ 0)}
    {2 := (succ 1)}
    {3 := (succ 2)}
    {is_two := {x => (match x
        (succ 1) True
        False
    )}}

    (is_two 3)
)