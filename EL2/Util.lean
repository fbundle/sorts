import Std.Data
namespace EL2

namespace Util

def optionMap (xs: List α) (f: α → Option β): List β :=
  let rec loop (ys: Array β) (xs: List α) (f: α → Option β): Array β :=
    match xs with
      | [] => ys
      | x :: xs =>
        match f x with
          | none => ys
          | some y => loop (ys.push y) xs f

  (loop #[] xs f).toList

def optionMapAll (xs: List α) (f: α → Option β): Option (List β) :=
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
def applyAtmostOnce {α: Type} {β} (fs: List (α → Option β)) (x: α): Option β :=
  match fs with
    | [] => none
    | f :: fs =>
      match f x with
        | some y => some y
        | none => applyAtmostOnce fs x

structure ParseAllResult (α: Type) (β: Type) where
  remaining: List α
  items: List β
  deriving Repr

def ParseAllResult.ok (r: ParseAllResult α β): Bool := r.remaining.length = 0

partial def parseAll (parse: List α → Option (List α × β)) (tokens: List α): ParseAllResult α β :=
  let rec loop (items : Array β) (tokens: List α): List α × Array β :=
    match parse tokens with
      | none => (tokens, items)
      | some (tokens, item) =>  loop (items.push item) tokens

  let (remaining, items) := loop #[] tokens
  {remaining := remaining, items := items.toList}


def Map (α) (β) [BEq α] [Hashable α] := Std.HashMap α β

structure Counter (α: Type) where
  field: α
  count: Nat := 0

def Counter.with (counter: Counter α) (field: β): Counter β := {
  counter with
  field := field,
}

def Counter.next (counter: Counter α): Counter α := {
  counter with
  count := counter.count + 1,
}










end Util
end EL2
