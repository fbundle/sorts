import Parser.Combinator
import EL2.Core

namespace EL2.Parser1
open Parser
open EL2.Core

abbrev Parser := Combinator.Parser Char

mutual

def parseExp: Parser Exp :=
  sorry

def parseUniv: Parser Exp := λ xs => do
  let (name, rest) ← Combinator.parseName xs
  if "Type".isPrefixOf name then
    let levelStr := name.stripPrefix "Type"
    let level ← levelStr.toNat?
    some (Exp.typ level, rest)
  else
    none




end
#eval parseUniv "Type123".toList


end EL2.Parser1
