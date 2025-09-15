package main

import (
	"fmt"
)

// https://chatgpt.com/share/68c81f91-4248-8009-a62d-f4d050fb4937

// ---------- AST ----------

// Form is the general AST node interface.
type Form interface {
	String() string
}

// Sort represents the universe of types.
type Sort struct {
	Level int
}

func (s *Sort) String() string {
	return fmt.Sprintf("Sort%d", s.Level)
}

// Var is a variable (term or type variable).
type Var struct {
	Name string
}

func (v *Var) String() string {
	return v.Name
}

// Pi represents dependent function types: Π(x:A). B(x).
type Pi struct {
	Var string
	A   Form
	B   Form
}

func (p *Pi) String() string {
	return fmt.Sprintf("(Π (%s : %s). %s)", p.Var, p.A, p.B)
}

// Lambda represents λ-abstraction.
type Lambda struct {
	Var  string
	A    Form // type of Var
	Body Form
}

func (l *Lambda) String() string {
	return fmt.Sprintf("(λ (%s : %s). %s)", l.Var, l.A, l.Body)
}

// App represents application f t.
type App struct {
	F Form
	T Form
}

func (a *App) String() string {
	return fmt.Sprintf("(%s %s)", a.F, a.T)
}

// ---------- Context ----------

// Context maps variable names to their types.
type Context struct {
	Vars map[string]Form
	Prev *Context
}

func (ctx *Context) Push(name string, ty Form) *Context {
	return &Context{Vars: map[string]Form{name: ty}, Prev: ctx}
}

func (ctx *Context) Lookup(name string) (Form, bool) {
	for c := ctx; c != nil; c = c.Prev {
		if ty, ok := c.Vars[name]; ok {
			return ty, true
		}
	}
	return nil, false
}

// ---------- Substitution ----------

// Subst replaces variable x with replacement inside body.
func Subst(body Form, x string, replacement Form) Form {
	switch b := body.(type) {
	case *Var:
		if b.Name == x {
			return replacement
		}
		return b
	case *Pi:
		if b.Var == x {
			return b
		}
		return &Pi{
			Var: b.Var,
			A:   Subst(b.A, x, replacement),
			B:   Subst(b.B, x, replacement),
		}
	case *Lambda:
		if b.Var == x {
			return b
		}
		return &Lambda{
			Var:  b.Var,
			A:    Subst(b.A, x, replacement),
			Body: Subst(b.Body, x, replacement),
		}
	case *App:
		return &App{
			F: Subst(b.F, x, replacement),
			T: Subst(b.T, x, replacement),
		}
	default:
		return body
	}
}

// ---------- Typechecking ----------

func TypeOf(ctx *Context, f Form) (Form, error) {
	switch t := f.(type) {
	case *Sort:
		return &Sort{Level: t.Level + 1}, nil

	case *Var:
		ty, ok := ctx.Lookup(t.Name)
		if !ok {
			return nil, fmt.Errorf("unbound variable: %s", t.Name)
		}
		return ty, nil

	case *Pi:
		// Check A : Sort
		aTy, err := TypeOf(ctx, t.A)
		if err != nil {
			return nil, err
		}
		if _, ok := aTy.(*Sort); !ok {
			return nil, fmt.Errorf("domain %s is not a type", t.A)
		}

		// Check B : Sort under x:A
		ctx2 := ctx.Push(t.Var, t.A)
		bTy, err := TypeOf(ctx2, t.B)
		if err != nil {
			return nil, err
		}
		if _, ok := bTy.(*Sort); !ok {
			return nil, fmt.Errorf("codomain %s is not a type", t.B)
		}

		return &Sort{Level: 0}, nil

	case *Lambda:
		// Expect A : Sort
		aTy, err := TypeOf(ctx, t.A)
		if err != nil {
			return nil, err
		}
		if _, ok := aTy.(*Sort); !ok {
			return nil, fmt.Errorf("lambda parameter %s is not a type: %s", t.Var, t.A)
		}

		// Body under x:A
		ctx2 := ctx.Push(t.Var, t.A)
		bodyTy, err := TypeOf(ctx2, t.Body)
		if err != nil {
			return nil, err
		}

		return &Pi{Var: t.Var, A: t.A, B: bodyTy}, nil

	case *App:
		fTy, err := TypeOf(ctx, t.F)
		if err != nil {
			return nil, err
		}
		pi, ok := fTy.(*Pi)
		if !ok {
			return nil, fmt.Errorf("function position not a Pi type: %s", fTy)
		}

		// Argument must match domain
		tTy, err := TypeOf(ctx, t.T)
		if err != nil {
			return nil, err
		}
		if fmt.Sprint(tTy) != fmt.Sprint(pi.A) {
			return nil, fmt.Errorf("argument type %s does not match domain %s", tTy, pi.A)
		}

		// Result type = B[x := t]
		return Subst(pi.B, pi.Var, t.T), nil
	}

	return nil, fmt.Errorf("unknown term: %T", f)
}

// ---------- Example ----------

func main() {
	ctx := &Context{}

	// Π(A:Sort0). Π(x:A). A
	idType := &Pi{
		Var: "A",
		A:   &Sort{Level: 0},
		B: &Pi{
			Var: "x",
			A:   &Var{Name: "A"},
			B:   &Var{Name: "A"},
		},
	}

	// λA. λx:A. x
	id := &Lambda{
		Var: "A",
		A:   &Sort{Level: 0},
		Body: &Lambda{
			Var:  "x",
			A:    &Var{Name: "A"},
			Body: &Var{Name: "x"},
		},
	}

	ty, err := TypeOf(ctx, id)
	if err != nil {
		panic(err)
	}

	fmt.Println("Identity function:", id)
	fmt.Println("Has type:         ", ty)
	fmt.Println("Expected type:    ", idType)
}
