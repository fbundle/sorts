import EL2.Term.Term
import EL2.Term.TermUtil
import EL2.Term.Util
import EL2.Term.Print

namespace EL2.Term.Infer


structure InferedTerm where
  term?: Option Term
  type: Term -- type of term
  level: Int -- level of term
  deriving Repr


end EL2.Term.Infer
