namespace Util

def optionMap (xs: List α) (f: α → Option β): List β :=
  let rec loop (ys: Array β) (xs: List α) (f: α → Option β): Array β :=
    match xs with
      | [] => ys
      | x :: xs =>
        match f x with
          | none => loop ys xs f
          | some y => loop (ys.push y) xs f

  (loop #[] xs f).toList

def optionMapAll (xs: List α) (f: α → Option β): Option (List β) :=
  let ys := optionMap xs f
  if ys.length ≠ xs.length then
    none
  else
    ys

def applyOnce {α: Type} {β} (fs: List (α → Option β)) (x: α): Option β :=
  match fs with
    | [] => none
    | f :: fs =>
      match f x with
        | some y => some y
        | none => applyOnce fs x

end Util
