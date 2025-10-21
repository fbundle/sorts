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

partial def PrintCtx.printAnn (ctx: PrintCtx) (x: Ann Term): String :=
  let contentList: List String := [x.name, ":", ctx.print x.type]
  printList contentList

partial def PrintCtx.print (ctx: PrintCtx) (term: Term): String :=
  let contentList: List String := match term with
    | .inh type =>
      ["inh", ctx.print type]

    | .univ level =>
      [s!"U_{level}"]

    | .var name =>
      [name]

    | .lst {init, last} =>
      if init.length = 0 then
        [ctx.print last]
      else
        let parts := (init ++ [last]).map (λ x => ctx.indentStr ++ (ctx.withIndent.print x) ++ "\n")
        ["\n" ++ String.join parts]

    | .bind_val {name, value} =>
      ["bind_val", name, ctx.print value]

    | .bind_typ {name, params, level} =>
      ["bind_typ", name] ++
      params.map ctx.printAnn ++
      [s!"U_{level}"]

    | .bind_mk {name, params, type} =>
      let {cmd, args} := type
      ["bind_mk", name] ++
      params.map ctx.printAnn ++
      ["->"] ++
      [printList (
        [cmd] ++
        args.map ctx.print
      )]

    | .typ {value} =>
      ["type", ctx.print value]

    | .app {cmd, args} =>
      [ctx.print cmd] ++
      args.map ctx.print

    | .lam {params, body} =>
      params.map ctx.printAnn ++
      ["=>", ctx.print body]

    | .mat {cond, cases} =>
      ["match", ctx.print cond, "with"] ++
      cases.map (λ {pattern, value} =>

        "\n" ++ ctx.indentStr ++ printList (
          [pattern.cmd] ++
          pattern.args ++
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
