namespace EL2.Term.Util

def optionMap (xs: List α) (f: α → Option β): List β :=
  -- TODO - remove this
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

partial def statefulMap (xs: List α) (state: State) (f: State → α → State × β): State × List β :=
  let rec loop (ys: Array β) (xs: List α) (state: State): State × List β :=
    match xs with
      | [] => (state, ys.toList)
      | x :: xs =>
        let (state, y) := f state x
        loop (ys.push y) xs state

  loop #[] xs state


def statefulMap? (xs: List α) (state: State) (f: State → α → Option (State × β)) : Option (State × List β) :=
  let rec loop (ys: Array β) (xs: List α) (state: State): Option (State × List β) :=
    match xs with
      | [] => some (state, ys.toList)
      | x :: xs =>
        match f state x with
          | none => none
          | some (state, y) => loop (ys.push y) xs state
  loop #[] xs state

end EL2.Term.Util
