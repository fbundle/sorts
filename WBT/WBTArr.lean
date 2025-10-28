import WBT.WBT

namespace WBT

structure WBTArr (α: Type u) where
  node? : Option (Node.Node α)

def WBTArr.fromNode (node?: Option (Node.Node α)): WBTArr α :=
  {node? := node?}

def WBTArr.empty : WBTArr α :=
  {node? := none}

partial def WBTArr.toArray (a: WBTArr α): Array α :=
  Node.iterate a.node?

def WBTArr.toList (a: WBTArr α): List α :=
  a.toArray.toList

def WBTArr.length (a: WBTArr α): Nat :=
  Node.weight a.node?
def WBTArr.depth (a: WBTArr α): Nat :=
  Node.height a.node?

instance [Repr α]: Repr (WBTArr α) where
  reprPrec (a: WBTArr α) (_: Nat): Std.Format :=
    s!"WBTMap {repr a.toArray}"

partial def WBTArr.get? (a: WBTArr α) (i: Nat): Option α :=
  match a.node? with
    | none => none
    | some n =>
      let (leftWeight, rightWeight) := (Node.weight n.left?, Node.weight n.right?)
      if i < leftWeight then
        WBTArr.get? (WBTArr.fromNode n.left?) i
      else if i = leftWeight then
        some n.entry
      else if i < 1 + leftWeight + rightWeight then
        WBTArr.get? (WBTArr.fromNode n.right?) (i - 1 - leftWeight)
      else
        none

partial def WBTArr.set? (a: WBTArr α) (i: Nat) (x: α): Option (WBTArr α) := do
  match a.node? with
    | none => none
    | some n =>
      let (leftWeight, rightWeight) := (Node.weight n.left?, Node.weight n.right?)
      if i < leftWeight then
        let l1 ← WBTArr.set? (WBTArr.fromNode n.left?) i x
        let n1 := Node.makeNode n.entry l1.node? n.right?
        pure (WBTArr.fromNode (Node.balance Node.δ n1))
      else if i = leftWeight then
        let n1 := Node.makeNode x n.left? n.right?
        pure (WBTArr.fromNode (Node.balance Node.δ n1))
      else if i < 1 + leftWeight + rightWeight then
        let r1 ← WBTArr.set? (WBTArr.fromNode n.right?) (i - 1 - leftWeight) x
        let n1 := Node.makeNode n.entry n.left? r1.node?
        pure (WBTArr.fromNode (Node.balance Node.δ n1))
      else
        none

partial def WBTArr.insert? (a: WBTArr α) (i: Nat) (x: α): Option (WBTArr α) := do
  match a.node? with
    | none =>
      if i ≠ 0 then none else some (WBTArr.fromNode (Node.makeNode x none none))
    | some n =>
      let (leftWeight, rightWeight) := (Node.weight n.left?, Node.weight n.right?)
      if i ≤ leftWeight then
        let l1 ← WBTArr.insert? (WBTArr.fromNode n.left?) i x
        let n1 := Node.makeNode n.entry l1.node? n.right?
        pure (WBTArr.fromNode (Node.balance Node.δ n1))
      else if i ≤ 1 + leftWeight + rightWeight then
        let r1 ← WBTArr.insert? (WBTArr.fromNode n.right?) (i - 1 - leftWeight) x
        let n1 := Node.makeNode n.entry n.left? r1.node?
        pure (WBTArr.fromNode (Node.balance Node.δ n1))
      else
        none

partial def WBTArr.delete? (a: WBTArr α) (i: Nat) : Option (WBTArr α) := do
  match a.node? with
    | none => none
    | some n =>
      let (leftWeight, rightWeight) := (Node.weight n.left?, Node.weight n.right?)
      if i < leftWeight then
        let l1 ← WBTArr.delete? (WBTArr.fromNode n.left?) i
        let n1 := Node.makeNode n.entry l1.node? n.right?
        pure (WBTArr.fromNode (Node.balance Node.δ n1))
      else if i = leftWeight then
        match n.right? with
          | none => pure (WBTArr.fromNode n.left?)
          | some r =>
            let x ← WBTArr.get? (WBTArr.fromNode r) 0
            let r1 ← WBTArr.delete? (WBTArr.fromNode r) 0
            let n1 := Node.makeNode x n.left? r1.node?
            pure (WBTArr.fromNode (Node.balance Node.δ n1))
      else if i < 1 + leftWeight + rightWeight then
        let r1 ← WBTArr.delete? (WBTArr.fromNode n.right?) (i - 1 - leftWeight)
        let n1 := Node.makeNode n.entry n.left? r1.node?
        pure (WBTArr.fromNode (Node.balance Node.δ n1))
      else
        none


partial def WBTArr.mapM [Monad m] (a: WBTArr α) (f: α → m β): m (WBTArr β) := do
  let rec loop (n?: Option (Node.Node α)): m (Option (Node.Node β)) := do
    match n? with
      | none => pure none
      | some n =>
        let entry ← f n.entry
        let left? ← loop n.left?
        let right? ← loop n.right?
        pure (Node.makeNode entry left? right?)

  let node? ← loop a.node?
  pure {node? := node? : WBTArr β}

def WBTArr.push (a: WBTArr α) (x: α): WBTArr α :=
  (a.insert? a.length x).get sorry

partial def WBTArr.fromList (xs: List α): WBTArr α :=
  let rec loop (a: WBTArr α) (xs: List α): WBTArr α :=
    match xs with
    | [] => a
    | x :: xs =>
      loop (a.push x) xs

  loop WBTArr.empty xs

def WBTArr.fromArray (xs: Array α): WBTArr α :=
  WBTArr.fromList xs.toList


#eval (WBTArr.fromArray (Array.replicate 1000 1)).node?

end WBT
