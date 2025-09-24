(let
    {Nil   := (inh Nil Any_2)}
    {nil   := (inh nil Nil)}
    {Bool  := (inh Bool Any_2)}
    {True  := (inh True Bool)}
    {False := (inh False Bool)}

    {Nat   := (inductive {Nat : Any_2}
        {n0 : Nat}
        {succ : {Nat -> Nat}}
    )}

    {n1 := (succ n0)}
    {n2 := (succ n1)}
    {n3 := (succ n2)}

    {x := {n1 ⊕ n2}}
    {x := {n1 ⊗ n2 ⊗ n3}}

    {is_pos := {{x : Nat} => (match x
        {(succ z) => True}
        {n0       => False}
    )}}

    {must_pos := {{x : Nat} => (match x
        {(succ z) => x}
        {n0       => nil}
    )}}


    {_ := (inspect is_pos)}                    # resolved type as       Nat -> Bool
    {_ := (inspect must_pos)}                  # resolved type as       Nat -> (Nat ⊕ Nil)
                                               # better to resolve as   Π_{x: Nat} B(x) where B(x) = (type (must_pos x))

    {_ := (inspect (is_pos n0))}
    {_ := (inspect (is_pos n2))}

    (type Unit_0)
)






