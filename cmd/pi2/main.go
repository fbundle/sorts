// skeleton_dep_compiler.go
// Minimal skeleton of a dependent-type-aware compiler in Go.
// Features:
// - AST for terms: variables, lambda, application, Pi (dependent function), Nat (zero/succ), NatRec
// - Simple parser for a tiny surface syntax (very limited)
// - A naive evaluator (normalization by evaluation not implemented, but we do reduction)
// - A basic typechecker that performs conversion checking by normalizing terms
// - Example usage in main: a couple of terms and typechecking

package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// -------------------- AST --------------------

type Term interface {
	String() string
}

// Variables (named)
type Var struct{ Name string }

func (v Var) String() string { return v.Name }

// Lambda: \x : A. body
type Lam struct {
	Param     string
	ParamType Term
	Body      Term
}

func (l Lam) String() string {
	return fmt.Sprintf("(\\%s:%s. %s)", l.Param, l.ParamType.String(), l.Body.String())
}

// Application: f a
type App struct {
	Fun Term
	Arg Term
}

func (a App) String() string { return fmt.Sprintf("(%s %s)", a.Fun.String(), a.Arg.String()) }

// Pi-type: Pi (x : A). B(x)
type Pi struct {
	Param     string
	ParamType Term
	Return    Term
}

func (p Pi) String() string {
	return fmt.Sprintf("(Pi (%s:%s). %s)", p.Param, p.ParamType.String(), p.Return.String())
}

// Universe: Type0, Type1 etc (we'll only use Type0)
type Univ struct{ Level int }

func (u Univ) String() string { return fmt.Sprintf("IType%d", u.Level) }

// Nat constructors
type Zero struct{}

func (z Zero) String() string { return "Zero" }

type Succ struct{ N Term }

func (s Succ) String() string { return fmt.Sprintf("(Succ %s)", s.N.String()) }

// Nat elimination (dependent recursor)
// NatRec P z_case s_case n  -- P : Nat -> IType, z_case : P Zero, s_case : forall (k : Nat). P k -> P (Succ k), n : Nat
type NatRec struct {
	P Term
	Z Term
	S Term
	N Term
}

func (r NatRec) String() string {
	return fmt.Sprintf("(NatRec %s %s %s %s)", r.P.String(), r.Z.String(), r.S.String(), r.N.String())
}

// Annotation (term : type)
type Ann struct {
	Term Term
	Type Term
}

func (a Ann) String() string { return fmt.Sprintf("(%s : %s)", a.Term.String(), a.Type.String()) }

// -------------------- Environment and Utilities --------------------

type Env map[string]Term

func (e Env) Copy() Env {
	ne := make(Env)
	for k, v := range e {
		ne[k] = v
	}
	return ne
}

// Substitute occurrences of a variable name with a term (naive capture-unsafe substitution)
// Note: In a production compiler use De Bruijn indices to avoid capture. This is a simple skeleton.
func subst(t Term, name string, val Term) Term {
	switch u := t.(type) {
	case Var:
		if u.Name == name {
			return val
		}
		return u
	case Lam:
		if u.Param == name { // shadowed
			return u
		}
		return Lam{Param: u.Param, ParamType: subst(u.ParamType, name, val), Body: subst(u.Body, name, val)}
	case App:
		return App{Fun: subst(u.Fun, name, val), Arg: subst(u.Arg, name, val)}
	case Pi:
		if u.Param == name {
			return u
		}
		return Pi{Param: u.Param, ParamType: subst(u.ParamType, name, val), Return: subst(u.Return, name, val)}
	case Univ:
		return u
	case Zero:
		return u
	case Succ:
		return Succ{N: subst(u.N, name, val)}
	case NatRec:
		return NatRec{P: subst(u.P, name, val), Z: subst(u.Z, name, val), S: subst(u.S, name, val), N: subst(u.N, name, val)}
	case Ann:
		return Ann{Term: subst(u.Term, name, val), Type: subst(u.Type, name, val)}
	default:
		return t
	}
}

// -------------------- Evaluator (small-step -> big-step) --------------------

// isValue: for our system, lambdas, zero, succ of value are values
func isValue(t Term) bool {
	switch u := t.(type) {
	case Lam:
		return true
	case Zero:
		return true
	case Succ:
		return isValue(u.N)
	default:
		return false
	}
}

