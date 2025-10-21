namespace EL2

inductive Term where
  | univ: (level: Int) → Term
  | var: (name: String) → Term
  | inh: (type: Term) → (cons: String) → (args: List Term) → Term
  | typ: (value: Term) → Term
  | list: (init: List Term) → (last: Term) → Term
  | bind: (name: String) → (value: Term) → Term
  | lam: (params: List (String × Term)) → (body: Term) → Term -- param is (name, type)
  | app: (cmd: Term) → (args: List Term) → Term
  | mat: (cond: Term) → (cases:  List (String × (List String) × Term)) → Term -- case is (patCmd, patArgs, value)
  deriving Repr, BEq -- BEq is computationally equal == DecidableEq is logical equal = and strictly stronger than ==

end EL2
