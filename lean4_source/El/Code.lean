structure Beta (α: Type) where
  cmd: α
  args: List α

structure Annot (α: Type) where
  name: String
  type: α

structure Binding (α: Type) where
  name: String
  value: α

structure Type (α: Type) where
  value: α

structure Inhabited (α: Type) where
  value: α

structure Pi (α: Type) where
  params: List (Annot α)
  body: α

structure Inductive (α: Type) where
  type: Annot α
  cons: List (Annot α)
