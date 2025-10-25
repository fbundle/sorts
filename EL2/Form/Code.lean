import EL2.Form.Form
import EL2.Term.Term
import EL2.Form.Util

namespace EL2.Form
open EL2.Term

mutual

def parseInh? := {
  parseHead := ["inh"],
  parseList (list: List Form): Option (Inh Term) := do
    let typeForm ← list[0]?
    let type ← parse? typeForm
    let consForm ← list[1]?
    let cons ← isName? consForm
    let argsForm := list.extract 2
    let args ← Util.optionMap? argsForm parse?
    pure {
      type := type,
      cons := cons,
      args := args,
    }
  : ParseList (Inh Term)
}.parseForm

def parseTyp? := {
  parseHead := ["typ"],
  parseList (list: List Form): Option (Typ Term) := do
    let valueForm ← list[0]?
    let value ← parse? valueForm
    pure {
      value := value,
    }
  : ParseList (Typ Term)
}.parseForm

def parseBind? := {
  parseHead := ["lam", ":="],
  parseList (list: List Form): Option (Bind Term) := do
    let nameForm ← list[0]?
    let name ← isName? nameForm
    let valueForm ← list[1]?
    let value ← parse? valueForm
    pure {
      name := name,
      value := value,
    }
  : ParseList (Bind Term)
}.parseForm

def parseBnd? := {
  parseHead := ["let"],
  parseList (list: List Form): Option (Bnd Term) := do
    let initForm := list.extract 0 (list.length - 1)
    let init ← Util.optionMap? initForm parseBind?
    let lastForm ← list.getLast?
    let last ← parse? lastForm
    pure {
      init := init,
      last := last,
    }
  : ParseList (Bnd Term)
}.parseForm

def parseAnn? := {
  parseHead := ["ann", ":"],
  parseList (list: List Form): Option (Ann Term) := do
    let nameForm ← list[0]?
    let name ← isName? nameForm
    let typeForm ← list[1]?
    let type ← parse? typeForm
    pure {
      name := name,
      type := type,
    }
  : ParseList (Ann Term)
}.parseForm

def parseLam? := {
  parseHead := ["lam", "=>"],
  parseList (list :List Form): Option (Lam Term) := do
    let paramsForm := list.extract 0 (list.length - 1)
    let params ← Util.optionMap? paramsForm parseAnn?
    let bodyForm ← list.getLast?
    let body ← parse? bodyForm
    pure {
      params := params,
      body := body,
    }
  : ParseList (Lam Term)
}.parseForm

def parseApp? (list: List Form): Option (App Term) := do
  let cmdForm ← list[0]?
  let cmd ← parse? cmdForm
  let argsForm := list.extract 1
  let args ← Util.optionMap? argsForm parse?
  pure {
    cmd := cmd,
    args := args,
  }

def parseCase? := {
  parseHead := ["case", "=>"],
  parseList (list: List Form): Option (Case Term) := do
    let patCmdForm ← list[0]?
    let patCmd ← isName? patCmdForm
    let patArgsForm := list.extract 1 (list.length - 1)
    let patArgs ← Util.optionMap? patArgsForm isName?
    let valueForm ← list.getLast?
    let value ← parse? valueForm
    pure {
      patCmd := patCmd,
      patArgs := patArgs,
      value := value,
    }
  : ParseList (Case Term)
}.parseForm

def parseMat? := {
  parseHead := ["match"],
  parseList (list: List Form): Option (Mat Term) := do
    let condForm ← list[0]?
    let cond ← parse? condForm
    let casesForm := list.extract 1
    let cases ← Util.optionMap? casesForm parseCase?
    pure {
      cond := cond,
      cases := cases,
    }
  : ParseList (Mat Term)
}

def parse? (form: Form): Option Term :=



  none

end



end EL2.Form
