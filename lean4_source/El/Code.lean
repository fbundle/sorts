import El.Form

structure Beta (α: Type) where
  cmd: α
  args: List α

structure Annot (α: Type) where
  name: String
  type: α

structure Binding (α: Type) where
  name: String
  value: α

structure Typeof (α: Type) where
  value: α

structure Inh (α: Type) where -- Inhabited
  value: α

structure Pi (α: Type) where
  params: List (Annot α)
  body: α

structure Ind (α: Type) where -- Inductive
  type: Annot α
  cons: List (Annot α)

inductive Code where
  | name: String → Code
  | beta: Beta Code → Code
  | annot: Annot Code → Code
  | binding: Binding Code → Code
  | typeof: Typeof Code → Code
  | inh: Inh Code → Code
  | pi: Pi Code → Code
  | ind: Ind Code → Code


abbrev getName := Form.getName
abbrev getList := Form.getName
abbrev Form := Form.Form

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
            -- TODO add more cases
            | _ =>
              let code ← parseBeta list
              pure (.beta code)

end
