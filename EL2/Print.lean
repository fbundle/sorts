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
    | .atom x =>
      [toString x]

    | .var n =>
      [n]

    | .list init tail =>
      let parts := (init ++ [tail]).map (λ x => ctx.indentStr ++ (ctx.withIndent.print x) ++ "\n")
      ["\n" ++ String.join parts]

    | .bind_val x =>
      ["bind_val", x.name, ctx.print x.value]

    | .bind_typ x =>
      ["bind_typ"] ++
      [x.name] ++
      x.params.map ctx.printAnn ++
      [ctx.print x.parent]

    | .bind_mk x =>
      ["bind_mk"] ++
      [x.name] ++
      x.params.map ctx.printAnn ++
      ["->"] ++
      [printList (
        [x.type.cmd] ++
        x.type.args.map ctx.print
      )]

    | .app x =>
      [ctx.print x.cmd] ++
      x.args.map ctx.print

    | .lam x =>
      x.params.map ctx.printAnn ++
      ["=>"] ++
      [ctx.print x.body]

    | .mat x =>
      ["match"] ++
      [ctx.print x.cond] ++
      x.cases.map (λ case =>
        "\n" ++ ctx.indentStr ++ printList (
          [case.pattern.cmd] ++
          case.pattern.args ++
          ["=>"] ++
          [ctx.print case.value]
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
