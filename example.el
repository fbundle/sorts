(let
    {Nil   := (inh Any_2)}
    {nil   := (inh Nil)}
    {Bool  := (inh Any_2)}
    {True  := (inh Bool)}
    {False := (inh Bool)}

    {Nat   := (inh Any_2)}
    {n0    := (inh Nat)}
    {succ  := (inh {Nat -> Nat})}

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

    Unit_0
)






