namespace EL2.Term.Util

def statefulMapM [Monad m] (xs: List α) (state: State) (f: State → α → m (State × β)) : m (State × List β) :=
  let rec loop (ys: Array β) (xs: List α) (state: State): m (State × List β) := do
    match xs with
      | [] => pure (state, ys.toList)
      | x :: xs =>
        let (state, y) ← f state x
        loop (ys.push y) xs state

  loop #[] xs state

def statefulMap (xs: List α) (state: State) (f: State → α → State × β): State × List β :=
  Id.run (statefulMapM xs state (λ s x => pure (f s x)))

def liftExcept (o: Option β) (e: α): Except α β :=
  match o with
    | some v => Except.ok v
    | none => Except.error e

end EL2.Term.Util
