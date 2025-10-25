namespace EL2.Form.Util

def optionMap (xs: List α) (f: α → Option β): List β :=
  let rec loop (ys: Array β) (xs: List α) (f: α → Option β): Array β :=
    match xs with
      | [] => ys
      | x :: xs =>
        match f x with
          | none => ys
          | some y => loop (ys.push y) xs f

  (loop #[] xs f).toList

def optionMap? (xs: List α) (f: α → Option β): Option (List β) :=
  let ys := optionMap xs f
  if ys.length ≠ xs.length then
    none
  else
    ys

partial def statefulMap (xs: List α) (state: State) (f: State → α → Option (State × β)): State × List β :=
  let rec loop (state: State) (ys: Array β) (listA: List α): State × Array β :=
    match listA with
      | [] => (state, ys)
      | x :: xs =>
        match f state x with
          | none => (state, ys)
          | some (state, b) => loop state (ys.push b) xs

  let (state, ys) := loop state #[] xs
  (state, ys.toList)

def statefulMap? (xs: List α) (state: State) (f: State → α → Option (State × β)) : Option (State × List β) :=
  let (state, ys) := statefulMap xs state f
  if ys.length ≠ xs.length then
    none
  else
    some (state, ys)

partial def applyAtMostOnce? {α: Type} {β: Type} (x: α) (fs: List (α → Option β)): Option β :=
  match fs with
    | [] => none
    | f :: fs =>
      match f x with
        | none => applyAtMostOnce? x fs
        | some y => some y

def chain (f: α → Option β) (g: β → Option γ) (a: α): Option γ := do
  let b ← f a
  let c ← g b
  pure c

end EL2.Form.Util
