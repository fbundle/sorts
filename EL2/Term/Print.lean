import EL2.Term.Term

namespace EL2.Term

structure PrintCtx where
  indentNum: Nat
  indentSize: Nat

def PrintCtx.next (ctx: PrintCtx): PrintCtx := {
  indentNum := ctx.indentNum+1,
  indentSize := ctx.indentSize,
}

def PrintCtx.indentStr (ctx: PrintCtx): String :=
  String.mk (List.replicate (ctx.indentNum * ctx.indentSize) ' ')


def printList (l: List String) (withParens: Bool := true): String :=
  match l with
    | [] => ""
    | x :: [] => x
    | _ =>
      let content := String.join (l.intersperse " ")
      if withParens then
        "(" ++ content ++ ")"
      else
        content

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
          ctx.indentStr ++ "let " ++ (printList [name, ":=", ctx.print value] (withParens := false)) ++ "\n"
        )
        let lastStr := ctx.indentStr ++ (ctx.print last) ++ "\n"

        String.join (initStrList ++ [lastStr])

    | lam {params, body} =>
      printList (
        ["fun"] ++ params.map (
          λ {name, type} => s!"({name}: {ctx.print type})"
        ) ++ ["=>", ctx.print body]
      )

    | app {cmd, args} =>
      printList ( [ctx.print cmd] ++ args.map ctx.print )

    | mat {cond, cases} =>
      let caseCtx := ctx.next
      let casesStrList := cases.map (λ {patCons, patArgs, value} =>
        caseCtx.indentStr ++ patCons ++ " " ++ (String.join (patArgs.intersperse " "))
        ++  " => " ++ (caseCtx.print value) ++ "\n"
      )

      "match " ++ (ctx.print cond) ++ " with" ++ "\n"
      ++ (String.join casesStrList)
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
