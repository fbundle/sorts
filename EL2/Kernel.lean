import EL2.Term
import EL2.Util

namespace EL2


class Context Ctx where
  set: Ctx → String → Term × Term × Int → Ctx
  get?: Ctx → String → Option (Term × Term × Int)

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

partial def infer? [Repr Ctx] [Context Ctx] (reduce: Bool) (ctx: Ctx) (term: Term): Option (Ctx × Term × Term × Int) := do
  let isLam? (term: Term): Option (List (String × Term) × Term) :=
    match term with
      | .lam params body => some (params, body)
      | _ => none

  -- return (ctx, term, type, level)
  match term with
    | .univ level =>
      pure (ctx, term, Term.univ (level+1), level+1)

    | .var name =>
      let (term, type, level) ← Context.get? ctx name
      pure (ctx, term, type, level)

    | .inh type cons name =>
      let (_, typeTerm, _, typeLevel) ← infer? reduce ctx type
      pure (ctx, .inh typeTerm cons name, typeTerm, typeLevel - 1)

    | .typ value =>
      let (_, _, valueType, _) ← infer? reduce ctx value
      infer? reduce ctx valueType

    | .list init last =>
      let (listCtx, _) ← Util.optionCtxMap? init ctx (infer? reduce)
      infer? reduce listCtx last

    | .bind name value =>
      let (_, term, type, level) ← infer? reduce ctx value
      pure (Context.set ctx name (term, type, level), term, type, level)

    | .lam params body =>
      let counterCtx := {field := ctx, count := 0 : Counter Ctx}

      let (counterCtx, paramsType) ← -- : Ctx × List (String × Term × Term × Int)
        Util.optionCtxMap? params counterCtx (λ counterCtx (name, type) => do
          let (ctx, type, typeType, typeLevel) ← infer? reduce counterCtx.field type
          let dummyCons := counterCtx.dummyName
          let counterCtx := counterCtx.next
          -- set dummy arg
          let (term, level) := (Term.inh type dummyCons [], typeLevel - 1)
          let ctx := Context.set ctx dummyCons (term, type, level)

          pure (counterCtx.with ctx, (name, type, typeType, typeLevel))
        )

      let dummyArgCtx := counterCtx.field
      let (_, bodyTerm, bodyType, bodyLevel) ← infer? reduce dummyArgCtx body

      -- type of Lam is Pi
      let params := paramsType.map (λ (name, type, _, _) => (name, type))
      let term := Term.lam params body
      let type := Term.lam params (.typ body)

      let paramsLevel := paramsType.map (λ (_, _, _, typeLevel) => typeLevel - 1)
      let level := paramsLevel.foldl max bodyLevel

      pure (ctx, term, type, level)

    | .app cmd args =>
      let (_, cmdTerm, cmdType, cmdLevel) ← infer? reduce ctx cmd
      let (_, args) ← Util.optionCtxMap? args ctx (λ ctx arg => do
        let (ctx, argTerm, argType, argLevel) ← infer? reduce ctx arg
        pure (ctx, (argTerm, argType, argLevel))
      )

      let (params, body) ← isLam? cmdTerm






      none
    | _ => none



end EL2
