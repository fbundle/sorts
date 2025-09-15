# ROADMAP (Detailed Implementation Plan)

This document expands the README roadmap into concrete implementation steps. It focuses on minimal, composable kernel extensions, typed front-end forms, and a pragmatic testing strategy. Each section lists:

- Problem statement and scope
- Data model updates (Go types, interfaces)
- Typing rules and subtyping implications
- Reduction/operational semantics
- Parser/EL forms (and ELT expectations)
- Context/universe interactions
- Logging, errors, and modes
- Tests and examples
- Milestones and incremental PR plan

The goal is to ship features in small, verifiable slices without regressing existing behavior.

## 1) Sum, Unit, and Empty types as first-class sorts

### Problem
Introduce canonical inhabitants for unit `Unit` and empty `Empty`, and add binary sum `A ⊕ B` at the sort level with intro/elim (inl/inr, match). The README already exposes `(⊕ A B)`; we will complete its semantics and evaluation.

### Data model
- In `sorts/` add concrete nodes for sums and canonical terms:
  - `type Sum struct { Atom; Left Sort; Right Sort }`
  - `type Unit struct { Atom }`, with a single canonical inhabitant `unit`.
  - `type Empty struct { Atom }`, with no inhabitants.
  - Intro forms for sum inhabitants:
    - `type Inl struct { Atom; Value Sort; Left Sort; Right Sort }`
    - `type Inr struct { Atom; Value Sort; Left Sort; Right Sort }`

- Universe levels:
  - `GetLevel` of `Sum` = `max(level(A), level(B))`.
  - `Parent(Sum(A,B))` = a universe sort at that level.
  - `Unit`, `Empty` live in universe levels like other atomic types.

### Typing rules
- Unit/Empty:
  - `unit : Unit`.
  - No intro rule for `Empty` inhabitant.
- Sum:
  - If `a : A` then `inl a : A ⊕ B`.
  - If `b : B` then `inr b : A ⊕ B`.
  - Elim (match): for `v : A ⊕ B` and branches `x:A ⊢ t1 : C` and `y:B ⊢ t2 : C`, result type is `C`. Implementation chooses LUB of branch types as current kernel does; keep this but ensure branches are checked in extended context.
- Subtyping:
  - Covariant in both arguments: `A1 ≤ A2 ∧ B1 ≤ B2 ⇒ A1 ⊕ B1 ≤ A2 ⊕ B2`.

### Reduction
- In `ModeEval`:
  - `match (inl v) (=> (inl x) t1) (=> (inr y) t2) ...` reduces to `[x := v] t1`.
  - Likewise for `inr`.
  - If default `_` branch exists, keep current behavior as fallback.

### Parser/EL
- Extend list compilers:
  - `(⊕ A B)` → `Sum{Left:A, Right:B}` already registered; verify typing.
  - `(inl A B v)` and `(inr A B v)` as explicit forms. Also allow sugar `(inl v (⊕ A B))` later if useful; first implement explicit arity to keep typing simple.
  - Reuse existing `(match ...)` with `=>` cases; add pattern constructors `inl`, `inr`, and `_`.

### Context/universe
- Add builtins `Unit`, `Empty` to `universe.GetBuiltin`.
- Provide `unit` as builtin inhabitant with parent `Unit`.

### Logging, errors, modes
- Type errors: `type_check` if branch type mismatches or intro value is not subtype of the declared side.
- `ToString` should show form and type when in `ModeEval`/`ModeDebug`.

### Tests/examples
- EL:
  ```el
  (inspect unit)                        # (form unit - type Unit - level 0)
  (match (inl A B x)
    (=> (inl a) a)
    (=> (inr b) (inh A)))               # type checks via LUB(A, A)
  ```
- Negative tests: constructing `(inl A B v)` when `v` not ≤ `A`.

### Milestones
- PR1: Introduce `Unit`, `Empty`, `unit` builtin and printing.
- PR2: Implement `Sum` sort typing and subtyping; compile `(⊕ A B)`.
- PR3: Add `inl`/`inr` intro nodes with typing and printing.
- PR4: Extend `match` reduction for `inl/inr` in `ModeEval`.

---

## 2) Dependent function (Π) and dependent pair (Σ)

### Problem
Add Π and Σ as dependent generalizations of `->` and product. Keep minimal kernel semantics; no general recursion.

