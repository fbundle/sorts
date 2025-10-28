import WBT.WBT

namespace WBT



-- TODO potentially let cmp hashed α so that the tree will be more balanced
structure WBTMap (α: Type u) (β: Type v) (cmp: α → α → Ordering) where
  node? : Option (Node (α × β))

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

-- def WBTMap.length (m: WBTMap α β cmp)






end WBT
