import EL2.Term

namespace EL2

structure PrintCtx where
  indentNum: Nat
  indentSize: Nat

def PrintCtx.withIndent (ctx: PrintCtx): PrintCtx := {
  ctx with
  indentNum := ctx.indentNum+1
}

def PrintCtx.indentStr (ctx: PrintCtx): String :=
  String.mk (List.replicate (ctx.indentNum * ctx.indentSize) ' ')


def printList (l: List String): String :=
  match l with
    | [] => ""
    | x :: [] => x
    | _ =>
      let content := String.join (l.intersperse " ")
      "(" ++ content ++ ")"

mutual

partial def PrintCtx.print (ctx: PrintCtx) (term: Term): String :=
  let contentList: List String := match term with
    | .inh type method values =>
      ["inh", ctx.print type, method] ++
      values.map ctx.print

    | .infer value =>
      ["infer", ctx.print value]

    | .univ level =>
      [s!"U_{level}"]

    | .var name =>
      [name]

    | .list init last =>
      if init.length = 0 then
        [ctx.print last]
      else
        let parts := (init ++ [last]).map (λ x => ctx.indentStr ++ (ctx.withIndent.print x) ++ "\n")
        ["\n" ++ String.join parts]

    | .bind name value =>
      ["bind", name, ctx.print value]

    | .app cmd args =>
      [ctx.print cmd] ++
      args.map ctx.print

    | .lam params body =>
      params.flatMap (λ (name, type) =>
        [name, ":", ctx.print type]
      ) ++
      ["=>", ctx.print body]

    | .mat cond cases =>
      ["match", ctx.print cond, "with"] ++
      cases.map (λ (patCmd, patArgs, value) =>

        "\n" ++ ctx.indentStr ++ printList (
          [patCmd] ++
          patArgs ++
          ["=>", ctx.print value]
        )
      ) ++ ["\n"]

  printList contentList

end

instance : ToString Term where
  toString (c: Term):= {
    indentNum := 0,
    indentSize := 2,
    :PrintCtx
  }.print c

end EL2
