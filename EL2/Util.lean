import Std.Data
namespace EL2

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

partial def ctxMap (listA: List α) (f: α → Option (Ctx × β)) (ctx: Ctx): Ctx × List β :=
  let rec loop (ctx: Ctx) (arrayB: Array β) (listA: List α): Ctx × Array β :=
    match listA with
      | [] => (ctx, arrayB)
      | a :: listA =>
        match f a with
          | none => (ctx, arrayB)
          | some (ctx, b) => loop ctx (arrayB.push b) listA

  let (ctx, arrayB) := loop ctx #[] listA
  (ctx, arrayB.toList)

def ctxMapAll (listA: List α) (f: α → Option (Ctx × β)) (ctx: Ctx): Option (Ctx × List β) :=
  let (ctx, listB) := ctxMap listA f ctx
  if listB.length ≠ listA.length then
    none
  else
    some (ctx, listB)









end Util
end EL2
