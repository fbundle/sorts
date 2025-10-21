import EL2.Term
import EL2.Util

namespace EL2


class Context Ctx where
  set: Ctx → String → Term × Term × Int → Ctx
  get?: Ctx → String → Option (Term × Term × Int)

structure InferedTerm where
  term: Term
  type: Term
  level: Int

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

def Counter.dummyName (counter: Counter α): String := s!"dummy_{counter.count}"

partial def infer? [Repr Ctx] [Context Ctx] (reduce: Bool) (ctx: Ctx) (term: Term): Option (Ctx × InferedTerm) := do
  let isLam? (term: Term): Option (Lam Term) :=
    match term with
      | lam l => some l
      | _ => none
  none



end EL2
