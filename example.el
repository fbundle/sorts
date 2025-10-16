(inductive {Bool: U_2}
  {true: Bool}
  {false: Bool}
)

(inductive {Nat: U_2}
  {zero: Nat}
  {succ: {{_: Nat} => Nat}}
)

{one := (succ zero)}
{two := (succ one)}

{is_pos :=
  {{n: Nat} => (match n
    { zero   -> false}
    { (succ x) -> true}
  )}
}

(is_pos zero)
(is_pos two)

(inductive {{{T: U_2} => (List T)} : U_2}
  {nil : (List T)}
  {cons : {{init: (List T)} {tail : T} => (List T)}}
)


{x := 3}
{y := 4}

{x => y => z}

