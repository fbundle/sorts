import EL2.Core.Term
import EL2.Core.TermNot


namespace EL2.Core

structure PrintCtx where
  indentNum: Nat
  indentSize: Nat

def PrintCtx.nextIndent (ctx: PrintCtx): PrintCtx := {
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
    | inh {type, cons, args} =>
      ["inh", ctx.print type, cons] ++
      args.map ctx.print

    | typ {value} =>
      ["typ", ctx.print value]

    | univ level =>
      [s!"U_{level}"]

    | var name =>
      [name]

    | bnd {init, last} =>
      if init.length = 0 then
        [ctx.print last]
      else
        let parts := init.map (λ {name, value} =>
          ctx.indentStr ++ (printList [name, ":=", ctx.nextIndent.print value]) ++ "\n"
        )
        ["\n" ++ String.join parts] ++ [ctx.indentStr ++ (ctx.print last) ++ "\n"]

    | lam {params, body} =>
      params.map (λ {name, type} =>
        s!"({name}: {ctx.print type})"
      ) ++
      ["=>", ctx.print body]

    | app {cmd, args} =>
      [ctx.print cmd] ++
      args.map ctx.print

    | mat {cond, cases} =>
      ["match", ctx.print cond, "with"] ++
      cases.map (λ {patCmd, patArgs, value} =>

        "\n" ++ ctx.indentStr ++ printList (
          [patCmd] ++
          patArgs ++
          ["=>", ctx.print value]
        )
      ) ++ ["\n"]

  printList contentList

end

instance : ToString Term where
  toString (c: Term) := {
    indentNum := 1,
    indentSize := 2,
    :PrintCtx
  }.print c

instance: Repr Term where
  reprPrec (term: Term) (prec: Nat): Std.Format := toString term


end EL2.Core
