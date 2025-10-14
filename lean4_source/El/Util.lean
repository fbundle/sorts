namespace Util

def applyAll (xs: List α) (f: α → Option β): Option (List β) :=
  let rec loop (ys: Array β) (xs: List α) (f: α → Option β): Option (Array β) :=
    match xs with
      | [] => some ys
      | x :: xs =>
        match f x with
          | none => none
          | some y => loop (ys.push y) xs f

  match loop #[] xs f with
    | none => none
    | some a => some a.toList

def applySome (xs: List α) (f: α → Option β): List β :=
  let rec loop (ys: Array β) (xs: List α) (f: α → Option β): Array β :=
    match xs with
      | [] => ys
      | x :: xs =>
        match f x with
          | none => loop ys xs f
          | some y => loop (ys.push y) xs f

  (loop #[] xs f).toList


end Util
