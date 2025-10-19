import EL2.Term

namespace EL2

structure PrintCtx where
  indentNum: Nat
  indentSize: Nat
  stripParens: Bool := false

def PrintCtx.withParens (ctx: PrintCtx): PrintCtx := {
  ctx with
  stripParens := false
}

def PrintCtx.withIndent (ctx: PrintCtx): PrintCtx := {
  ctx with
  indentNum := ctx.indentNum+1
  stripParens := true
}

def PrintCtx.indentStr (ctx: PrintCtx): String :=
  String.mk (List.replicate (ctx.indentNum * ctx.indentSize) ' ')


def printList (l: List String) (stripParens: Bool): String :=
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
      ["\n" ++ String.join (l.map (λ x => ctx.indentStr ++ (ctx.withIndent.print x) ++ "\n"))]

    | .ann x =>
      [x.name, ":", ctx.withParens.print x.type]

    | .bind_val x =>
      ["bind_val", x.name, ctx.withParens.print x.value]

    | .bind_typ x =>
      ["bind_typ"] ++
      [x.name] ++
      x.params.map (ctx.withParens.print ∘ (Term.ann ·)) ++
      [ctx.withParens.print x.parent]

    | .bind_mk x =>
      ["bind_mk"] ++
      [x.name] ++
      x.params.map (ctx.withParens.print ∘ (Term.ann ·)) ++
      ["->"] ++
      [printList (
        [x.type.cmd] ++
        x.type.args.map ctx.withParens.print
      ) false]

    | .app x =>
      [ctx.withParens.print x.cmd] ++
      x.args.map ctx.withParens.print

    | .lam x =>
      x.params.map (ctx.withParens.print ∘ (Term.ann ·)) ++
      ["=>"] ++
      [ctx.withParens.print x.body]

    | .mat x =>
      ["match"] ++
      [ctx.withParens.print x.cond] ++
      x.cases.map (λ case =>
        "\n" ++ ctx.indentStr ++ printList (
          [case.pattern.cmd] ++
          case.pattern.args ++
          ["=>"] ++
          [ctx.withParens.print case.value]
        ) true
      ) ++ ["\n"]

  printList contentList ctx.stripParens

instance [ToString β]: ToString (Term β) where
  toString (c: Term β):= {
    indentNum := 0,
    indentSize := 2
    stripParens := true
  :PrintCtx}.print c

end EL2
