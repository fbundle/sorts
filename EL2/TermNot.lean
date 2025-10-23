import EL2.Term

namespace EL2
-- TODO - move these notations into another file
-- make them function
notation "univ" x => Term.t (T.univ x)
notation "var" x => Term.t (T.var x)
notation "inh" x => Term.t (T.inh x)
notation "typ" x => Term.t (T.typ x)
notation "bnd" x => Term.t (T.bnd x)
notation "lam" x => Term.t (T.lam x)
notation "app" x => Term.t (T.app x)
notation "mat" x => Term.t (T.mat x)

end EL2
