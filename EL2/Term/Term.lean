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
    (univ: Term) →                                           -- universe
    (name: String) →                                         -- type name
    (params: List (String × Term)) →                        -- type params
    (cons: List (String × List (String × Term) × List Term)) →   -- constructors
    (body: Term) →
    Term
  deriving Nonempty


partial def chainLam (init: List String) (tail: Term): Term :=
  match init with
    | [] => tail
    | name :: rest =>
      Term.lam name (chainLam rest tail)

partial def chainPi (init: List (String × Term)) (tail: Term): Term :=
  match init with
    | [] => tail
    | (name, type) :: rest =>
      Term.pi name type (chainPi rest tail)

partial def chain (init: List (String × Term)) (tail: Term × Term): (Term × Term) :=
  match init with
    | [] => tail
    | (name, type) :: rest =>
      let (restValue, restType) := chain rest tail
      (Term.lam name restValue, Term.pi name type restType)

partial def curry (cmd: Term) (args: List Term): Term :=
  match args with
    | [] => cmd
    | arg :: rest =>
      curry (Term.app cmd arg) rest

partial def chainBind (init: List (String × Term × Term)) (tail: Term): Term :=
  match init with
    | [] => tail
    | (name, type, value) :: rest =>
      Term.bnd name type value (chainBind rest tail)

def scott
  (univ: Term)
  (name: String)
  (params: List (String × Term))
  (cons: List (String × List (String × Term) × List Term))
  (body: Term): Term :=
    -- Scott encoding for inductive type
    -- ind works like bnd - it binds type name, and constructor name
    -- then give body
    -- inductive T (p1: T1) (p2: T2) ... (pN: TN): Univ where
    --    | cons1: (c11: V11) -> (c21: V21) ... -> (T x11 x21 ...)
    --    | ...
    --    | consM: (cM1: VM1) -> (cM2: VM2) ... -> (T x1M x2M ...)
    -- params = (p1: T1) (p2: T2) ... (pN: TN)
    -- cons[1] := cons1 ((c11: V11) -> (c21: V21) ...) (x11 x21 ...)

  let consNameList := cons.map (λ (consName, _, _) => consName)
  let consParamList := cons.map (λ (_, consParams, _) => consParams)

  let R := Term.var "R"
  let T := chainPi (
    ("R", univ) :: List.zip consNameList (consParamList.map (chainPi · R))
  ) R

  let Tfunc := chainPi params T



  let consNameTermTypeList := cons.map (λ (consName, consParams, consBody) =>
    -- x is of type Tfunc
    let x := chainLam (params.map (Prod.fst))
      (chainLam ("r" :: consNameList) (Term.var consName)) -- pick one, e.g. λ a₁ a₂ a₃ => a₁
    -- consTerm1 := λ c11 λ c21 ... (T x11 x21 ...)
    let consTerm := chainLam (consParams.map (Prod.fst))
      -- need anotated type since our core doesn't support untyped app
      (curry (Term.ann x Tfunc) consBody)

    -- type of consTerm
    let consType := chainPi consParams T
    (consName, consTerm, consType)
  )

  let consBind := chainBind consNameTermTypeList body
  let typeBind := Term.bnd name Tfunc univ consBind

  typeBind



partial def Term.toExp (term: Term): Exp :=
  match term with
    | Term.typ level => Exp.typ level
    | Term.var name => Exp.var name
    | Term.app cmd arg => Exp.app cmd.toExp arg.toExp
    | Term.bnd name value type body => Exp.bnd name value.toExp type.toExp body.toExp
    | Term.ann term type => Exp.ann term.toExp type.toExp
    | Term.lam name body => Exp.lam name body.toExp
    | Term.pi name type body => Exp.pi name type.toExp body.toExp
    | Term.ind univ name params cons body => (scott univ name params cons body).toExp



def test1: Term := -- Nat and Vec
  id
  $ Term.ind (Term.typ 0) "Nat" [] [
    ("zero", [], []),
    ("succ", [("prev", Term.var "Nat")], []),
  ]
  $ Term.ind (Term.typ 0) "Vec0" [("n", Term.var "Nat"), ("T", Term.typ 0)] [
    (
      "nil",
      [("T", Term.typ 0)],
      [Term.var "zero", Term.var "T"],
    ),
    (
      "push",
      [
        ("n", Term.var "Nat"), ("T", Term.typ 0),
        ("v", curry (Term.var "Vec0") [Term.var "n", Term.var "T"]),
        ("x", Term.var "T"),
      ],
      [curry (Term.var "succ") [Term.var "n"], Term.var "T"],
    ),
  ]
  $ Term.bnd "one" (curry (Term.var "succ") [Term.var "zero"]) (Term.var "Nat")
  $ Term.bnd "empty_vec" (curry (Term.var "nil") [Term.var "Nat"]) (curry (Term.var "Vec0") [Term.var "zero", Term.var "Nat"])
  $ Term.bnd "single_vec" (curry (Term.var "push") [Term.var "empty_vec"]) (curry (Term.var "Vec0") [Term.var "one", Term.var "Nat"])
  $ Term.typ 0

#eval test1.toExp
#eval typeCheck? test1.toExp (Term.typ 1).toExp


end EL2.Term
