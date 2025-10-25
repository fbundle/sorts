import EL2.Term.Term

namespace EL2.Term


def T.optionMap? (t: T α) (f: α → Option β) : Option (T β) := do
  match t with
    | T.inh x =>
      let type ← f x.type
      let args ← Util.optionMap? x.args f
      T.inh {
        type := type,
        cons := x.cons,
        args := args
      }
    | T.typ x =>
      let value ← f x.value
      T.typ {
        value := value,
      }
    | T.bnd x =>
      let init ← Util.optionMap? x.init (λ bind => do
        let value ← f bind.value
        pure {
          name := bind.name,
          value := value,
          : Bind β
        }
      )
      let last ← f x.last
      T.bnd {
        init := init,
        last := last,
      }
    | T.lam x =>
      let params ← Util.optionMap? x.params (λ param => do
        let type ← f param.type
        pure {
          name := param.name,
          type := type,
          : Ann β
        }
      )
      let body ← f x.body
      T.lam {
        params := params,
        body := body,
      }
    | T.app x =>
      let cmd ← f x.cmd
      let args ← Util.optionMap? x.args f
      T.app {
        cmd := cmd,
        args := args,
      }
    | T.mat x =>
      let cond ← f x.cond
      let cases ← Util.optionMap? x.cases (λ case => do
        let value ← f case.value
        pure {
          patCmd := case.patCmd,
          patArgs := case.patArgs,
          value := value
          : Case β
        }
      )
      T.mat {
        cond := cond,
        cases := cases,
      }


def T.map (t: T α) (f: α → β) : T β :=
  match t with
    | T.inh x =>
      let type := f x.type
      let args := x.args.map f
      T.inh {
        type := type,
        cons := x.cons,
        args := args
      }
    | T.typ x =>
      let value := f x.value
      T.typ {
        value := value,
      }
    | T.bnd x =>
      let init := x.init.map (λ bind =>
        let value := f bind.value
        {
          name := bind.name,
          value := value,
          : Bind β
        }
      )
      let last := f x.last
      T.bnd {
        init := init,
        last := last,
      }
    | T.lam x =>
      let params := x.params.map (λ param =>
        let type := f param.type
        {
          name := param.name,
          type := type,
          : Ann β
        }
      )
      let body := f x.body
      T.lam {
        params := params,
        body := body,
      }
    | T.app x =>
      let cmd := f x.cmd
      let args := x.args.map f
      T.app {
        cmd := cmd,
        args := args,
      }
    | T.mat x =>
      let cond := f x.cond
      let cases := x.cases.map (λ case =>
        let value := f case.value
        {
          patCmd := case.patCmd,
          patArgs := case.patArgs,
          value := value
          : Case β
        }
      )
      T.mat {
        cond := cond,
        cases := cases,
      }


def Term.map (term: Term) (f: Term → Term): Term :=
  match term with
    | .univ _ => term
    | .var _ => term
    | .t x => Term.t (x.map f)

def Term.optionMap? (term: Term) (f: Term → Option Term): Option Term := do
  match term with
    | .univ _ => term
    | .var _ => term
    | .t x =>
      let y ← x.optionMap? f
      Term.t y


namespace EL2.Term
