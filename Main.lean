import EL2.Core.CoreV2
open EL2.Core

def main  : IO Unit := do
  IO.println "--------------------------------------"
  IO.println s!"{test0.toExp.toString}"
  IO.println s!"{typeCheck? test0 (Term.typ 0)}"
