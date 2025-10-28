import WBT.WBT

namespace WBT

structure WBTArr (α: Type u) where
  node? : Option (Node α)

instance: Coe (Option (Node α)) (WBTArr α) where
  coe (node?: Option (Node α)): WBTArr α := {node? := node?}

def WBTArr.empty : WBTArr α :=
  {node? := none}

partial def WBTArr.toArray (a: WBTArr α): Array α :=
  iterate a.node?

def WBTArr.toList (a: WBTArr α): List α :=
  a.toArray.toList

def WBTArr.length (a: WBTArr α): Nat :=
  weight a.node?
def WBTArr.depth (a: WBTArr α): Nat :=
  height a.node?

instance [Repr α]: Repr (WBTArr α) where
  reprPrec (a: WBTArr α) (_: Nat): Std.Format :=
    s!"WBTArr(l={a.length}, d={a.depth})"

partial def WBTArr.get? (a: WBTArr α) (i: Nat): Option α :=
  match a.node? with
    | none => none
    | some n =>
      let (leftWeight, rightWeight) := (weight n.left?, weight n.right?)
      if i < leftWeight then
        WBTArr.get? n.left? i
      else if i = leftWeight then
        some n.entry
      else if i < 1 + leftWeight + rightWeight then
        WBTArr.get? n.right? (i - 1 - leftWeight)
      else
        none

partial def WBTArr.set? (a: WBTArr α) (i: Nat) (x: α): Option (WBTArr α) := do
  let rec loop (n?: Option (Node α)) (i: Nat): Option (Option (Node α)) := do
    match n? with
      | none => none
      | some n =>
        let (leftWeight, rightWeight) := (weight n.left?, weight n.right?)
        if i < leftWeight then
          let l1 ← loop n.left? i
          let n1 := makeNode n.entry l1 n.right?
          some (balance δ n1)
        else if i = leftWeight then
          let n1 := makeNode x n.left? n.right?
          some n1
        else if i < 1 + leftWeight + rightWeight then
          let r1 ← loop n.right? (i - 1 - leftWeight)
          let n1 := makeNode n.entry n.left? r1
          some (balance δ n1)
        else
          none

  let node? ← loop a.node? i
  some {node? := node?}

partial def WBTArr.insert? (a: WBTArr α) (i: Nat) (x: α): WBTArr α :=
  let rec loop (n?: Option (Node α)) (i: Nat): Option (Node α) :=
    match n? with
      | none =>
        if i = 0 then
          some (makeNode x none none)
        else
          none
      | some n =>
        let (leftWeight, rightWeight) := (weight n.left?, weight n.right?)
        if i ≤ leftWeight then
          let l1 := loop n.left? i
          match loop n.left? i with
            | none => none
            | some l1 =>
              let n1 := makeNode n.entry l1 n.right?
              some (balance δ n1)
        else if i ≤ 1 + leftWeight + rightWeight then
          match loop n.right? (i - 1 - leftWeight) with
            | none => none
            | some r1 =>
              let n1 := makeNode n.entry n.left? r1
              some (balance δ n1)
        else
          none

  {node? := loop a.node? i}

partial def WBTArr.delete? (a: WBTArr α) (i: Nat) : Option (WBTArr α) :=
  let rec loop (n?: Option (Node α)) (i: Nat): Option (Node α) :=
    match n? with
      | none => none
      | some n =>
        let (leftWeight, rightWeight) := (weight n.left?, weight n.right?)
        if i < leftWeight then
          match loop n.left? i with
            | none => none
            | some l1 =>
              let n1 := makeNode n.entry l1 n.right?
              some (balance δ n1)
        else if i = leftWeight then
          if leftWeight = 0 then
            n.right?
          else
            let newEntry? := {node? := n.left? : WBTArr α}.get? (leftWeight - 1)
            match newEntry? with
              | none => none
              | some newEntry =>
                let l1 := loop n.left? (leftWeight - 1)
                let n1 := makeNode newEntry l1 n.right?
                some (balance δ n1)
        else if i < 1 + leftWeight + rightWeight then
          match loop n.right? (i - 1 - leftWeight) with
            | none => none
            | some r1 =>
              let n1 := makeNode n.entry n.left? r1
              some (balance δ n1)
        else
          none
  some {node? := loop a.node? i}

partial def WBTArr.mapM [Monad m] (a: WBTArr α) (f: α → m β): m (WBTArr β) := do
  let rec loop (n?: Option (Node α)): m (Option (Node β)) := do
    match n? with
      | none => pure none
      | some n =>
        let entry ← f n.entry
        let left? ← loop n.left?
        let right? ← loop n.right?
        pure (makeNode entry left? right?)

  let node? ← loop a.node?
  pure {node? := node? : WBTArr β}

def WBTArr.push (a: WBTArr α) (x: α): WBTArr α :=
  match a.insert? a.length x with
    | none => sorry
    | some a => a

partial def WBTArr.fromList (xs: List α): WBTArr α :=
  let rec loop (a: WBTArr α) (xs: List α): WBTArr α :=
    match xs with
    | [] => a
    | x :: xs =>
      loop (a.push x) xs

  loop WBTArr.empty xs

def WBTArr.fromArray (xs: Array α): WBTArr α :=
  WBTArr.fromList xs.toList

#eval balanceCond 3 (WBTArr.fromArray (Array.replicate 1000 1)).node?

end WBT
