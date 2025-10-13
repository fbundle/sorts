structure Atom where
  name: String

structure Annot (α: Type) where
  name: String
  type: α

structure Pi (α: Type) where
  params: List (Annot α)
  body: α
