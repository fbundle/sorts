
import Lean.Data


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

partial def iterateAll (parse: List α → Option (List α × β)) (tokens: List α): List α × List β :=
  let rec loop (items : Array β) (tokens: List α): List α × Array β :=
    match parse tokens with
      | none => (tokens, items)
      | some (tokens, item) =>  loop (items.push item) tokens

  let (rest, items) := loop #[] tokens
  (rest, items.toList)


def Frame β := Lean.PersistentHashMap String β
def emptyFrame: Frame β := Lean.PersistentHashMap.empty
def Frame.set (f: Frame β) (key: String) (val: β): Frame β := f.insert key val
def Frame.get? (f: Frame β) (key: String): Option β := f.find? key

end Util
