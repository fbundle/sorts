namespace EL2

structure Ann (α: Type) where -- (2: Nat)
  name: String
  type: α
  deriving Repr, BEq

structure Case (α: Type) where
  patCmd: String
  patArgs : List String
  value: α
  deriving Repr, BEq

inductive Term where
  | inh: (method: String) → (values: List Term) → (type: Term) → Term
  | infer: (term: Term) → Term -- type of
  | univ: (level: Int) → Term
  | var: (name: String) → Term
  | list: (init: List Term) → (last: Term) → Term
  | bind: (name: String) → (value: Term) → Term
  | typ: (name: String) → (args: List Term) → (level: Int) → Term
  | app: (cmd: Term) → (args: List Term) → Term
  | lam: (params: List (Ann Term)) → (body: Term) → Term
  | mat: (cond: Term) → (cases:  List (Case Term)) → Term
  deriving Repr, BEq -- BEq is computationally equal == DecidableEq is logical equal = and strictly stronger than ==

end EL2
