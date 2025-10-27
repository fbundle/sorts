namespace WBT

structure Node (α: Type u) where
  weight: Nat
  height: Nat
  entry: α
  left?: Option (Node α)
  right?: Option (Node α)

partial def iterate (n?: Option (Node α)): List α :=
  match n? with
    | none => []
    | some n =>
      let (l, r) := (iterate n.left?, iterate n.right?)
      l ++ [n.entry] ++ r

def weight (n?: Option (Node α)): Nat :=
  match n? with
    | none => 0
    | some n => n.weight

def height (n?: Option (Node α)): Nat :=
  match n? with
    | none => 0
    | some n => n.height

def makeNode (entry: α) (left?: Option (Node α)) (right?: Option (Node α)): Node α :=
  {
    weight := 1 + weight left? + weight right?,
    height := 1 + max (weight left?) (weight right?),
    entry := entry,
    left? := left?,
    right? := right?,
  }



def rightRotate (n: Node α) (hl: n.left?.isSome): Node α :=
  -- right rotate
  --         n
  --   l           r
  -- ll lr
  --
  --      becomes
  --
  --         l
  --   ll          n
  --             lr r
  let (l, r?) := (n.left?.get hl, n.right?)
  let (ll?, lr?) := (l.left?, l.right?)

  let n1 := makeNode n.entry lr? r?
  let l1 := makeNode l.entry ll? n1
  l1

def leftRotate (n: Node α) (hr: n.right?.isSome): Node α :=
  -- left rotate
  --         n
  --   l           r
  --             rl rr
  --
  --      becomes
  --
  --         r
  --   n          rr
  --  l rl
  let (l?, r) := (n.left?, n.right?.get hr)
  let (rl?, rr?) := (r.left?, r.right?)

  let n1 := makeNode n.entry l? rl?
  let r1 := makeNode r.entry n1 rr?
  r1







partial def strongCmp (δ: Nat) (n: Option (Node α)): Ordering :=
  match n with
    | none => Ordering.eq
    | some n =>
      let (l, r) := (weight n.left?, weight n.right?)
      if l + r ≤ 1 then
        Ordering.eq
      else if (l > δ * r) then
        Ordering.gt
      else if (δ * l < r) then
        Ordering.lt
      else
        Ordering.eq

partial def weakCmp (δ: Nat) (n: Option (Node α)): Ordering :=
  match n with
    | none => Ordering.eq
    | some n =>
      let (l, r) := (weight n.left?, weight n.right?)
      if l + r ≤ 1 then
        Ordering.eq
      else if (l > δ * r + 1) then
        Ordering.gt
      else if (δ * l + 1 < r) then
        Ordering.lt
      else
        Ordering.eq

partial def rotateAtMostOnce (δ: Nat) (n: Node α): Node α :=
  -- assuming δ ≥ 3
  -- assuming the two subtrees n.left and n.right are balanced
  -- a single rotation is sufficient to make the whole tree balanced
  let (l, r) := (n.left?, n.right?)
  match strongCmp δ (some n) with
    | Ordering.gt =>
      -- right rotate
      --         n
      --   l           r
      -- ll lr
      --
      --      becomes
      --
      --         l
      --   ll          n
      --             lr r
      let l := l.get sorry -- add proof from leftHeavy
      let (ll, lr) := (l.left?, l.right?)
      let n1 := makeNode n.entry lr r
      let n1 := rotateAtMostOnce δ n1
      let l1 := makeNode l.entry ll n1
      l1
    | Ordering.lt =>
      -- left rotate
      --         n
      --   l           r
      --             rl rr
      --
      --      becomes
      --
      --         r
      --   n          rr
      --  l rl
      let r := r.get sorry -- TODO -- add proof from _right_heavy
      let (rl, rr) := (r.left?, r.right?)
      let n1 := makeNode n.entry l rl
      let n1 := rotateAtMostOnce δ n1
      let r1 := makeNode r.entry n1 rr
      r1
    | Ordering.eq => n

-- theorem for _balance_once
-- assuming the two subtrees n.left and n.right are balanced
-- with δ ≥ 3, a single rotation is sufficient to make the whole tree balanced
def balanceThm (δ: Nat) (n: Node α):
  δ ≥ 3
  → Ordering.eq = strongCmp δ n.left?
  → Ordering.eq = strongCmp δ n.right?
  → Ordering.eq = weakCmp δ (some n)
  → Ordering.eq = strongCmp δ (some (rotateAtMostOnce δ n))
  := sorry

end WBT
