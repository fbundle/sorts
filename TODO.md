## TODO

### Quick wins
- Remove unused constant `dataLevel` in `cmd/sort_v2/main.go`.
- Guard `printSorts` against `nil` parent: avoid `Parent().Name()` panic.
- Add brief doc comments for `weakSort` and `strongSort` explaining semantics.

### Ergonomics and safety
- Replace panic-based control flow in helpers with error returns:
  - `chain` → return `(sorts.Sort, error)` (e.g., `composeArrows`).
  - `mustAtom` → return `(sorts.Sort, error)` (e.g., `atomOrError`).
  - Handle errors in `main` and fail fast with context.
- Consider validating `level` inputs if the API expects constrained levels.

### Performance/structure (low priority)
- Cache `any` and `unit` atoms once per `SortSystem` and pass into helpers to avoid repeated lookups.
- Optional micro-optimization: in `weakSort`/`strongSort`, preallocate and assign by index instead of `append` in loops.

### Naming
- Consider renaming:
  - `chain` → `composeArrows` (or `compose`) to clarify intent.
  - `mustAtom` → `atomOrError` if switching to error returns.

### Testing
- Add `_test.go` covering:
  - Cast rule: `AddRule("bool", "int")` implies `bool ≤ int` and not vice versa.
  - `weakSort`/`strongSort` relationships against examples like `int→int→int`.
  - Error paths for failed compositions (once helpers return errors).

### CLI output
- Keep `fmt.Printf` for demo, but consider `log` or a `-v` flag if this evolves into a tool.

---

## Roadmap: Extend `sorts_v2` into a (practical) dependent type theory

### Phase 0: Groundwork and architecture
- Define a core syntax for terms and types: variables, applications, lambdas, Pi, Sigma, universes, literals.
- Introduce `Context` (Γ) with typed bindings, scoping, and lookup APIs.
- Implement capture-avoiding substitution and De Bruijn indices or a safe name-layer.
- Encode sorts/universes as a hierarchy: `Type0 : Type1 : Type2 ...` with cumulativity.

### Phase 1: Typing and convertibility
- Implement typing judgment Γ ⊢ t : A covering:
  - Variables, λ, application, Π, Σ, literals, and annotations.
  - Universe formation/levels and cumulativity (Γ ⊢ A : Type i ⇒ Γ ⊢ A : Type j for i ≤ j).
  - Leverage existing `Level() int` on sorts to enforce universe rules and detect mismatches early.
- Implement definitional equality (convertibility): β-reduction, η (for functions), δ (unfold definable constants when needed).
- Add a normalizer (weak-head first) and a decision procedure for convertibility.

### Phase 2: Elaboration and inference
- Add implicit arguments and instance/auto arguments where appropriate.
- Implement metavariables (holes) and a constraint store.
- Add higher-order unification (pattern-fragment first, fallback heuristics later).
- Elaboration pipeline: surface syntax → constraints + metas → solve → core terms.

### Phase 3: Inductive families and pattern matching
- Add inductive type declarations (simple → indexed families): Nat, Vec, Fin, List.
- Positivity checker for inductive definitions.
- Derive recursors/eliminators and compile pattern matching to eliminators.
- Add computation rules for eliminators (β-like rules) to convertibility.

### Phase 4: Equality and rewriting
- Add identity type (Martin-Löf equality) with `refl` and `J` eliminator.
- Enable `rewrite`/transport by equalities in the checker/elaborator.
- Optional: Propositional irrelevance flag for specific universes.

### Phase 5: Safety and termination
- Termination checker for recursive definitions (structural or sized types-lite).
- Coverage checker for pattern matches.
- Universe polymorphism and level inference with constraints.

### Phase 6: Usability and libraries
- Pretty-printer for terms, types, and contexts; error messages with ranges and hints.
- Prelude: Bool, Nat, List, Maybe/Option, Sigma/Pi helpers, equality lemmas, vectors.
- Examples: length-indexed vectors, safe head, `map` fusion, simple proofs.

### Engineering details to integrate with `sorts_v2`
- Unify `SortSystem` with universe levels: expose APIs for creating `Type i` and coercions (cumulativity) as rules.
- Ensure `Pi`/`Sigma` construction in `sorts_v2` validates domain/codomain inhabitation and universe levels.
- Represent terms distinctly from sorts; add a core AST package: `core/term.go`, `core/typechecker.go`.
- Provide an `env`/`context` module with persistent data structures for immutability.
- Isolate reduction/equality in `core/conv.go` with configurable unfolding (fuel/flags).
 - Use `Sort.Level() int` consistently:
   - Define `Type(i)` as sorts with `Level()==i` and ensure `i<j ⇒ Type(i) ≤ Type(j)` via cumulativity.
   - In `Pi(A,B)`, compute `level := max(Level(A), Level(B))` (with lifting for dependency) and set the resulting level accordingly.
   - Add checks so user-visible constructors cannot fabricate invalid `Level()` values.

### Testing milestones
- Unit tests for substitution, normalization, convertibility, and type formation.
- Property tests for level monotonicity and cumulativity using `Level()`.
- Golden tests for elaboration of examples with implicit arguments and holes.
- Inductive definitions validated by positivity/termination; proofs using identity type and rewriting.

### Stretch goals
- Records (Σ with field names) and projections with η for records.
- Quotients or HITs (later): start with setoids and rewriting by relations.
- Tactic-style proof scripting for small automation.


