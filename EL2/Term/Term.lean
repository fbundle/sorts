import EL2.Core.CoreV2

namespace EL2.Term
open EL2.Core

inductive T α where
  | typ: (level: Nat) → T α
  | var: (name: String) → T α
  | app: (cmd: α) → (arg: α) → T α
  | bnd: (name: String) → (value: α) → (type: α) → (body: α) → T α
  | ann: (term: α) → (type: α) → T α
  | lam: (name: String) → (body: α) → T α
  | pi:  (name: String) → (type: α) → (body: α) → T α
  deriving Nonempty

def T.toExp (t: T α) (f: α → Exp) : Exp :=
  match t with
    | T.typ level => Exp.typ level
    | T.var name => Exp.var name
    | T.app cmd arg => Exp.app (f cmd) (f arg)
    | T.bnd name value type body => Exp.bnd name (f value) (f type) (f body)
    | T.ann term type => Exp.ann (f term) (f type)
    | T.lam name body => Exp.lam name (f body)
    | T.pi name type body => Exp.pi name (f type) (f body)


inductive Term where
  | t: (t: T Term) → Term
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
        Term.t $ T.pi name type (chain rest tail)

partial def curry (cmd: Term) (args: List Term): Term :=
  match args with
    | [] => cmd
    | arg :: rest =>
      curry (Term.t $ T.app cmd arg) rest

partial def Term.toExp (term: Term): Exp :=
  match term with
    | Term.t tt => tt.toExp Term.toExp
    | Term.ind name params cons body =>
      -- Scott encoding
      sorry


def test1 := -- Nat and Vec
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
  $ Term.ann (Term.var "single_vec") (curry (Term.var "Vec0") [Term.var "one", Term.var "Nat"])

#eval (scott 0 test1)

end EL2.Term
