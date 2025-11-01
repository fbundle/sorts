import EL2.Core.CoreV2

namespace EL2.Term
open EL2.Core

inductive Term where
  | typ : (level: Nat) → Term
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

partial def scott (k: Nat) (term: Term): Term :=
  -- Scott encoding - turn inductive type into pi
  match term with
    | Term.ind indName params cons body =>
      let R := Term.var s!"_{k}"  -- fresh result type variable

      -- Scott-encode each constructor type
      let encCons := cons.map (fun (cName, args, cRetType) =>
        (cName, chain args cRetType)  -- constructor args -> constructor return type
      )

      -- Chain constructors into a single Pi to R (this is the Scott type body)
      let scottType := chain params (chain encCons R)

      -- Fold right: bind constructors first
      let withConstructors := encCons.foldr (fun (cName, cType) acc =>
        Term.bnd cName cType (Term.typ 0) acc
      ) body

      -- Finally bind the inductive type itself
      Term.bnd indName scottType (Term.typ 0) withConstructors

    | _ => term

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
