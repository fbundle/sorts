package sorts

/*
// Inductive - inductive type
type Inductive interface {
	Sort
	Iter(yield func(name form.Name, constr func([]Sort) Inductive) bool)
}
*/

// Dept - represent a type B(x) depends on Sort x
// Dept is not a type/sort, it is a family of types indexed by A
type Dept func(Sort) Sort // Lambda, Match (if we consider match is kinda function) match (x: A) | case ...

/*
must_pos = lambda (x: Nat)
  match x with
	 | succ z    => x
	 | n0        => nil

must_pos_type = lambda (x: Nat)
	match x with
	 | succ z    => Nat
	 | n0        => Nil

in this example must_pos is of type Π_{x: Nat} must_pos_type(x)
from current type checking, we know that this is a subtype of Nat -> (Nat ⊕ Nil)
*/

// Pi - dependent function type Π_{x: A} B(x)
type Pi struct {
	H Name
	A Sort
	B Dept
}

func (s Pi) sortAttr(a SortAttr) sortAttr {
	return sortAttr{
		form: List{s.H, a.Form(s.A), a.Form(s.B)},
	}
}

/*
must_pos = lambda (x: Nat)
  match x with
	 | succ z    => (x ⊗ x)
	 | n0        => (x ⊗ nil)

must_pos_type = lambda (x: Nat)
	match x with
	 | succ z    => (Nat ⊗ Nat)
	 | n0        => (Nat ⊗ Nil)

in this example, must_pos is of type Σ_{x: Nat} must_pos_type(x)
from type checking we know that this is a subtype of (Nat ⊗ (Nat ⊕ Nil))
*/

// Sigma - dependent pair type Σ_{x: A} B(x)
type Sigma struct {
	H Name
	A Sort
	B Dept
}

func (s Sigma) sortAttr(a SortAttr) sortAttr {
	return sortAttr{
		form: List{s.H, a.Form(s.A), a.Form(s.B)},
	}
}
