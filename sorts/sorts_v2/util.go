package sorts

func Test() {
	var ss = NewSortSystem("type")

	induceLeft := func(x Sort) Sort {
		// get type of parent (A \times B) -> return a term for A
		p := x.Parent().Sort().(prod)
		return ss.Atom(p.a.Level()-1, "dummy", p.a).MustUnwrap()
	}
	prop := ss.Atom(1, "Prop", nil).MustUnwrap()
	trueProp := ss.Atom(0, "True", prop).MustUnwrap()
	trueProof := ss.Atom(-1, "true_proof", trueProp).MustUnwrap()

	falseProp := ss.Atom(0, "False", prop).MustUnwrap() // non-inhabited

	_ = induceLeft
	_ = trueProof
	_ = falseProp

}
