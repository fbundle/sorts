import EL.EL
import EL2.EL2
import Std

namespace EL2_EXAMPLE
open EL2

def el2 (_: Unit): List (Code Atom) :=
  [
    .bind_typ {name := "typ_Bool", params := []},
    .bind_val {name := "Bool", value := (.app {cmd := (.name "typ_Bool"), args := []})},   -- Bool := typ_Bool ()
    .bind_mk {name := "mk_true", params := [], type := {cmd := "Bool", args := []}},      -- mk_true := () => Bool
    .bind_mk {name := "mk_false", params := [], type := {cmd := "Bool", args := []}},     -- mk_false := () => Bool
    .bind_val {name := "true", value := (.app {cmd := (.name "mk_true"), args := []})},   -- true := mk_true ()
    .bind_val {name := "false", value := (.app {cmd := (.name "mk_false"), args := []})}, -- false := mk_false ()

    .bind_typ {name := "typ_Nat", params := []},
    .bind_val {name := "Nat", value := (.app {cmd := (.name "typ_Nat"), args := []})},
    .bind_mk {name := "mk_zero", params := [], type := {cmd := "Nat", args := []}},
    .bind_mk {
      name := "mk_cons",
      params := [{name := "n", type := (.name "Nat")}],
      type := {cmd := "Nat", args := []},
    },
    .bind_val {name := "zero", value := (.app {cmd := (.name "mk_zero"), args := []})},
    .bind_val {name := "cons", value := (.abst {
      params := [{name := "x", type := (.name "Nat")}],
      body := (.app {cmd := (.name "mk_zero"), args := [(.name "x")]}),
    })},



  ]


end EL2_EXAMPLE





def main (args : List String) : IO UInt32 := do
  match args with
  | [fileName] => do
      let content ← IO.FS.readFile fileName
      let tokens := EL.tokenize content
      let result := Util.parseAll EL.parse tokens
      if result.remaining.length ≠ 0 then
        let remaining := String.join (result.remaining.intersperse " ")
        IO.println s!"{repr result.items}\nerror at {remaining}"
        return 1
      else
        IO.println s!"{repr result.items}"
        return 0
  | _ => do
      IO.eprintln "Usage: el <file>"
      return 1
