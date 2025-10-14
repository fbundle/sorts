namespace Util

partial def applyMany (xs: List α) (f: α → Option β): Option (List β) :=
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

end Util
