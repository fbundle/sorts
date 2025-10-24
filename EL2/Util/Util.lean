import Std.Data
namespace EL2.Util

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










end EL2.Util
