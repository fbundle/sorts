import EL2.Term

namespace EL2

structure PrintCtx where
  indentNum: Nat
  indentSize: Nat
  stripParens: Bool

def PrintCtx.next (ctx: PrintCtx): PrintCtx := {
  ctx with
  stripParens := false
}

def PrintCtx.nextIndent (ctx: PrintCtx): PrintCtx := {
  ctx with
  indentNum := ctx.indentNum+1
  stripParens := true
}

def PrintCtx.indentStr (ctx: PrintCtx): String :=
  String.mk (List.replicate (ctx.indentNum * ctx.indentSize) ' ')


def printList (l: List String) (stripParens: Bool := False): String :=
  match l with
    | [] => ""
    | x :: [] => x
    | _ =>
      let content := String.join (l.intersperse " ")
      if stripParens then
        content
      else
        "(" ++ content ++ ")"


partial def PrintCtx.print [ToString β] (ctx: PrintCtx) (c: Term β): String :=
  let contentList: List String := match c with
    | .atom x =>
      [toString x]

    | .var n =>
      [n]

    | .list l =>
      ["\n" ++ String.join (l.map (λ x => ctx.indentStr ++ (ctx.nextIndent.print x) ++ "\n"))]

    | .ann x => [x.name, ":", ctx.next.print x.type]

    | .bind_val x => [x.name, ":=", ctx.next.print x.value]

    | .bind_typ x =>
      ["type"] ++
      [x.name] ++
      x.params.map (ctx.next.print ∘ (Term.ann ·))


    | .bind_mk x =>
      ["type_mk"] ++
      [x.name] ++
      x.params.map (ctx.next.print ∘ (Term.ann ·)) ++
      [printList (
        [x.type.cmd] ++
        x.type.args.map ctx.next.print
      ) false]

    | .app x =>
      [ctx.next.print x.cmd] ++
      x.args.map ctx.next.print



    | .lam x =>
      x.params.map (ctx.next.print ∘ (Term.ann ·)) ++
      [ctx.next.print x.body]

  printList contentList ctx.stripParens

instance [ToString β]: ToString (Term β) where
  toString (c: Term β):= {
    indentNum := 0,
    indentSize := 2
    stripParens := true
  :PrintCtx}.print c

end EL2
