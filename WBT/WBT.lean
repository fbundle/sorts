namespace Persistent

structure Node (α: Type u) where
  weight: Nat
  height: Nat
  entry: α
  left: Option (Node α)
  right: Option (Node α)

def weight (n: Option (Node α)): Nat :=
  match n with
    | none => 0
    | some n => n.weight

def height (n: Option (Node α)): Nat :=
  match n with
    | none => 0
    | some n => n.height


def makeNode (entry: α) (left: Option (Node α)) (right: Option (Node α)): Node α :=
  {
    weight := 1 + weight left + weight right,
    height := 1 + max (weight left) (weight right),
    entry := entry,
    left := left,
    right := right,
  }

def leftHeavy (n: Node α): Bool :=
  let (l, r) := (weight n.left, weight n.right)
  (l + r ≥ 2) ∧ (l ≥ 3 * r + 1)

def rightHeavy (n: Node α): Bool :=
  let (l, r) := (weight n.left, weight n.right)
  (l + r ≥ 2) ∧ (3 * l + 1 ≤ r)

partial def rotateAtMostOnce (n: Node α): Node α :=
  -- assuming δ = 3
  -- assuming the two subtrees n.left and n.right are balanced
  -- a single rotation is sufficient to make the whole tree balanced
  let (l, r) := (n.left, n.right)
  if leftHeavy n then
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
    let n1 := makeNode n.entry lr r
    let n1 := rotateAtMostOnce n1
    let l1 := makeNode l.entry ll n1
    l1
  else if rightHeavy n then
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
    let n1 := makeNode n.entry l rl
    let n1 := rotateAtMostOnce δ n1
    let r1 := makeNode r.key r.val n1 rr
    r1
  else
    n

-- theorem for _balance_once
-- assuming the two subtrees n.left and n.right are balanced
-- with δ ≥ 3, a single rotation is sufficient to make the whole tree balanced
def balanceThm (δ: Nat) (n: Node α):
  δ ≥ 3 →
  balanceCond δ n.left → balanceCond δ n.right →
  balanceCond δ (some (rotateAtMostOnce δ n))
  := sorry



def makeNodeBalance (key: α) (val: β) (left: Option (Node α)) (right: Option (Node α)): Node α :=
  let δ := 3
  let n1 := makeNode key val left right
  let n2 := rotateAtMostOnce δ n1
  n2











structure WBTMap (α: Type u) (β: Type v) (cmp: α → α → Ordering): Type (max u v) where
  node: Option (Node α)








-- all these for this
structure A where
  val: Nat
  map: WBTMap String A compare

end Persistent
