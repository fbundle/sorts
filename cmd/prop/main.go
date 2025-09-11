package main

import (
	"fmt"

	sorts2 "github.com/fbundle/sorts/obsolete/sorts_v3"
)

func main() {
	prop := sorts2.NewAtom(1, "Prop", nil)

	P := sorts2.NewAtom(0, "P", prop)
	Q := sorts2.NewAtom(0, "Q", prop)

	PorQ := sorts2.Sum{P, Q}
	QorP := sorts2.Sum{Q, P}

	P_implies_PorQ := sorts2.Arrow{P, PorQ}
	Q_implies_PorQ := sorts2.Arrow{Q, PorQ}
	QorP_implies_PorQ := sorts2.Arrow{QorP, PorQ}

	// I did it 😅 - x is a proof for (Q or P) -> (P or Q)
	x := QorP_implies_PorQ.Intro("proof1", func(term_QorP sorts2.Sort) sorts2.Sort {
		term_P_implies_PorQ := P_implies_PorQ.Intro("proof2", func(term_P sorts2.Sort) sorts2.Sort {
			return PorQ.Intro(term_P, nil)
		})
		term_Q_implies_PorQ := Q_implies_PorQ.Intro("proof3", func(term_Q sorts2.Sort) sorts2.Sort {
			return PorQ.Intro(nil, term_Q)
		})

		term_PorQ := QorP.ByCases(term_QorP, term_Q_implies_PorQ, term_P_implies_PorQ)
		return term_PorQ
	})

	fmt.Println(sorts2.Level(x), sorts2.Name(x))
}
