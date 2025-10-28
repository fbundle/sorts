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







end WBT
