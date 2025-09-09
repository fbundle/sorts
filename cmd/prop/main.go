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
	QorP_implies_PorQ := sorts.Arrow{QorP, PorQ}

	// I did it 😅 - x is a proof for (Q or P) -> (P or Q)
	x := QorP_implies_PorQ.Intro("proof1", func(term_QorP sorts.Sort) sorts.Sort {
		term_P_implies_PorQ := P_implies_PorQ.Intro("proof2", func(term_P sorts.Sort) sorts.Sort {
			return PorQ.Intro(term_P, nil)
		})
		term_Q_implies_PorQ := Q_implies_PorQ.Intro("proof3", func(term_Q sorts.Sort) sorts.Sort {
			return PorQ.Intro(nil, term_Q)
		})

		term_PorQ := QorP.ByCases(term_QorP, term_Q_implies_PorQ, term_P_implies_PorQ)
		return term_PorQ
	})

	fmt.Println(sorts.Level(x), sorts.Name(x), sorts.Name(sorts.Parent(x)))
}
