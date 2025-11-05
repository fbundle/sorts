namespace EL2.Reducer

inductive Exp where
  -- typ 0 is type of small types: Nat, Pi, etc.
  -- typ 0 is at level 2
  -- typ N is at level N + 2
  -- small types are at level 1
  -- terms are at level 0
  | typ : (level: Nat) → Exp
  -- variable
  | var: (name: String) → Exp
  -- application
  | app: (cmd: Exp) → (arg: Exp) → Exp
  -- λ abstraction
  | lam: (name: String) → (body: Exp) → Exp
  -- let binding: let name: type := value
  | bnd: (name: String) → (value: Exp) → (body: Exp) → Exp
  -- inh - const
  | inh: (name: String) → (body: Exp) → Exp
  deriving Repr

end EL2.Reducer
