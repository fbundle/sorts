namespace Persistent

def Pos := Nat

private structure _node (α: Type u) where
  weight: Nat -- TODO change to Pos as weight is positive
  height: Nat -- TODO change to Pos as height is positive
  entry: α
  left: Option (_node α)
  right: Option (_node α)

private def _weight (n: Option (_node α)): Nat :=
  match n with
    | none => 0
    | some n => n.weight

private def _height (n: Option (_node α)): Nat :=
  match n with
    | none => 0
    | some n => n.height


-- balanced condition for tree with more than 3 values
private def _left_heavy (δ: Nat) (n: _node α): Bool :=
  let (l, r) := (_weight n.left, _weight n.right)
  (l + r ≥ 2) ∧ (l > δ * r)

private def _right_heavy (δ: Nat) (n: _node α) : Bool :=
  let (l, r) := (_weight n.left, _weight n.right)
  (l + r ≥ 2) ∧ (δ * l < r)

-- balanced tree condition
private partial def _balance_cond (δ: Nat) (n: Option (_node α)): Bool :=
  match n with
    | none => True
    | some n =>
      ¬ ((_left_heavy δ n) ∨ (_right_heavy δ n))
      ∧
      _balance_cond δ n.left ∧ _balance_cond δ n.right


private partial def _make_node (e: α) (l: Option (_node α)) (r: Option (_node α)) : _node α :=
  let _make_node_without_balance (e: α) (l: Option (_node α)) (r: Option (_node α)) : _node α :=
  {
    weight := 1 + _weight l + _weight r,
    height := 1 + max (_height l) (_height r),
    entry := e,
    left := l,
    right := r,
  }

  let rec _balance (δ: Nat) (n: _node α): _node α :=
  -- assuming δ ≥ 3
  -- assuming the two subtrees n.left and n.right are balanced
  -- a single rotation is sufficient to make the whole tree balanced
  let (l, r) := (n.left,  n.right)
  if _left_heavy δ n then
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
    let l := l.get sorry -- TODO - add proof from _left_heavy
    let (ll, lr) := (l.left, l.right)
    let n1 := _make_node_without_balance n.entry lr r
    let n1 := _balance δ n1
    let l1 := _make_node_without_balance l.entry ll n1
    l1
  else if _right_heavy δ n then
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
    let n1 := _make_node_without_balance n.entry l rl
    let n1 := _balance δ n1
    let r1 := _make_node_without_balance r.entry n1 rr
    r1
  else
    n

  -- theorem for _balance_once
  -- assuming the two subtrees n.left and n.right are balanced
  -- with δ ≥ 3, a single rotation is sufficient to make the whole tree balanced
  let _balance_thm (δ: Nat) (n: _node α):
    δ ≥ 3 →
    _balance_cond δ n.left → _balance_cond δ n.right →
    _balance_cond δ (some (_balance δ n))
    := sorry

  let δ := 3
  let n1 := _make_node_without_balance e l r
  let n2 := _balance δ n1

  n2

private def _empty {α}: Option (_node α) := none
private def _singleton (e: α): _node α := _make_node e none none

private partial def _get (n: Option (_node α)) (i: Nat) : Option α :=
  match n with
    | none => none
    | some n =>
      let leftWeight := _weight n.left
      if i < leftWeight then
        _get n.left i
      else if i < leftWeight + 1 then
        n.entry
      else
        _get n.right (i - (leftWeight + 1))

private partial def _set (n: Option (_node α)) (i: Nat) (e: α) : Option (_node α) :=
  match n with
    | none => none
    | some n =>
      let leftWeight := _weight n.left
      if i < leftWeight then
        let l1 := _set n.left i e
        let n1 := _make_node n.entry l1 n.right
        n1
      else if i = leftWeight then
        let n1 := _make_node e n.left n.right
        n1
      else
        let r1 := _set n.right (i - (leftWeight + 1)) e
        let n1 := _make_node n.entry n.left r1
        n1

