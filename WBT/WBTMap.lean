import WBT.WBT

namespace WBT

-- TODO potentially let cmp hashed α so that the tree will be more balanced
structure WBTMap (α: Type u) (β: Type v) (cmp: α → α → Ordering) where
  node? : Option (Node.Node (α × β))

-- as Lean enforces type to be strictly positive, sometimes recursive structure doesn't work
-- e.g
-- structure A where
--   val : Nat
--   map : Std.HashMap String A compare

-- e.g
-- structure A where
--   val : Nat
--   map : Lean.RBTree String A compare

-- somehow, List (String × A) and Array (String × A) work but it requires O(n) look up time

-- the whole purpose of this self-balancing tree is to do this
private structure A where
  val : Nat
  map : WBTMap String A compare

def WBTMap.min? (m: WBTMap α β cmp): Option (α × β) :=
  Node.left m.node?

def WBTMap.length (m: WBTMap α β cmp): Nat :=
  Node.weight m.node?

def WBTMap.depth (m: WBTMap α β cmp): Nat :=
  Node.height m.node?

instance: Coe (Option (Node.Node (α × β))) (WBTMap α β cmp) where
  coe (node?: Option (Node.Node (α × β))): WBTMap α β cmp := {node? := node?}

def WBTMap.empty : WBTMap α β cmp :=
  {node? := none}

partial def WBTMap.toArray (m: WBTMap α β cmp): Array (α × β) :=
  Node.iterate m.node?

def WBTArr.toList (m: WBTMap α β cmp): List (α × β) :=
  m.toArray.toList

instance [Repr α]: Repr (WBTMap α β cmp) where
  reprPrec (m: WBTMap α β cmp) (_: Nat): Std.Format :=
    s!"WBTArr(l={m.length}, d={m.depth})"

partial def WBTMap.get? (m: WBTMap α β cmp) (key: α): Option β :=
  match m.node? with
    | none => none
    | some n =>
      let (ekey, eval) := n.entry
      match cmp key ekey with
        | Ordering.lt =>
          WBTMap.get? (cmp := cmp) n.left? key
        | Ordering.eq =>
          eval
        | Ordering.gt =>
          WBTMap.get? (cmp := cmp) n.right? key

partial def WBTMap.set (m: WBTMap α β cmp) (key: α) (val: β): WBTMap α β cmp :=
  match m.node? with
    | none => Node.makeNode (key, val) none none
    | some n =>
      let (ekey, _) := n.entry
      match cmp key ekey with
        | Ordering.lt =>
          let l1 := WBTMap.set (cmp := cmp) n.left? key val
          let n1 := Node.makeNode n.entry l1.node? n.right?
          Node.balance Node.δ n1
        | Ordering.eq =>
          let n1 := Node.makeNode (key, val) n.left? n.right?
          Node.balance Node.δ n1
        | Ordering.gt =>
          let r1 := WBTMap.set (cmp := cmp) n.right? key val
          let n1 := Node.makeNode n.entry n.left? r1.node?
          Node.balance Node.δ n1

partial def WBTMap.del? (m: WBTMap α β cmp) (key: α): Option (WBTMap α β cmp) := do
  match m.node? with
    | none => none
    | some n =>
      let (ekey, _) := n.entry
      match cmp key ekey with
        | Ordering.lt =>
          let l1 ← WBTMap.del? (cmp := cmp) n.left? key
          let n1 := Node.makeNode n.entry l1.node? n.right?
          Node.balance Node.δ n1
        | Ordering.eq =>
          match n.right? with
            | none => n.left?
            | some r => -- by default, remove from the right
              let (rMinKey, rMinVal) ← WBTMap.min? (cmp := cmp) r
              let r1 ← WBTMap.del? (α := α) (β := β) (cmp := cmp) r rMinKey
              let n1 := Node.makeNode (rMinKey, rMinVal) n.left? r1.node?
              Node.balance Node.δ n1
        | Ordering.gt =>
          let r1 ← WBTMap.del? (cmp := cmp) n.right? key
          let n1 := Node.makeNode n.entry n.left? r1.node?
          Node.balance Node.δ n1

end WBT