// eval: naive evaluator that reduces beta and natrec where possible. Performs full normalization by repeated evaluation.
func eval(env Env, t Term) Term {
	prev := t
	for {
		next := evalOnce(env, prev)
		if fmt.Sprintf("%#v", next) == fmt.Sprintf("%#v", prev) {
			return next
		}
		prev = next
	}
}

func evalOnce(env Env, t Term) Term {
	switch u := t.(type) {
	case Var:
		if v, ok := env[u.Name]; ok {
			return v
		}
		return u
	case Ann:
		// erase annotation
		return evalOnce(env, u.Term)
	case App:
		f := eval(env, u.Fun)
		a := eval(env, u.Arg)
		if lam, ok := f.(Lam); ok {
			// beta-reduce (capture-unsafe)
			return subst(lam.Body, lam.Param, a)
		}
		return App{Fun: f, Arg: a}
	case Lam:
		// normalize type and body
		return Lam{Param: u.Param, ParamType: eval(env, u.ParamType), Body: eval(env, u.Body)}
	case Pi:
		return Pi{Param: u.Param, ParamType: eval(env, u.ParamType), Return: eval(env, u.Return)}
	case Succ:
		return Succ{N: eval(env, u.N)}
	case NatRec:
		n := eval(env, u.N)
		switch nn := n.(type) {
		case Zero:
			// return Z
			return eval(env, u.Z)
		case Succ:
			// evaluate S applied to predecessor and recursive call
			k := nn.N
			// recursive call: NatRec P Z S k
			rec := NatRec{P: u.P, Z: u.Z, S: u.S, N: k}
			recval := eval(env, rec)
			// apply S to k and recval: S k recval
			return eval(env, App{Fun: App{Fun: u.S, Arg: k}, Arg: recval})
		default:
			return NatRec{P: eval(env, u.P), Z: eval(env, u.Z), S: eval(env, u.S), N: n}
		}
	default:
		return t
	}
}

// -------------------- Typechecker --------------------