### Data model
- `type Pi struct { Atom; Param TypeAnnot; BodyType Sort }`
  - `Param: Name : A` where `A` is the domain.
  - `BodyType` is the type of the result depending on `Param.Name`.
- `type Sigma struct { Atom; Param TypeAnnot; BodyType Sort }`
- Intro/elim:
  - Π-intro is `Lambda` we already have; enrich typing to allow dependent codomain.
  - Π-elim is application; extend beta to substitute argument into body type and term.
  - Σ-intro: `Pair{ First Sort; Second Sort; Param TypeAnnot; BodyType Sort }` where `first : A` and `second : BodyType[x := first]`.
  - Σ-elim: `fst`, `snd` projections with appropriate parents.

### Typing rules
- Π:
  - If `x:A ⊢ B(x) : Type l` then `Π (x:A). B(x) : Type max(level(A), level(B))`.
  - `λ (x:A). t(x) : Π (x:A). B(x)` when `x:A ⊢ t(x) : B(x)`.
  - Application: if `f : Π (x:A). B(x)` and `a : A'` with `A' ≤ A`, then `f a : B(a)`.
- Σ:
  - If `x:A ⊢ B(x) : Type l` then `Σ (x:A). B(x) : Type max(level(A), level(B))`.
  - Pair: if `a:A` and `b:B(a)`, then `(a,b) : Σ (x:A). B(x)`.
  - Projections: `fst : Σ (x:A). B(x) -> A`; `snd : Π (p:Σ (x:A). B(x)). B(fst p)`.

### Reduction
- Π: `(λ x. t) a  ↦  t[x := a]` (β). Capture-avoiding substitution required.
- Σ: `fst (a,b) ↦ a`, `snd (a,b) ↦ b`.

### Parser/EL
- Add compilers:
  - `(Π (: x A) B)` → `Pi{Param:x:A, BodyType:B}`. Infix sugar: `{(: x A) Π B}` optional later.
  - `(Σ (: x A) B)` → `Sigma{...}`.
  - Pairs: `(pair a b (: x A) B)` initially explicit. Projections: `(fst p)`, `(snd p)`.
  - Extend `(=> (: x A) body)` to allow dependent result typing: when computing parent for `Lambda`, set to `Π (: x A) type(body)`.

### Context/universe
- Context extension during typing: when checking `B`, use `ctx.Set(x, NewTerm(x, A))`.
- Subtyping for Π, Σ:
  - Π is contravariant in domain and covariant in codomain.
  - Σ is covariant in both domain and codomain.

### Implementation details
- Add a small, hygienic substitution utility in `sorts/util.go` that:
  - Takes a `Sort` term, a `form.Name` and a replacement `Sort`, returns a new `Sort` with bound occurrences replaced, respecting binders.
  - Track binders via node types that introduce names (`Lambda`, `Pi`, `Sigma`, `Let`, `Match` patterns).
- Ensure `ctx.ToString` shows Π/Σ forms.

### Tests/examples
- Π:
  ```el
  (=> (: x A) (=> (: y (B x)) y))        # Π x:A. Π y:B x. B x
  ((=> (: x A) x) a)                      # type is A, β-reduces in eval mode
  ```
- Σ:
  ```el
  (pair a b (: x A) (B x))
  (fst (pair a b (: x A) (B x)))          # reduces to a
  ```

### Milestones
- PR1: Define `Pi`/`Sigma` types and printing, without reduction.
- PR2: Extend `Lambda` typing to produce `Pi` parents when body depends on param.
- PR3: Add application β-substitution support; implement capture-safe substitution.
- PR4: Implement `pair`, `fst`, `snd` with typing and eval.

---

## 3) Definitional equality and normalization

### Problem
Introduce judgmental equality for types/terms and a lightweight Normalization by Evaluation (NbE) scaffold to decide convertibility where needed (e.g., during type checking for Π application, branch type comparison, etc.).

### Design
- Keep the current fast path (structural checks). Fall back to NbE when shapes differ but may normalize to equal forms.

### Data model
- Add an `Equal(ctx, x, y Sort) bool` on `Context` delegating to `universe` policy.
- Introduce `sorts/normalize.go`:
  - `Normalize(ctx, t Sort) Sort` reduces β/let/match to weak head normal form (WHNF) first; optionally full normal form for small terms.
  - `AlphaEq(x, y Sort) bool` for alpha-equivalence on binders.

