package main

import sorts "github.com/fbundle/sorts/sorts_v3"

func main() {
	prop := sorts.NewAtom(1, "Prop", nil)

	P := sorts.NewAtom(0, "P", prop)
	Q := sorts.NewAtom(0, "Q", prop)

	PorQ := sorts.Sum{P, Q}
	QorP := sorts.Sum{Q, P}

	P_implies_PorQ := sorts.Arrow{P, PorQ}
	Q_implies_PorQ := sorts.Arrow{Q, PorQ}
	QorP_implies_PoQ := sorts.Arrow{QorP, PorQ}

	QorP_implies_PoQ.Intro(func(term_QorP sorts.Sort) sorts.Sort {
		term_Q := QorP.IntroLeft(term_QorP)
		term_P := QorP.IntroLeft(term_QorP)
		term_Q_implies_PorQ := PorQ.IntroLeft(term_Q)
		term_P_implies_PorQ := PorQ.IntroLeft(term_P)

		PorQ.ByCases()
	})


	myProp :=

}
