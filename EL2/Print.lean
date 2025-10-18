import EL2.Term

namespace EL2

structure PrintCtx where
  indentNum: Nat
  indentSize: Nat
  stripParens: Bool

def PrintCtx.next (ctx: PrintCtx): PrintCtx := {ctx with indentNum := ctx.indentNum+1}

def PrintCtx.indentStr (ctx: PrintCtx): String :=
  String.mk (List.replicate (ctx.indentNum * ctx.indentSize) ' ')

partial def PrintCtx.print [ToString β] (ctx: PrintCtx) (c: Term β): String :=
  match c with
    | .atom x =>
      toString x

    | .var n =>
      n

    | .list l =>

      "\n"
      ++
      String.join (l.map (λ x => ctx.indentStr ++ (ctx.next.print x) ++ "\n"))

    | .ann x => s!"({x.name}: {ctx.print x.type})"

    | .bind_val x => s!"{x.name} := {ctx.print x.value}"

    | .bind_typ x => s!"type {x.name} {
      String.join ((x.params.map (ctx.print ∘ (Term.ann ·))).intersperse " ")
    }"

    | .bind_mk x => s!"type_mk {x.name} {
      String.join ((x.params.map (ctx.print ∘ (Term.ann ·))).intersperse " ")
    } -> " ++
    match x.type.args with
      | [] => x.type.cmd
      | _ => s!"({x.type.cmd} {
          String.join ((x.type.args.map ctx.print).intersperse " ")
        })"

    | .app x =>
      match x.args with
        | [] => ctx.print x.cmd
        | _ =>
          s!"({ctx.print x.cmd} {
            String.join ((x.args.map ctx.print).intersperse " ")
          })"

    | .lam x => s!"{
      String.join ((x.params.map (ctx.print ∘ (Term.ann ·))).intersperse " ")
    } => {ctx.print x.body}"

instance [ToString β]: ToString (Term β) where
  toString (c: Term β):= {
    indentNum := 0,
    indentSize := 2
    stripParens := true
  :PrintCtx}.print c

end EL2