// check returns the type of term t under context ctx (maps names to types). It also needs a value environment to normalize definitions.
func typeCheck(ctx Env, venv Env, t Term) (Term, error) {
	switch u := t.(type) {
	case Var:
		typ, ok := ctx[u.Name]
		if !ok {
			return nil, fmt.Errorf("unbound variable %s", u.Name)
		}
		return typ, nil
	case Ann:
		// check that annotation is a valid type
		typOfType, err := typeCheck(ctx, venv, u.Type)
		if err != nil {
			return nil, err
		}
		// For simplicity ensure annotation is a universe
		if _, ok := typOfType.(Univ); !ok {
			return nil, fmt.Errorf("annotation is not a type (universe)")
		}
		// check term against annotation
		err = checkAgainst(ctx, venv, u.Term, u.Type)
		if err != nil {
			return nil, err
		}
		return u.Type, nil
	case Lam:
		// ParamType must be a type
		ptyp, err := typeCheck(ctx, venv, u.ParamType)
		if err != nil {
			return nil, err
		}
		if _, ok := ptyp.(Univ); !ok {
			return nil, fmt.Errorf("parameter type not a universe: %s", ptyp.String())
		}
		// extend context
		nctx := ctx.Copy()
		nctx[u.Param] = u.ParamType
		// If body type T, then lam has type Pi
		bodyType, err := typeCheck(nctx, venv, u.Body)
		if err != nil {
			return nil, err
		}
		return Pi{Param: u.Param, ParamType: u.ParamType, Return: bodyType}, nil
	case App:
		ftyp, err := typeCheck(ctx, venv, u.Fun)
		if err != nil {
			return nil, err
		}
		// normalize ftyp
		ftypn := eval(venv, ftyp)
		if pi, ok := ftypn.(Pi); ok {
			// check argument against pi.ParamType
			if err := checkAgainst(ctx, venv, u.Arg, pi.ParamType); err != nil {
				return nil, err
			}
			// return result with substitution
			argEval := eval(venv, u.Arg)
			res := subst(pi.Return, pi.Param, argEval)
			return res, nil
		}
		return nil, fmt.Errorf("application of non-function: %s", ftypn.String())
	case Pi:
		// param type must be a universe
		ptyp, err := typeCheck(ctx, venv, u.ParamType)
		if err != nil {
			return nil, err
		}
		if _, ok := ptyp.(Univ); !ok {
			return nil, fmt.Errorf("Pi parameter not a type")
		}
		// extend ctx and check return is a universe
		nctx := ctx.Copy()
		nctx[u.Param] = u.ParamType
		rtyp, err := typeCheck(nctx, venv, u.Return)
		if err != nil {
			return nil, err
		}
		if _, ok := rtyp.(Univ); !ok {
			return nil, fmt.Errorf("Pi return not a universe")
		}
		// universe level handling omitted; return Type0
		return Univ{Level: 0}, nil
	case Univ:
		// Type0 : Type1 but we simplify
		return Univ{Level: u.Level + 1}, nil
	case Zero:
		return Var{Name: "Nat"}, nil
	case Succ:
		// check inner is Nat
		if err := checkAgainst(ctx, venv, u.N, Var{Name: "Nat"}); err != nil {
			return nil, err
		}
		return Var{Name: "Nat"}, nil
	case NatRec:
		// Check P : Nat -> IType
		pType, err := typeCheck(ctx, venv, u.P)
		if err != nil {
			return nil, err
		}
		// expect pType to be Pi (k : Nat) -> IType
		pTypeN := eval(venv, pType)
		pi, ok := pTypeN.(Pi)
		if !ok {
			return nil, fmt.Errorf("P must be a Pi from Nat to IType; got %s", pTypeN.String())
		}
		// check Z : P Zero
		zExpected := subst(pi.Return, pi.Param, Zero{})
		if err := checkAgainst(ctx, venv, u.Z, zExpected); err != nil {
			return nil, fmt.Errorf("Z check failed: %w", err)
		}
		// check S : forall k. P k -> P (Succ k)
		// S should have type Pi(k:Nat). Pi(_ : P k). P (Succ k)
		sType, err := typeCheck(ctx, venv, u.S)
		if err != nil {
			return nil, err
		}
		_ = sType
		// Check n : Nat
		if err := checkAgainst(ctx, venv, u.N, Var{Name: "Nat"}); err != nil {
			return nil, err
		}
		// return P n
		return App{Fun: u.P, Arg: u.N}, nil
	default:
		return nil, fmt.Errorf("unknown term in typeCheck: %T", t)
	}
}

// checkAgainst: ensure term t has type expected (by normalizing both and checking equality)
func checkAgainst(ctx Env, venv Env, t Term, expected Term) error {
	got, err := typeCheck(ctx, venv, t)
	if err != nil {
		return err
	}
	gn := eval(venv, got)
	en := eval(venv, expected)
	if equalTerms(gn, en) {
		return nil
	}
	return fmt.Errorf("type mismatch: expected %s but got %s", en.String(), gn.String())
}

// equalTerms: naive structural equality on normalized terms
func equalTerms(a Term, b Term) bool {
	aStr := a.String()
	bStr := b.String()
	return aStr == bStr
}

// -------------------- Very Small Parse --------------------

// We implement just enough parsing for examples used in main. Not robust.

func tokenize(s string) []string {
	s = strings.ReplaceAll(s, "(", " ( ")
	s = strings.ReplaceAll(s, ")", " ) ")
	toks := strings.Fields(s)
	return toks
}

