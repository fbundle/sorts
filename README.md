# SORTS
started from some basic dependent type in [lean4](https://github.com/leanprover/lean4), I stated this project with the goal to implement the full dependent type system so that it is capable for mathemtical proof. This is probably a decade-long project, hope it would last.


This repository currently focuses on:
- Basic sort universe machinery with levels and a partial order (subtyping/cumulativity)
- Core type constructors: arrows (function types) and products (pairs)
- A tiny s-expression based front-end (EL) and a Python helper to transpile a Python-like syntax into EL

The vision is to iteratively extend this core into a dependable, minimal kernel with an ergonomic front-end suitable for mathematical proofs.

## Quick start

- Ensure Go is installed (see `go.mod` for version).
- Run the EL interpreter on the provided examples:

```bash
cat example.el | go run cmd/el/main.go
python cmd/elt/elt.py example.elt | go run cmd/el/main.go
```

## Project layout

- `sorts/`: core sort/type constructors and subtyping
  - `atom.go`: atomic sorts and universe chaining
  - `arrow.go`: function types, intro/elim
  - `prod.go`: product types, intro/elim
  - `sorts.go`: shared interfaces and helpers
- `universe/`: universe implementation and builtin rules (initial/terminal chains, cumulativity)
- `form/`, `form_processor/`: tokenization and parsing of s-expressions used by EL
- `cmd/el/`: CLI for interpreting EL input
- `cmd/elt/`: Python helper to transpile a lightweight Python-like syntax into EL
- `example.el`, `example.py`: sample programs

## Current capabilities (work-in-progress)

- Universe levels with initial/terminal chains and basic cumulativity checks
- Subtyping via a partial order on sorts, including function-type variance
- Intro/elim rules for arrows and products (at the sort level)
- EL front-end to construct and inspect terms; Python-like transpiler for convenience

## Roadmap toward a dependent type theory proof assistant

Short-to-medium term:
- Add sum types and unit/empty as proper inhabitants
- Dependent function types (Π-types) and dependent pairs (Σ-types)
- Judgmental equality/definitional equality and normalization for the kernel
- Pattern matching and inductive families (starting with simple inductives)
- Contexts, typing judgments, and well-scoped variables with hygienic names
- Improved parser and elaboration from surface syntax to core terms

Long-term milestones:
- Universe polymorphism and cumulativity scaling
- General inductive types (W-types), recursion/recursors, and termination checking
- Tactic framework and interactive proving experience
- Standard library of datatypes and theorems
- Proof irrelevance/prop hierarchy exploration and potential classical axioms as options

## Contributing

This is an early-stage exploration. Issues, ideas, and small PRs are welcome—especially around tests, examples, and small core extensions. If proposing larger changes (new connectives, typing rules, or parser features), please open an issue first to discuss the design.

## License

TBD. For now, assume all rights reserved by the author until a license is added.

