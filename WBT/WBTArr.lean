import WBT.WBT

namespace WBT

structure WBTArr (α: Type u) where
  node? : Option (Node α)

def WBTArr.empty : WBTArr α := {node? := none}

partial def WBTArr.toArray (a: WBTArr α): Array α :=
  iterate a.node?

def WBTArr.size (a: WBTArr α): Nat := weight a.node?

partial def WBTArr.get? (a: WBTArr α) (i: Nat): Option α :=
  let rec loop (n?: Option (Node α)) (i: Nat): Option α :=
    match n? with
      | none => none
      | some n =>
        let (leftWeight, rightWeight) := (weight n.left?, weight n.right?)
        if i < leftWeight then
          loop n.left? i
        else if i = leftWeight then
          n.entry
        else if i < 1 + leftWeight + rightWeight then
          loop n.right? (i - 1 - leftWeight)
        else
          none
  loop a.node? i

partial def WBTArr.set? (a: WBTArr α) (i: Nat) (x: α): Option (WBTArr α) :=
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
              balance δ n1
        else if i = leftWeight then
          let n1 := makeNode x n.left? n.right?
          some n1
        else if i < 1 + leftWeight + rightWeight then
          match loop n.right? (i - 1 - leftWeight) with
            | none => none
            | some r1 =>
              let n1 := makeNode n.entry n.left? r1
              balance δ n1
        else
          none

  some {node? := loop a.node? i}

partial def WBTArr.insert? (a: WBTArr α) (i: Nat) (x: α): Option (WBTArr α) :=
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
          match loop n.left? i with
            | none => none
            | some l1 =>
              let n1 := makeNode n.entry l1 n.right?
              balance δ n1
        else if i ≤ 1 + leftWeight + rightWeight then
          match loop n.right? (i - 1 - leftWeight) with
            | none => none
            | some r1 =>
              let n1 := makeNode n.entry n.left? r1
              balance δ n1
        else
          none

  some {node? := loop a.node? i}

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
              balance δ n1
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
                balance δ n1
        else if i < 1 + leftWeight + rightWeight then
          match loop n.right? (i - 1 - leftWeight) with
            | none => none
            | some r1 =>
              let n1 := makeNode n.entry n.left? r1
              balance δ n1
        else
          none
  some {node? := loop a.node? i}

partial def WBTArr.merge (a: WBTArr α) (b: WBTArr α): WBTArr α :=
  sorry

partial def WBTArr.split? (a: WBTArr α) (i: Nat): (WBTArr α) × (WBTArr α) :=
  sorry

partial def WBTArr.mapM [Monad m] (a: WBTArr α) (f: α → m β): m (WBTArr β) := do
  let rec loop (n?: Option (Node α)): m (Option (Node β)) := do
    match n? with
      | none => pure none
      | some n =>
        let entry ← f n.entry
        let left? ← loop n.left?
        let right? ← loop n.right?
        pure $ makeNode entry left? right?

  let node? ← loop a.node?
  pure {node? := node? : WBTArr β}

end WBT
