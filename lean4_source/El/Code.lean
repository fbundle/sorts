import El.Form

namespace Code

structure Beta (α: Type) where
  cmd: α
  args: List α
  deriving Repr

structure Annot (α: Type) where
  name: String
  type: α
  deriving Repr

structure Binding (α: Type) where
  name: String
  value: α
  deriving Repr

structure Typeof (α: Type) where
  value: α
  deriving Repr

structure Inh (α: Type) where -- Inhabited
  type: α
  deriving Repr

structure Pi (α: Type) where -- Pi or Lambda
  params: List (Annot α)
  body: α
  deriving Repr

inductive Code where
  | name: String → Code
  | beta: Beta Code → Code
  | annot: Annot Code → Code
  | binding: Binding Code → Code
  | typeof: Typeof Code → Code
  | inh: Inh Code → Code
  | pi: Pi Code → Code
  deriving Repr


abbrev getName := Form.getName
abbrev getList := Form.getName
abbrev Form := Form.Form

private partial def applyMany (xs: List α) (f: α → Option β): Option (List β) :=
  let rec loop (ys: Array β) (xs: List α) (f: α → Option β): Option (Array β) :=
    match xs with
      | [] => some #[]
      | x :: xs =>
        match f x with
          | none => none
          | some y => loop (ys.push y) xs f

  match loop #[] xs f with
    | none => none
    | some a => some a.toList

mutual

private partial def parseBeta (list: List Form): Option (Beta Code) := do
  match list with
    | [] => none
    | x :: xs =>
      let cmd ← parse x
      let args ← applyMany xs parse
      pure {cmd := cmd, args := args}

private partial def parseAnnot (list: List Form): Option (Annot Code) := do
  let nameForm ← list[0]?
  let name ← getName nameForm
  let typeForm ← list[1]?
  let type ← parse typeForm
  pure {name := name, type := type}

private partial def parseBinding (list: List Form): Option (Binding Code) := do
  let nameForm ← list[0]?
  let name ← getName nameForm
  let valueForm ← list[1]?
  let value ← parse valueForm
  pure {name := name, value := value}

private partial def parseTypeof (list: List Form): Option (Typeof Code) := do
  let valueForm ← list[0]?
  let value ← parse valueForm
  pure {value := value}

private partial def parseInh (list: List Form): Option (Inh Code) := do
  let typeForm ← list[0]?
  let type ← parse typeForm
  pure {type := type}

private partial def parsePi (list: List Form): Option (Pi Code) := do
  if list.length = 0 then
    none
  else
    let paramForms := list.extract 0 (list.length-1)
    let params ← applyMany paramForms ((λ form =>
      match form with
        | .name name => none
        | .list list => parseAnnot list
    ): Form → Option (Annot Code))
    let bodyForm ← list[list.length-1]?
    let body ← parse bodyForm
    pure {params := params, body := body}

partial def parse (form: Form): Option Code := do
  match form with
    | .name name => pure (Code.name name)
    | .list list =>
      match list with
        | [] => none
        | head :: rest =>
          match head with
            | .name ":" =>
              let code ← parseAnnot rest
              pure (.annot code)
            | .name ":=" =>
              let code ← parseBinding rest
              pure (.binding code)
            | .name "&" =>
              let code ← parseTypeof rest
              pure (.typeof code)
            | .name "*" =>
              let code ← parseInh rest
              pure (.inh code)
            | .name "=>" =>
              let code ← parsePi rest
              pure (.pi code)
            -- TODO add more cases
            | _ =>
              let code ← parseBeta list
              pure (.beta code)
end



def _example := "
  (:= Nat (* U_2))
  (:= 0 (* Nat))
  (:= succ (* (-> Nat)))

  (:= 1 (succ 0))
  (:= 2 (succ 0))
  (:= 3 (succ 0))
  (:= 4 (succ 0))
  (:= x 3)
  (:= y 4)

  (+ x y)
"

#eval _example

#eval Form.defaultParser.tokenize _example

#eval (Form.defaultParseAll _example).get!




end Code
