# example of EL code

(succ 0)
(succ (succ 0))

(let
	{x := 0}
	x
)

(add 3 5)

(let
	{Nat := (* Any_2)}
	{0 := (* Nat)}
	{NatToNat := {{_: Nat} => Nat}}		# NatToNat is Nat -> Nat
	{succ := (* NatToNat)}

	(succ (succ 0))
)

(let
	{NatPair := {{_: Nat} × Nat} }		# NatPair is Nat × Nat
	{x := (* NatPair)}
	x
)


