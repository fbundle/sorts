import EL2.Core.CoreV2
import EL2.Term.Term

open EL2.Term

def main  : IO Unit := do
  IO.println "--------------------------------------"
  IO.println s!"{test0.toExp.toString}"
  IO.println s!"{typeCheck? test0 (Term.typ 0)}"