// parseTerm parses a token list into a Term; supports forms used below
func parseTerm(tokens []string, pos *int) (Term, error) {
	if *pos >= len(tokens) {
		return nil, errors.New("unexpected EOF")
	}
	ok := tokens[*pos]
	// parentheses
	if ok == "(" {
		*pos++
		// lookahead simple forms
		if *pos < len(tokens) && tokens[*pos] == "\\" { // lambda syntax (\x : A . body)
			*pos++
			if *pos >= len(tokens) {
				return nil, errors.New("expected param")
			}
			param := tokens[*pos]
			*pos++
			if tokens[*pos] != ":" {
				return nil, errors.New("expected :")
			}
			*pos++
			ptype, err := parseTerm(tokens, pos)
			if err != nil {
				return nil, err
			}
			if tokens[*pos] != "." {
				return nil, errors.New("expected . after lambda type")
			}
			*pos++
			body, err := parseTerm(tokens, pos)
			if err != nil {
				return nil, err
			}
			if tokens[*pos] != ")" {
				return nil, errors.New("expected )")
			}
			*pos++
			return Lam{Param: param, ParamType: ptype, Body: body}, nil
		}
		// simple function application or grouping: parse first then others
		first, err := parseTerm(tokens, pos)
		if err != nil {
			return nil, err
		}
		// parse possible second as argument
		second, err := parseTerm(tokens, pos)
		if err != nil {
			return nil, err
		}
		if tokens[*pos] != ")" {
			return nil, errors.New("expected ) at application")
		}
		*pos++
		return App{Fun: first, Arg: second}, nil
	}
	// special names
	switch ok {
	case "Type0":
		*pos++
		return Univ{Level: 0}, nil
	case "Nat":
		*pos++
		return Var{Name: "Nat"}, nil
	case "Zero":
		*pos++
		return Zero{}, nil
	case "Succ":
		*pos++
		if *pos >= len(tokens) {
			return nil, errors.New("Succ expects argument")
		}
		argTok := tokens[*pos]
		*pos++
		// parse simple numeric literal
		if n, err := strconv.Atoi(argTok); err == nil {
			// build nested Succ
			var term Term = Zero{}
			for i := 0; i < n; i++ {
				term = Succ{N: term}
			}
			return term, nil
		}
		return Var{Name: argTok}, nil
	default:
		// identifier or numeric literal
		*pos++
		if n, err := strconv.Atoi(ok); err == nil {
			var term Term = Zero{}
			for i := 0; i < n; i++ {
				term = Succ{N: term}
			}
			return term, nil
		}
		return Var{Name: ok}, nil
	}
}

// -------------------- Example and Main --------------------

func main() {
	// Context with Nat : Type0
	ctx := make(Env)
	ctx["Nat"] = Univ{Level: 0}

	// Value environment initially empty
	venv := make(Env)

	// Example: 2 = Succ (Succ Zero)
	two := Succ{N: Succ{N: Zero{}}}

	// Example: simple function: \x:Nat. x
	id := Lam{Param: "x", ParamType: Var{Name: "Nat"}, Body: Var{Name: "x"}}

	// Check id has type Pi(x:Nat). Nat
	typeOfId, err := typeCheck(ctx, venv, id)
	if err != nil {
		fmt.Println("type error:", err)
		return
	}
	fmt.Println("id type:", typeOfId.String())

	// Apply id to 2
	app := App{Fun: id, Arg: two}
	typApp, err := typeCheck(ctx, venv, app)
	if err != nil {
		fmt.Println("type error app:", err)
		return
	}
	fmt.Println("(id 2) type:", typApp.String())

	// A tiny NatRec example: define plus1 using NatRec
	// plus1 := \n:Nat. NatRec (\k:Nat. Nat) Zero (\k:Nat. \rec:Nat. Succ rec) n
	P := Lam{Param: "k", ParamType: Var{Name: "Nat"}, Body: Var{Name: "Nat"}} // P(k)=Nat (non-dependent simplification)
	zcase := Zero{}
	scase := Lam{Param: "k", ParamType: Var{Name: "Nat"}, Body: Lam{Param: "r", ParamType: Var{Name: "Nat"}, Body: Succ{N: Var{Name: "r"}}}}
	plus1 := Lam{Param: "n", ParamType: Var{Name: "Nat"}, Body: NatRec{P: App{Fun: P, Arg: Var{Name: "n"}}, Z: zcase, S: scase, N: Var{Name: "n"}}}

	// typecheck plus1
	typPlus1, err := typeCheck(ctx, venv, plus1)
	if err != nil {
		fmt.Println("plus1 type error:", err)
		return
	}
	fmt.Println("plus1 type:", typPlus1.String())

	// apply plus1 to 2 -> should yield 3
	threeTerm := App{Fun: plus1, Arg: two}
	res := eval(venv, threeTerm)
	fmt.Println("plus1 2 eval ->", res.String())

	// Note: This skeleton is intentionally simplified. It uses naive substitution and
	// string-based term equality. To make this production-ready:
	// - Use De Bruijn indices or another capture-avoiding representation
	// - Implement normalization by evaluation (NbE) for efficient conversion checking
	// - Implement universe levels and cumulativity correctly
	// - Add positivity checks for inductive definitions
	// - Add a proper parser, pattern matching and tactics if desired
}
