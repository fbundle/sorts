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

end EL2.Form.Util
