import EL2.Term

-- some COE to make writing code less painful
namespace EL2

instance: Coe String (Term β) where
  coe := (.var ·)

instance: Coe (Lst (Term β)) (Term β) where
  coe := (.lst ·)

instance: Coe (BindVal (Term β)) (Term β) where
  coe := (.bind_val ·)

instance: Coe (BindTyp (Term β)) (Term β) where
  coe := (.bind_typ ·)

instance: Coe (BindMk (Term β)) (Term β) where
  coe := (.bind_mk ·)

instance: Coe (Lam (Term β)) (Term β) where
  coe := (.lam ·)

instance: Coe (App (Term β) (Term β)) (Term β) where
  coe := (.app ·)

instance: Coe (Mat (Term β)) (Term β) where
  coe := (.mat ·)

end EL2
