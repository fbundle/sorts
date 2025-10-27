import WBT.WBT

namespace WBT

structure WBTMap (α: Type u) (β: Type v) (cmp: α → α → Ordering): Type (max u v) where
  node?: Option (Node (α × β))

def WBTMap.empty {α: Type u} {β: Type v} {cmp: α → α → Ordering}: WBTMap α β cmp := {node? := none}

def WBTMap.entries (map: WBTMap α β cmp): List (α × β) := iterate map.node?

partial def WBTMap.get? (map: WBTMap α β cmp) (key: α): Option β := do
  let node ← map.node?
  let (eKey, eVal) := node.entry
  match cmp key eKey with
    | Ordering.eq => pure eVal
    | Ordering.lt => WBTMap.get? (cmp := cmp) {node? := node.left?} key
    | Ordering.gt => WBTMap.get? (cmp := cmp) {node? := node.right?} key

partial def WBTMap.set (map: WBTMap α β cmp) (key: α) (val: β): WBTMap α β cmp :=
  match map.node? with
    | none =>
      {node? := makeNode (key, val) none none}
    | some node =>
      let (eKey, _) := node.entry
      match cmp key eKey with
        | Ordering.eq =>
          {node? := makeNode (key, val) node.left? node.right?}
        | Ordering.lt =>
          let newLeft := (WBTMap.set (cmp := cmp) {node? := node.left?} key val).node?
          {node? := makeNode node.entry newLeft node.right?}
        | Ordering.gt =>
          let newRight := (WBTMap.set (cmp := cmp) {node? := node.right?} key val).node?
          {node? := makeNode node.entry node.left? newRight}

partial def WBTMap.del (map: WBTMap α β cmp) (key: α) (val: β): WBTMap α β cmp :=
  match map.node? with
    | none => {node? := none}
    | some node =>ßß
      sorry

-- all these for this
structure A where
  val: Nat
  map: WBTMap String A compare

def x :=
  let y : WBTMap Nat String compare := WBTMap.empty
  let y := y.set 123 "123"
  let y := y.set 456 "123"
  let y := y.set 456 "123"
  y

#eval x.entries

end WBT
