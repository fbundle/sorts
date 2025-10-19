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
partial def PrintCtx.printAnn [ToString β] (ctx: PrintCtx) (x: Ann (Term β)): String :=
  let contentList: List String := [x.name, ":", ctx.print x.type]
  printList contentList

partial def PrintCtx.print [ToString β] (ctx: PrintCtx) (c: Term β): String :=
  let contentList: List String := match c with
    | .atom a =>
      [toString a]

    | .var name =>
      [name]

    | .list {init, tail} =>
      if init.length = 0 then
        [ctx.print tail]
      else
        let parts := (init ++ [tail]).map (λ x => ctx.indentStr ++ (ctx.withIndent.print x) ++ "\n")
        ["\n" ++ String.join parts]

    | .bind_val {name, value} =>
      ["bind_val", name, ctx.print value]

    | .bind_typ {name, params, parent} =>
      ["bind_typ"] ++
      [name] ++
      params.map ctx.printAnn ++
      [ctx.print parent]

    | .bind_mk {name, params, type} =>
      let {cmd, args} := type
      ["bind_mk"] ++
      [name] ++
      params.map ctx.printAnn ++
      ["->"] ++
      [printList (
        [cmd] ++
        args.map ctx.print
      )]

    | .app {cmd, args} =>
      [ctx.print cmd] ++
      args.map ctx.print

    | .lam {params, body} =>
      params.map ctx.printAnn ++
      ["=>"] ++
      [ctx.print body]

    | .mat {cond, cases} =>
      ["match"] ++
      [ctx.print cond] ++
      cases.map (λ {pattern, value} =>

        "\n" ++ ctx.indentStr ++ printList (
          [pattern.cmd] ++
          pattern.args ++
          ["=>"] ++
          [ctx.print value]
        )
      ) ++ ["\n"]

  printList contentList

end

instance [ToString β]: ToString (Term β) where
  toString (c: Term β):= {
    indentNum := 0,
    indentSize := 2
  :PrintCtx}.print c

end EL2