private partial def _ins (n: Option (_node α)) (i: Nat) (e: α) : Option (_node α) :=
  match n with
    | none =>
      if i = 0 then
        _singleton e
      else
        none
    | some n =>
      let leftWeight := _weight n.left
      if i < leftWeight then
        let l1 := _ins n.left i e
        let n1 := _make_node n.entry l1 n.right
        n1
      else if i = leftWeight then
        let r1 := _ins n.right 0 n.entry
        let n1 := _make_node e n.left r1
        n1
      else
        let r1 := _ins n.right (i - (leftWeight + 1)) e
        let n1 := _make_node n.entry n.left r1
        n1


private partial def _del (n: Option (_node α)) (i: Nat) : Option (_node α) :=
  match n with
    | none => none
    | some n =>
      let leftWeight := _weight n.left
      if i < leftWeight then
        let l1 := _del n.left i
        let n1 := _make_node n.entry l1 n.right
        n1
      else if i = leftWeight then
        let e := _get n.right 0
        match e with
          | none => n.left
          | some e =>
            let r1 := _del n.right 0
            let n1 := _make_node e n.left r1
            n1
      else
        let r1 := _del n.right (i - (leftWeight + 1))
        let n1 := _make_node n.entry n.left r1
        n1

private partial def _to_array (n: Option (_node α)) : Array α :=
  let rec loop (a: Array α) (i: Nat) (n: Option (_node α)) : Array α :=
    let e := _get n i
    match e with
      | none => a
      | some e => loop (a.push e) (i+1) n

  loop Array.empty 0 n

private partial def _from_array (a: Array α): Option (_node α) :=
  let rec loop (n: Option (_node α)) (i: Nat) (a: Array α): Option (_node α) :=
    match a[i]? with
      | none => n
      | some x => loop (_ins n (_weight n) x) (i+1) a

  loop none 0 a

-- ArrayWBT
structure ArrayWBT (α: Type u) where
  node: Option (_node α)

def singleton (e: α): ArrayWBT α := {node := _singleton e}
def empty {α: Type u}: ArrayWBT α := {node := none}
def fromArray (a: Array α): ArrayWBT α := {node := _from_array a}
def ArrayWBT.toArray (s: ArrayWBT α) : Array α := _to_array s.node

def ArrayWBT.length (s: ArrayWBT α) := _weight s.node
def ArrayWBT.depth (s: ArrayWBT α) := _height s.node

def ArrayWBT.get (s: ArrayWBT α) (i: Nat) : Option α := _get s.node i
def ArrayWBT.set (s: ArrayWBT α) (i: Nat) (e: α): ArrayWBT α := {node := _set s.node i e}
def ArrayWBT.insert (s: ArrayWBT α) (i: Nat) (e: α): ArrayWBT α := {node := _ins s.node i e}
def ArrayWBT.delete (s: ArrayWBT α) (i: Nat) : ArrayWBT α := {node := _del s.node i}
def ArrayWBT.append (s: ArrayWBT α) (e: α) : ArrayWBT α := s.insert s.length e

-- ArrayWBT ToString
--instance [ToString α] : ToString (ArrayWBT α) where
  -- toString (a: ArrayWBT α): String := toString a.toArray


-- ArrayWBT is applicative functor
private partial def _map (f: α → β) (n: Option (_node α)): Option (_node β) :=
  match n with
    | none => none
    | some n => _make_node (f n.entry) (_map f n.left) (_map f n.right)

def ArrayWBT.map (s: ArrayWBT α) (f: α → β): ArrayWBT β := {node := _map f s.node}

instance : Functor ArrayWBT where
  map := λ f s => s.map f


def _z2 :=
  let z := fromArray (Array.replicate 1024 1)
  z

#eval (_z2.length, _z2.depth)
#eval _balance_cond 3 _z2.node

end Persistent
