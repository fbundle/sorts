import EL2.Core.CoreV2

namespace EL2.Term
open EL2.Core

inductive Term where
  | typ: (level: Nat) → Term
  | var: (name: String) → Term
  | app: (cmd: Term) → (arg: Term) → Term
  | bnd: (name: String) → (value: Term) → (type: Term) → (body: Term) → Term
  | ann: (term: Term) → (type: Term) → Term
  | lam: (name: String) → (body: Term) → Term
  | pi:  (name: String) → (type: Term) → (body: Term) → Term
  | ind:
    (name: String) →
    (params: List (String × Term)) →                        -- type params
    (cons: List (String × List (String × Term) × Term)) →   -- constructors
    (body: Term) →
    Term
  deriving Nonempty




partial def chain (init: List (String × Term)) (tail: Term): Term :=
  match init with
    | [] => tail
    | (name, type) :: rest =>
        Term.pi name type (chain rest tail)

partial def curry (cmd: Term) (args: List Term): Term :=
  match args with
    | [] => cmd
    | arg :: rest =>
      curry (Term.app cmd arg) rest

def scott
  (name: String)
  (params: List (String × Term))
  (cons: List (String × List (String × Term) × Term))
  (body: Term): Term :=
    -- Scott encoding
    -- ind works like bnd - it binds type name, and constructor name
    -- then give body
    sorry

partial def Term.toExp (term: Term): Exp :=
  match term with
    | Term.typ level => Exp.typ level
    | Term.var name => Exp.var name
    | Term.app cmd arg => Exp.app cmd.toExp arg.toExp
    | Term.bnd name value type body => Exp.bnd name value.toExp type.toExp body.toExp
    | Term.ann term type => Exp.ann term.toExp type.toExp
    | Term.lam name body => Exp.lam name body.toExp
    | Term.pi name type body => Exp.pi name type.toExp body.toExp
    | Term.ind name params cons body => (scott name params cons body).toExp



def test1: Term := -- Nat and Vec
  id
  $ Term.ind "Nat" [] [
    ("zero", [], Term.var "Nat"),
    ("succ", [("prev", Term.var "Nat")], Term.var "Nat"),
  ]
  $ Term.ind "Vec0" [("n", Term.var "Nat"), ("T", Term.typ 0)] [
    (
      "nil",
      [("T", Term.typ 0)],
      curry (Term.var "Vec0") [Term.var "zero", Term.var "T"],
    ),
    (
      "push",
      [
        ("n", Term.var "Nat"), ("T", Term.typ 0),
        ("v", curry (Term.var "Vec0") [Term.var "n", Term.var "T"]),
        ("x", Term.var "T"),
      ],
      curry (Term.var "Vec0") [curry (Term.var "succ") [Term.var "n"], Term.var "T"],
    ),
  ]
  $ Term.bnd "one" (curry (Term.var "succ") [Term.var "zero"]) (Term.var "Nat")
  $ Term.bnd "empty_vec" (curry (Term.var "nil") [Term.var "Nat"]) (curry (Term.var "Vec0") [Term.var "zero", Term.var "Nat"])
  $ Term.bnd "single_vec" (curry (Term.var "push") [Term.var "empty_vec"]) (curry (Term.var "Vec0") [Term.var "one", Term.var "Nat"])
  $ Term.typ 0

#eval typeCheck? test1.toExp (Term.typ 1).toExp

end EL2.Term
