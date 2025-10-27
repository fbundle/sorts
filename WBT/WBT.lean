namespace WBT

structure Node (α: Type u) where
  weight: Nat
  height: Nat
  entry: α
  left?: Option (Node α)
  right?: Option (Node α)

partial def iterate (n?: Option (Node α)): Array α :=
  match n? with
    | none => #[]
    | some n =>
      let (l, r) := (iterate n.left?, iterate n.right?)
      l ++ #[n.entry] ++ r

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

partial def cmp (δ: Nat) (n: Option (Node α)): Ordering :=
  match n with
    | none => Ordering.eq
    | some n =>
      let (l, r) := (weight n.left?, weight n.right?)
      if l + r ≤ 1 then
        Ordering.eq
      else if (l > δ * r) then
        Ordering.gt -- left heavy
      else if (δ * l < r) then
        Ordering.lt -- right heavy
      else
        Ordering.eq

def rightRotate (n: Node α): Option (Node α) :=
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
  match n.left? with
    | none => none
    | some l =>
      let r? := n.right?
      let (ll?, lr?) := (l.left?, l.right?)

      let n1 := makeNode n.entry lr? r?
      let l1 := makeNode l.entry ll? n1
      l1

def leftRotate (n: Node α): Option (Node α) :=
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
  match n.right? with
    | none => none
    | some r =>
      let l? := n.left?
      let (rl?, rr?) := (r.left?, r.right?)

      let n1 := makeNode n.entry l? rl?
      let r1 := makeNode r.entry n1 rr?
      r1

partial def balance (δ: Nat) (n: Node α): Option (Node α) := do
  -- assuming δ ≥ 3
  -- assuming the two subtrees n.left and n.right are balanced
  -- double rotation is necessary - see `why_double_rotation.jpeg`
  match cmp δ (some n) with
    | Ordering.eq => n
    | Ordering.gt => -- left heavy
      let n1 ← rightRotate n
      if cmp δ (some n1) = Ordering.eq then
        n1
      else
        -- not balanced after one single rotation
        -- because lr too heavy
        -- double rotation effectively split lr in half
        match n.left? with
          | none => none
          | some l =>
            let l1 ← leftRotate l
            let n2 := makeNode n.entry l1 n.right?
            let n3 := rightRotate n2
            n3
    | Ordering.lt => -- right heavy
      let n1 ← leftRotate n
      if cmp δ (some n1) = Ordering.eq then
        n1
      else
        -- not balanced after one single rotation
        -- because rl too heavy
        -- double rotation effectively split rl in half
        match n.right? with
          | none => none
          | some r =>
            let r1 ←  rightRotate r
            let n2 := makeNode n.entry n.left? r1
            let n3 := leftRotate n2
            n3

-- theorem for _balance_once
-- assuming the two subtrees n.left and n.right are balanced
-- with δ ≥ 3, a single rotation is sufficient to make the whole tree balanced
def balanceThm (δ: Nat) (n: Node α):
  δ ≥ 3
  → Ordering.eq = cmp δ n.left?
  → Ordering.eq = cmp δ n.right?
  → (
    match balance δ n with
      | none => false
      | some n1 => Ordering.eq = cmp δ (some n1)
  )
  := sorry

def δ := 3

end WBT
