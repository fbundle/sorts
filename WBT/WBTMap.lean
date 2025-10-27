import WBT.WBT

namespace WBT



structure WBTMap (α: Type u) (β: Type v) (cmp: α → α → Ordering) where
  node? : Option (Node (α × β))

-- the whole purpose of this self-balancing tree is to do this
private structure A where
  val : Nat
  -- this only works for List, Array and don't work for Std.HashMap and Lean.RBTree
  map : WBTMap String A compare






end WBT
