# python-like syntax using transpiler

# start of block patterns: let, match, (name:type) =>, (

let
   Nil = inh Any_2
   nil = inh Nil

   Bool = inh Any_2
   True = inh Bool
   False = inh Bool

   Nat = inh Any_2
   n0 = inh Nat
   succ = inh (Nat -> Nat)

   n1 = succ n0
   n2 = succ n1
   n3 = succ n2
   n4 = succ n3
   x = n1 ⊕ n2 ⊕ n3
   x = n1 ⊗ n2 ⊗ n3 ⊗ n4

   is_pos = (x: Nat) => match x with
      | succ z    => True
      | n0        => False

   must_pos = (x: Nat) => match x with
      | succ z    => True
      | n0        => nil


   print is_pos                    # resolved type as Nat -> Bool
   print must_pos                  # resolved type as Nat -> (Nat + Nil)
