# SORTS

started from some basic dependent type in (lean4](https://github.com/leanprover/lean4), I stated this project with the goal to implement the full dependent type system so that it is capable for mathemtical proof. This is probably a decade-long project, hope it would last.  


# EXAMPLES

## PROPOSITION LOGIC

proving `(Q or P) -> (P or Q)`

```go
package main

import (
	"fmt"

	sorts "github.com/fbundle/sorts/sorts_v3"
)

func main() {
	prop := sorts.NewAtom(1, "Prop", nil)

	P := sorts.NewAtom(0, "P", prop)
	Q := sorts.NewAtom(0, "Q", prop)

	PorQ := sorts.Sum{P, Q}
	QorP := sorts.Sum{Q, P}

	P_implies_PorQ := sorts.Arrow{P, PorQ}
	Q_implies_PorQ := sorts.Arrow{Q, PorQ}
	QorP_implies_PoQ := sorts.Arrow{QorP, PorQ}

	// I did it 😅 - x is a proof for (Q or P) -> (P or Q)
	x := QorP_implies_PoQ.Intro(func(term_QorP sorts.Sort) sorts.Sort {
		term_P_implies_PorQ := P_implies_PorQ.Intro(func(term_P sorts.Sort) sorts.Sort {
			return PorQ.IntroLeft(term_P)
		})
		term_Q_implies_PorQ := Q_implies_PorQ.Intro(func(term_Q sorts.Sort) sorts.Sort {
			return PorQ.IntroRight(term_Q)
		})

		return QorP.ByCases(term_QorP, term_Q_implies_PorQ, term_P_implies_PorQ)
	})

	fmt.Println(sorts.Name(sorts.Parent(x)))
}
```

