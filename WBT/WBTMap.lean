namespace Persistent

structure Node (α: Type u) (β: Type v) (cmp: α → α → Ordering): Type (max u v) where
  weight: Nat
  height: Nat
  key: α
  val: β
  left?: Option (Node α β cmp)
  right?: Option (Node α β cmp)

def weight (n?: Option (Node α β cmp)): Nat :=
  match n? with
    | none => 0
    | some n => n.weight

def height (n?: Option (Node α β cmp)): Nat :=
  match n? with
    | none => 0
    | some n => n.height

partial def balanceCond (δ: Nat) (n?: Option (Node α β cmp)): Ordering :=
  match n? with
    | none => Ordering.eq
    | some n =>
      let (l, r) := (weight n.left?, weight n.right?)
      if l + r < 2 then
        Ordering.eq
      else
        sorry


-- δ ll ≥ lr
-- δ r + 1 ≥ 1 + ll + lr > δ r -- not too imbalanced might be (δ r + something)

-- 1 + ll + lr + r ≤ (δ + 1) r + 1
-- 1 + lr + r ≤ (δ + 1) r - ll + 1
--            < (δ + 1) (1 + ll + lr) / δ - ll + 1
--            ≤ 1 + (δ + 1)/δ + (δ + 1) lr /δ + 1/δ ll
--            ≤ 1 + (δ + 1)/δ + (δ + 1) ll + 1/δ ll
-- choose δ = 3
-- => 1 + lr + r ≤ 3 ll



def makeNode (key: α) (val: β) (left: Option (Node α β cmp)) (right: Option (Node α β cmp)): Node α β cmp :=
  {
    weight := 1 + weight left + weight right,
    height := 1 + max (weight left) (weight right),
    key := key,
    val := val,
    left := left,
    right := right,
  }

partial def rotateAtMostOnce (δ: Nat) (n: Node α β cmp): Node α β cmp :=
  -- assuming δ ≥ 3
  -- assuming the two subtrees n.left and n.right are balanced
  -- a single rotation is sufficient to make the whole tree balanced
  let (l, r) := (n.left, n.right)
  if leftHeavy δ n then
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
    let (ll, lr) := (l.left, l.right)
    let n1 := makeNode n.key n.val lr r
    let n1 := rotateAtMostOnce δ n1
    let l1 := makeNode l.key l.val ll n1
    l1
  else if rightHeavy δ n then
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
    let (rl, rr) := (r.left, r.right)
    let n1 := makeNode n.key n.val l rl
    let n1 := rotateAtMostOnce δ n1
    let r1 := makeNode r.key r.val n1 rr
    r1
  else
    n

-- theorem for _balance_once
-- assuming the two subtrees n.left and n.right are balanced
-- with δ ≥ 3, a single rotation is sufficient to make the whole tree balanced
def balanceThm (δ: Nat) (n: Node α β cmp):
  δ ≥ 3 →
  balanceCond δ n.left → balanceCond δ n.right →
  balanceCond δ (some (rotateAtMostOnce δ n))
  := sorry



def makeNodeBalance (key: α) (val: β) (left: Option (Node α β cmp)) (right: Option (Node α β cmp)): Node α β cmp :=
  let δ := 3
  let n1 := makeNode key val left right
  let n2 := rotateAtMostOnce δ n1
  n2











structure WBTMap (α: Type u) (β: Type v) (cmp: α → α → Ordering): Type (max u v) where
  node: Option (Node α β cmp)








-- all these for this
structure A where
  val: Nat
  map: WBTMap String A compare

end Persistent
