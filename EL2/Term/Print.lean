import EL2.Term.Term

namespace EL2.Term

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
  match term with
    | inh {type, cons, args} =>
      printList (["inh", ctx.print type, cons] ++ args.map ctx.print)

    | univ level => s!"U_{level}"

    | var name => name

    | bnd {init, last} =>
      if init.length = 0 then
        ctx.print last
      else
        let initStrList := init.map (λ {name, value} =>
          ctx.nextIndent.nextIndent.indentStr ++ (printList [name, ":=", ctx.nextIndent.print value]) ++ "\n"
        )
        let lastStr := ctx.nextIndent.nextIndent.indentStr ++ (ctx.nextIndent.print last) ++ "\n"

        "\n" ++ ctx.nextIndent.indentStr ++ "(" ++ "let" ++ "\n"
        ++ (String.join (initStrList ++ [lastStr]))
        ++ ctx.nextIndent.indentStr ++ ")"

    | lam {params, body} =>
      printList (
        ["λ"] ++ params.map (
          λ {name, type} => s!"({name}: {ctx.print type})"
        ) ++ ["=>", ctx.print body]
      )

    | app {cmd, args} =>
      printList ( [ctx.print cmd] ++ args.map ctx.print )

    | mat {cond, cases} =>
      let casesStrList := cases.map (λ {patCmd, patArgs, value} =>
        ctx.nextIndent.nextIndent.indentStr ++ patCmd ++ " " ++ (String.join (patArgs.intersperse " "))
        ++  " => " ++ (ctx.nextIndent.print value) ++ "\n"
      )

      "\n" ++ ctx.nextIndent.indentStr ++ "(" ++ "match " ++ (ctx.print cond) ++ " with" ++ "\n"
      ++ (String.join casesStrList)
      ++ ctx.nextIndent.indentStr ++ ")"
end

instance : ToString Term where
  toString (c: Term) := {
    indentNum := 0,
    indentSize := 2,
    :PrintCtx
  }.print c

instance: Repr Term where
  reprPrec (term: Term) (_: Nat): Std.Format := toString term


end EL2.Term