### Typing integration
- Replace direct `LessEqual`/parent comparisons in sensitive sites with `Equal` or `Convertible(ctx, x, y)` that uses normalization.
- Example: in application, allow `A'` convertible to `A` rather than merely `≤`.

### Reduction
- Complete `Reduce` for `Beta`, `Let`, `Match` to fuel normalization.

### Tests
- Equality across β-redexes, let-inlining, trivial match on constructors.
- Non-equal counterexamples.

### Milestones
- PR1: Implement WHNF with β and let.
- PR2: Add match-constructor folding for sums/unit.
- PR3: Alpha-equivalence for binders and paths that need it.

---

## 4) Pattern matching and small inductive families

### Problem
Extend `match` beyond sums to handle small first inductives (e.g., `Bool`, `Nat`) with constructors and simple eliminators.

### Data model
- Add builtins:
  - `Bool` with `true`, `false`.
  - `Nat` with `zero`, `succ`.
- Patterns:
  - Extend `MatchCase` to allow constructor patterns with binders, e.g., `(=> (succ n) t)`.

### Typing rules
- For `Bool`, branches must yield a common type (use LUB or stronger equal requirement under definitional equality).
- For `Nat`, allow primitive recursion later; initially only case analysis.

### Reduction
- In `ModeEval`, perform constructor-directed reduction.

### Parser/EL
- Compilers: `(Bool)`, `(Nat)`, constructors as names, and `(match ...)` extended.

### Tests
- `(match true (=> true a) (=> false b))`.
- `(match (succ zero) (=> zero a) (=> (succ n) (f n)))`.

### Milestones
- PR1: `Bool` + match reduction.
- PR2: `Nat` + non-recursive match.
- PR3: Optional: primitive recursion with fuelled evaluator only.

---

## 5) Contexts, typing judgments, hygienic names

### Problem
Strengthen context hygiene and error reporting; ensure capture-avoiding substitution and shadowing rules are explicit.

### Plan
- Introduce a `NameSupply` that generates fresh `form.Name` with source-span hints.
- Track binder scopes in nodes to support safe substitution and pretty-printing.
- Improve errors with source forms via `ctx.Form` in panics.

### Milestones
- PR1: Name supply and freshening utility, integrate with `Lambda`, `Pi`, `Sigma`, `Let`, `Match`.
- PR2: Replace ad-hoc parameter naming with supply-driven names.

---

## 6) Parser and elaboration improvements

### Problem
Make surface syntax more ergonomic while keeping a small core.

### Plan
- ELT: add sugar for tuples, `fun x : A => t`, `{A × B}`.
- EL: optional sugar for `(inl v (⊕ A B))`, `(inr ...)`.
- Elaboration: desugar sugar into core EL before compile.

### Milestones
- PR1: ELT additions with roundtrip tests.
- PR2: EL sugar for sum constructors.

---

## 7) Complete reductions in ModeEval and Debug

### Problem
Currently `Reduce` is TODO for several nodes.

### Plan
- Implement `Beta.Reduce`: evaluate `cmd` to WHNF; if lambda, substitute.
- Implement `Let.Reduce`: evaluate bindings lazily by substitution into `Final`.
- Implement `Match.Reduce`: evaluate `Cond` to WHNF; if constructor, pick branch.
- Ensure `ToString` in `ModeEval` reflects reduced form with type.

### Tests
- Golden tests comparing compile vs eval output across examples.

---

## 8) Universe scaling (future)

### Problem
Prepare for universe polymorphism and cumulative hierarchies.

### Plan
- Parameterize constructors with universe levels where necessary.
- Provide `lift` operations and ensure `LessEqualBasic` handles level shifts.

---

## Testing strategy

- Unit tests per constructor and reducer in `sorts/` using small EL snippets compiled through `el.Context`.
- Parser roundtrip tests for infix and block structures.
- Negative tests: ill-typed programs must panic with consistent error strings.
- Golden tests for CLI output in compile/eval/debug modes.

## Incremental delivery

Each subsection is broken into 2–4 PRs with full tests, no cross-feature coupling. Favor small interfaces and deterministic printing to keep diffs reviewable.
