import EL2.Term

namespace EL2

structure PrintCtx where
  indentNum: Nat
  indentSize: Nat
  stripParens: Bool

def PrintCtx.next (ctx: PrintCtx): PrintCtx := {ctx with indentNum := ctx.indentNum+1}

def PrintCtx.indentStr (ctx: PrintCtx): String :=
  String.mk (List.replicate (ctx.indentNum * ctx.indentSize) ' ')


def printList (l: List String): String :=
  match l with
    | [] => ""
    | x :: [] => x
    | _ =>
      "(" ++ String.join (l.intersperse " ") ++ ")"

partial def PrintCtx.print [ToString β] (ctx: PrintCtx) (c: Term β): String :=
  match c with
    | .atom x =>
      toString x

    | .var n =>
      n

    | .list l =>
      "\n" ++ String.join (l.map (λ x => ctx.indentStr ++ (ctx.next.print x) ++ "\n"))

    | .ann x => s!"({x.name}: {ctx.print x.type})"

    | .bind_val x => s!"({x.name} := {ctx.print x.value})"

    | .bind_typ x => printList (
      ["type"] ++
      [x.name] ++
      x.params.map (ctx.print ∘ (Term.ann ·))
    )

    | .bind_mk x => printList (
      ["type_mk"] ++
      [x.name] ++
      x.params.map (ctx.print ∘ (Term.ann ·)) ++
      [printList (
        [x.type.cmd] ++
        x.type.args.map ctx.print
      )]
    )
    | .app x => printList (
      [ctx.print x.cmd] ++
      x.args.map ctx.print
    )


    | .lam x => printList (
      x.params.map (ctx.print ∘ (Term.ann ·)) ++
      [ctx.print x.body]
    )

instance [ToString β]: ToString (Term β) where
  toString (c: Term β):= {
    indentNum := 0,
    indentSize := 2
    stripParens := true
  :PrintCtx}.print c

end EL2
